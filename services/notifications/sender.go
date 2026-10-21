package notifications

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"iuno-api/db"
	"iuno-api/services/wordoftheday"

	"github.com/jackc/pgx/v5"
	"github.com/sideshow/apns2"
	apnstoken "github.com/sideshow/apns2/token"
)

const (
	batchSize  = 100
	retryDelay = 5 * time.Minute
	claimDelay = 10 * time.Minute
)

type Device struct {
	ID          int64
	Token       string
	Environment string
	Timezone    string
}

type Sender struct {
	topic   string
	clients map[string]*apns2.Client
}

func NewSender() (*Sender, error) {
	keyID := strings.TrimSpace(os.Getenv("APNS_KEY_ID"))
	teamID := strings.TrimSpace(os.Getenv("APNS_TEAM_ID"))
	privateKey := strings.TrimSpace(os.Getenv("APNS_PRIVATE_KEY"))
	topic := strings.TrimSpace(os.Getenv("APNS_TOPIC"))

	privateKey = strings.ReplaceAll(privateKey, `\n`, "\n")

	if keyID == "" || teamID == "" || privateKey == "" || topic == "" {
		return nil, errors.New("missing APNs environment variables")
	}

	authKey, err := apnstoken.AuthKeyFromBytes([]byte(privateKey))
	if err != nil {
		return nil, fmt.Errorf("invalid APNs private key: %w", err)
	}

	authToken := &apnstoken.Token{
		AuthKey: authKey,
		KeyID:   keyID,
		TeamID:  teamID,
	}

	return &Sender{
		topic: topic,
		clients: map[string]*apns2.Client{
			"production": apns2.NewTokenClient(authToken).Production(),
			"sandbox":    apns2.NewTokenClient(authToken).Development(),
		},
	}, nil
}

type alertPayload struct {
	APS struct {
		Alert struct {
			Title string `json:"title"`
			Body  string `json:"body"`
		} `json:"alert"`
		Sound string `json:"sound"`
	} `json:"aps"`

	Type  string `json:"type"`
	Lemma string `json:"lemma"`
}

func (s *Sender) send(
	device Device,
	lemma string,
	meaning string,
) error {
	client, ok := s.clients[device.Environment]
	if !ok {
		return fmt.Errorf("unsupported APNs environment")
	}

	var payload alertPayload
	payload.APS.Alert.Title = "Word of the Day: " + lemma
	payload.APS.Alert.Body = meaning
	payload.APS.Sound = "default"
	payload.Type = "word_of_the_day"
	payload.Lemma = lemma

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	n := &apns2.Notification{
		DeviceToken: device.Token,
		Topic:       s.topic,
		Payload:     body,
		PushType:    apns2.PushTypeAlert,
		Priority:    apns2.PriorityHigh,
	}

	response, err := client.Push(n)
	if err != nil {
		return fmt.Errorf("APNs request failed: %w", err)
	}

	if response.StatusCode != 200 {
		return fmt.Errorf(
			"APNs rejected notification: status=%d reason=%s",
			response.StatusCode,
			response.Reason,
		)
	}

	return nil
}

func nextDailyNotification(timezone string, now time.Time) (time.Time, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}

	localNow := now.In(location)
	next := time.Date(
		localNow.Year(),
		localNow.Month(),
		localNow.Day(),
		9, 0, 0, 0,
		location,
	)

	if !localNow.Before(next) {
		next = next.AddDate(0, 0, 1)
	}

	return next.UTC(), nil
}

func claimDueDevices(ctx context.Context) ([]Device, error) {
	tx, err := db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	rows, err := tx.Query(ctx, `
		WITH due AS (
			SELECT id
			FROM device_notification_preferences
			WHERE daily_word_enabled = TRUE
			  AND next_notification_at <= NOW()
			  AND (
				last_sent_local_date IS NULL
				OR last_sent_local_date <
				   (NOW() AT TIME ZONE timezone)::date
			  )
			ORDER BY next_notification_at
			LIMIT $1
			FOR UPDATE SKIP LOCKED
		)
		UPDATE device_notification_preferences AS p
		SET
			next_notification_at = NOW() + $2::interval,
			updated_at = NOW()
		FROM due
		WHERE p.id = due.id
		RETURNING
			p.id,
			p.device_token,
			p.apns_environment,
			p.timezone
	`, batchSize, claimDelay.String())
	if err != nil {
		return nil, err
	}

	var devices []Device

	for rows.Next() {
		var d Device
		if err := rows.Scan(
			&d.ID,
			&d.Token,
			&d.Environment,
			&d.Timezone,
		); err != nil {
			rows.Close()
			return nil, err
		}
		devices = append(devices, d)
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	rows.Close()

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return devices, nil
}

func Run(ctx context.Context) error {
	sender, err := NewSender()
	if err != nil {
		return err
	}

	// Use the same Word of the Day service as the API.
	word, err := wordoftheday.GetToday(ctx)
	if err != nil {
		return fmt.Errorf("get today's word: %w", err)
	}

	if len(word.Meanings) == 0 ||
		strings.TrimSpace(word.Meanings[0]) == "" {
		return errors.New("today's word has no meaning to send")
	}

	devices, err := claimDueDevices(ctx)
	if err != nil {
		return fmt.Errorf("claim due devices: %w", err)
	}

	log.Printf("Word of the Day notification batch: %d devices", len(devices))

	for _, device := range devices {
		now := time.Now()

		location, err := time.LoadLocation(device.Timezone)
		if err != nil {
			log.Printf("Invalid timezone for notification record %d", device.ID)
			retryDevice(ctx, device.ID)
			continue
		}

		localDate := now.In(location).Format("2006-01-02")

		err = sender.send(
			device,
			word.Lemma,
			word.Meanings[0],
		)
		if err != nil {
			log.Printf(
				"Notification failed for record %d: %v",
				device.ID,
				err,
			)

			// APNs tells us this token is no longer registered.
			var responseErr interface {
				Error() string
			}
			_ = responseErr

			if strings.Contains(err.Error(), "status=410 reason=Unregistered") {
				disableDevice(ctx, device.ID)
			} else {
				retryDevice(ctx, device.ID)
			}
			continue
		}

		next, err := nextDailyNotification(device.Timezone, time.Now())
		if err != nil {
			retryDevice(ctx, device.ID)
			continue
		}

		_, err = db.Pool.Exec(ctx, `
			UPDATE device_notification_preferences
			SET
				last_sent_local_date = $2::date,
				next_notification_at = $3,
				updated_at = NOW()
			WHERE id = $1
		`, device.ID, localDate, next)

		if err != nil {
			// The push may already have been accepted by APNs.
			// Log this case for manual investigation.
			log.Printf(
				"Push accepted but delivery state update failed for record %d: %v",
				device.ID,
				err,
			)
		} else {
			log.Printf("Notification accepted by APNs for record %d", device.ID)
		}
	}

	return nil
}

func retryDevice(ctx context.Context, id int64) {
	_, err := db.Pool.Exec(ctx, `
		UPDATE device_notification_preferences
		SET
			next_notification_at = NOW() + $2::interval,
			updated_at = NOW()
		WHERE id = $1
		  AND daily_word_enabled = TRUE
	`, id, retryDelay.String())

	if err != nil {
		log.Printf("Could not schedule retry for record %d: %v", id, err)
	}
}

func disableDevice(ctx context.Context, id int64) {
	_, err := db.Pool.Exec(ctx, `
		UPDATE device_notification_preferences
		SET
			daily_word_enabled = FALSE,
			next_notification_at = NULL,
			updated_at = NOW()
		WHERE id = $1
	`, id)

	if err != nil {
		log.Printf("Could not disable invalid device record %d: %v", id, err)
	}
}
