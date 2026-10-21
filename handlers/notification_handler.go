package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"iuno-api/db"
	"iuno-api/middleware"
	"iuno-api/utils"
)

type notificationRequest struct {
	DeviceToken     string `json:"device_token"`
	APNSEnvironment string `json:"apns_environment"`
	Timezone        string `json:"timezone"`
	Enabled         *bool  `json:"enabled"`
}

func decodeNotificationRequest(
	w http.ResponseWriter,
	r *http.Request,
) (*notificationRequest, error) {
	defer r.Body.Close()

	var req notificationRequest
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16*1024))

	if err := decoder.Decode(&req); err != nil {
		return nil, err
	}

	req.DeviceToken = strings.TrimSpace(req.DeviceToken)
	req.APNSEnvironment = strings.TrimSpace(req.APNSEnvironment)
	req.Timezone = strings.TrimSpace(req.Timezone)

	if req.DeviceToken == "" {
		return nil, errors.New("device_token is required")
	}

	if len(req.DeviceToken) > 512 {
		return nil, errors.New("invalid device_token")
	}

	if req.APNSEnvironment == "" {
		req.APNSEnvironment = "production"
	}

	if req.APNSEnvironment != "sandbox" &&
		req.APNSEnvironment != "production" {
		return nil, errors.New("apns_environment must be sandbox or production")
	}

	if req.Timezone == "" {
		return nil, errors.New("timezone is required")
	}

	if _, err := time.LoadLocation(req.Timezone); err != nil {
		return nil, errors.New("invalid timezone")
	}

	if req.Enabled == nil {
		return nil, errors.New("enabled is required")
	}

	return &req, nil
}

func notificationUserID(r *http.Request) *int {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
	if !ok || claims == nil {
		return nil
	}

	id := claims.UserID
	return &id
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

func RegisterDeviceNotificationHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("NOTIFICATIONS: register handler reached")

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	req, err := decodeNotificationRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var nextNotificationAt *time.Time

	if *req.Enabled {
		next, err := nextDailyNotification(req.Timezone, time.Now())
		if err != nil {
			http.Error(w, "invalid timezone", http.StatusBadRequest)
			return
		}
		nextNotificationAt = &next
	}

	userID := notificationUserID(r)

	_, err = db.Pool.Exec(
		r.Context(),
		`
		INSERT INTO device_notification_preferences (
			user_id,
			device_token,
			apns_environment,
			timezone,
			daily_word_enabled,
			next_notification_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		ON CONFLICT (device_token, apns_environment)
		DO UPDATE SET
			user_id = EXCLUDED.user_id,
			timezone = EXCLUDED.timezone,
			daily_word_enabled = EXCLUDED.daily_word_enabled,
			next_notification_at = EXCLUDED.next_notification_at,
			updated_at = NOW()
		`,
		userID,
		req.DeviceToken,
		req.APNSEnvironment,
		req.Timezone,
		*req.Enabled,
		nextNotificationAt,
	)
	if err != nil {
		log.Printf("NOTIFICATIONS: database insert failed: %v", err)
		http.Error(w, "failed to register device", http.StatusInternalServerError)
		return
	}

	log.Println("NOTIFICATIONS: device registered successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success":              true,
		"enabled":              *req.Enabled,
		"next_notification_at": nextNotificationAt,
	})
}

func UpdateNotificationPreferencesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	req, err := decodeNotificationRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var nextNotificationAt *time.Time

	if *req.Enabled {
		next, err := nextDailyNotification(req.Timezone, time.Now())
		if err != nil {
			http.Error(w, "invalid timezone", http.StatusBadRequest)
			return
		}
		nextNotificationAt = &next
	}

	tag, err := db.Pool.Exec(
		r.Context(),
		`
		UPDATE device_notification_preferences
		SET
			timezone = $1,
			daily_word_enabled = $2,
			next_notification_at = $3,
			updated_at = NOW()
		WHERE device_token = $4
		  AND apns_environment = $5
		`,
		req.Timezone,
		*req.Enabled,
		nextNotificationAt,
		req.DeviceToken,
		req.APNSEnvironment,
	)
	if err != nil {
		http.Error(w, "failed to update notification preferences", http.StatusInternalServerError)
		return
	}

	if tag.RowsAffected() == 0 {
		http.Error(w, "device not registered", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success":              true,
		"enabled":              *req.Enabled,
		"next_notification_at": nextNotificationAt,
	})
}
