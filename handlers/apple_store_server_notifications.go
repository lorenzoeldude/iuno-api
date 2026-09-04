package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"

	"iuno-api/storekit"
)

type AppleServerNotificationRequest struct {
	SignedPayload string `json:"signedPayload"`
}

// =====================================================
// HANDLER
// =====================================================

func AppleStoreServerNotificationsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	var request AppleServerNotificationRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if request.SignedPayload == "" {
		http.Error(
			w,
			"signedPayload is missing",
			http.StatusBadRequest,
		)
		return
	}

	// =====================================================
	// VERIFY / DECODE APPLE NOTIFICATION
	// =====================================================

	notification, err :=
		storekit.VerifyServerNotification(
			request.SignedPayload,
		)

	if err != nil {
		log.Printf(
			"Apple notification verification failed: %v",
			err,
		)

		http.Error(
			w,
			"invalid signed payload",
			http.StatusBadRequest,
		)
		return
	}

	log.Printf(
		"Apple notification: type=%s subtype=%s uuid=%s environment=%s",
		notification.NotificationType,
		notification.Subtype,
		notification.NotificationUUID,
		notification.DataEnvironment(),
	)

	// =====================================================
	// TEST NOTIFICATION
	// =====================================================

	if notification.NotificationType == "TEST" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// =====================================================
	// NO DATA
	// =====================================================

	if notification.Data == nil {
		log.Printf(
			"Apple notification %s contains no data",
			notification.NotificationUUID,
		)

		w.WriteHeader(http.StatusOK)
		return
	}

	// =====================================================
	// IDEMPOTENCY
	// =====================================================

	tx, err := storekit.BeginTransaction(r.Context())

	if err != nil {
		log.Printf(
			"failed to begin Apple notification transaction: %v",
			err,
		)

		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	defer tx.Rollback(r.Context())

	var alreadyProcessed bool

	err = tx.QueryRow(
		r.Context(),
		`
		SELECT EXISTS (
			SELECT 1
			FROM apple_server_notifications
			WHERE notification_uuid = $1
		)
		`,
		notification.NotificationUUID,
	).Scan(&alreadyProcessed)

	if err != nil {
		log.Printf(
			"failed to check Apple notification idempotency: %v",
			err,
		)

		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	if alreadyProcessed {
		log.Printf(
			"Apple notification %s already processed",
			notification.NotificationUUID,
		)

		w.WriteHeader(http.StatusOK)
		return
	}

	// =====================================================
	// DECODE TRANSACTION
	// =====================================================

	if notification.Data.SignedTransactionInfo == "" {
		log.Printf(
			"Apple notification %s has no signedTransactionInfo",
			notification.NotificationUUID,
		)

		// Still record the notification so Apple doesn't
		// repeatedly send something we cannot process.
		_, err = tx.Exec(
			r.Context(),
			`
			INSERT INTO apple_server_notifications (
				notification_uuid,
				notification_type,
				subtype,
				signed_date,
				environment
			)
			VALUES ($1, $2, $3, $4, $5)
			`,
			notification.NotificationUUID,
			notification.NotificationType,
			notification.Subtype,
			notification.SignedDate,
			notification.Data.Environment,
		)

		if err != nil {
			log.Printf(
				"failed to save Apple notification: %v",
				err,
			)

			http.Error(
				w,
				"internal server error",
				http.StatusInternalServerError,
			)
			return
		}

		if err := tx.Commit(r.Context()); err != nil {
			log.Printf(
				"failed to commit Apple notification: %v",
				err,
			)

			http.Error(
				w,
				"internal server error",
				http.StatusInternalServerError,
			)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	transaction, err :=
		storekit.VerifyTransaction(
			notification.Data.SignedTransactionInfo,
		)

	if err != nil {
		log.Printf(
			"failed to decode Apple signedTransactionInfo: %v",
			err,
		)

		http.Error(
			w,
			"invalid transaction info",
			http.StatusBadRequest,
		)
		return
	}

	// =====================================================
	// FIND SUBSCRIPTION OWNER
	// =====================================================

	userID, err :=
		storekit.GetSubscriptionOwner(
			r.Context(),
			tx,
			transaction.Payload.OriginalTransactionID,
		)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// This can happen if Apple sends a notification
			// before our purchase endpoint has associated the
			// subscription with a user.
			//
			// We acknowledge it rather than making Apple retry
			// forever.
			log.Printf(
				"Apple notification for unknown subscription: %s",
				transaction.Payload.OriginalTransactionID,
			)

			_, insertErr := tx.Exec(
				r.Context(),
				`
				INSERT INTO apple_server_notifications (
					notification_uuid,
					notification_type,
					subtype,
					signed_date,
					environment
				)
				VALUES ($1, $2, $3, $4, $5)
				`,
				notification.NotificationUUID,
				notification.NotificationType,
				notification.Subtype,
				notification.SignedDate,
				notification.Data.Environment,
			)

			if insertErr != nil {
				log.Printf(
					"failed to save unknown Apple notification: %v",
					insertErr,
				)

				http.Error(
					w,
					"internal server error",
					http.StatusInternalServerError,
				)
				return
			}

			if err := tx.Commit(r.Context()); err != nil {
				log.Printf(
					"failed to commit Apple notification: %v",
					err,
				)

				http.Error(
					w,
					"internal server error",
					http.StatusInternalServerError,
				)
				return
			}

			w.WriteHeader(http.StatusOK)
			return
		}

		log.Printf(
			"failed to find Apple subscription owner: %v",
			err,
		)

		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	// =====================================================
	// SAVE SUBSCRIPTION
	// =====================================================

	isPremium := storekit.IsPremium(transaction.Payload)

	if err := storekit.SaveSubscription(
		r.Context(),
		tx,
		userID,
		transaction,
		isPremium,
	); err != nil {
		log.Printf(
			"failed to save Apple subscription: %v",
			err,
		)

		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	// =====================================================
	// RECALCULATE PREMIUM
	// =====================================================

	activeAppleSubscription, err :=
		storekit.HasActiveAppleSubscription(
			r.Context(),
			tx,
			userID,
		)

	if err != nil {
		log.Printf(
			"failed to determine Apple premium status: %v",
			err,
		)

		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	if err := storekit.UpdateUserPremium(
		r.Context(),
		tx,
		userID,
		activeAppleSubscription,
	); err != nil {
		log.Printf(
			"failed to update user premium status: %v",
			err,
		)

		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	// =====================================================
	// SAVE NOTIFICATION
	// =====================================================

	_, err = tx.Exec(
		r.Context(),
		`
		INSERT INTO apple_server_notifications (
			notification_uuid,
			notification_type,
			subtype,
			signed_date,
			environment
		)
		VALUES ($1, $2, $3, $4, $5)
		`,
		notification.NotificationUUID,
		notification.NotificationType,
		notification.Subtype,
		notification.SignedDate,
		notification.Data.Environment,
	)

	if err != nil {
		log.Printf(
			"failed to save Apple notification: %v",
			err,
		)

		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	// =====================================================
	// COMMIT
	// =====================================================

	if err := tx.Commit(r.Context()); err != nil {
		log.Printf(
			"failed to commit Apple notification: %v",
			err,
		)

		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	log.Printf(
		"Apple notification processed: type=%s subtype=%s user=%d premium=%t",
		notification.NotificationType,
		notification.Subtype,
		userID,
		activeAppleSubscription,
	)

	w.WriteHeader(http.StatusOK)
}
