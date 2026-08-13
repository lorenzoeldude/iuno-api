package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"iuno-api/db"

	"github.com/jackc/pgx/v5"
	stripeSDK "github.com/stripe/stripe-go/v86"
	"github.com/stripe/stripe-go/v86/webhook"
)

func StripeWebhookHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	// =====================================================
	// WEBHOOK SECRET
	// =====================================================

	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	if webhookSecret == "" {
		log.Println("STRIPE_WEBHOOK_SECRET is missing")

		http.Error(
			w,
			"webhook secret is not configured",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// READ RAW BODY
	// =====================================================

	payload, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(
			w,
			"failed to read request body",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================================
	// VERIFY STRIPE SIGNATURE
	// =====================================================

	signature := r.Header.Get("Stripe-Signature")

	event, err := webhook.ConstructEvent(
		payload,
		signature,
		webhookSecret,
	)

	if err != nil {
		log.Printf(
			"Stripe webhook signature verification failed: %v",
			err,
		)

		http.Error(
			w,
			"invalid webhook signature",
			http.StatusBadRequest,
		)

		return
	}

	log.Printf(
		"Stripe webhook received: %s (%s)",
		event.Type,
		event.ID,
	)

	// =====================================================
	// BEGIN DATABASE TRANSACTION
	// =====================================================

	tx, err := db.Pool.Begin(r.Context())

	if err != nil {
		log.Printf(
			"Failed to begin Stripe webhook transaction: %v",
			err,
		)

		http.Error(
			w,
			"failed to begin transaction",
			http.StatusInternalServerError,
		)

		return
	}

	committed := false

	defer func() {
		if !committed {
			if rollbackErr := tx.Rollback(r.Context()); rollbackErr != nil &&
				!errors.Is(rollbackErr, pgx.ErrTxClosed) {

				log.Printf(
					"Failed to rollback Stripe webhook transaction: %v",
					rollbackErr,
				)
			}
		}
	}()

	// =====================================================
	// ATOMIC IDEMPOTENCY CLAIM
	// =====================================================
	//
	// Only one concurrent request can successfully INSERT
	// this Stripe event ID.
	//
	// If another request is already processing the same
	// event, this INSERT returns no row.
	//
	// The unique constraint on stripe_event_id makes this
	// safe even when Stripe delivers the event concurrently.
	// =====================================================

	var webhookEventID int

	err = tx.QueryRow(
		r.Context(),
		`
		INSERT INTO stripe_webhook_events (
			stripe_event_id,
			event_type
		)
		VALUES ($1, $2)
		ON CONFLICT (stripe_event_id) DO NOTHING
		RETURNING id
		`,
		event.ID,
		event.Type,
	).Scan(&webhookEventID)

	if errors.Is(err, pgx.ErrNoRows) {

		log.Printf(
			"Ignoring already processed Stripe event: %s",
			event.ID,
		)

		// Nothing was changed by this transaction.
		// Rollback is intentional.
		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		w.WriteHeader(http.StatusOK)

		w.Write(
			[]byte(`{"received":true,"duplicate":true}`),
		)

		return
	}

	if err != nil {

		log.Printf(
			"Failed to claim Stripe webhook event %s: %v",
			event.ID,
			err,
		)

		http.Error(
			w,
			"failed to record webhook event",
			http.StatusInternalServerError,
		)

		return
	}

	log.Printf(
		"Claimed Stripe webhook event %s as database event %d",
		event.ID,
		webhookEventID,
	)

	// =====================================================
	// HANDLE EVENTS
	// =====================================================

	switch event.Type {

	// =====================================================
	// CHECKOUT COMPLETED
	// =====================================================

	case "checkout.session.completed":

		var session stripeSDK.CheckoutSession

		if err := json.Unmarshal(
			event.Data.Raw,
			&session,
		); err != nil {

			log.Printf(
				"Failed to parse checkout session: %v",
				err,
			)

			http.Error(
				w,
				"failed to parse event",
				http.StatusBadRequest,
			)

			return
		}

		log.Printf(
			"Checkout completed: %s",
			session.ID,
		)

		// =================================================
		// USER ID
		// =================================================

		userIDString := session.Metadata["user_id"]

		if userIDString == "" {

			log.Printf(
				"Checkout session %s has no user_id metadata",
				session.ID,
			)

			http.Error(
				w,
				"missing user_id metadata",
				http.StatusBadRequest,
			)

			return
		}

		userID, err := strconv.Atoi(userIDString)

		if err != nil {

			log.Printf(
				"Invalid user_id: %s",
				userIDString,
			)

			http.Error(
				w,
				"invalid user_id",
				http.StatusBadRequest,
			)

			return
		}

		// =================================================
		// PRICE ID
		// =================================================

		priceID := session.Metadata["price_id"]

		if priceID == "" {

			log.Printf(
				"Checkout session %s has no price_id metadata",
				session.ID,
			)

			http.Error(
				w,
				"missing price_id metadata",
				http.StatusBadRequest,
			)

			return
		}

		// =================================================
		// CUSTOMER
		// =================================================

		customerID := ""

		if session.Customer != nil {
			customerID = session.Customer.ID
		}

		// =================================================
		// SUBSCRIPTION
		// =================================================

		subscriptionID := ""

		if session.Subscription != nil {
			subscriptionID = session.Subscription.ID
		}

		log.Printf(
			"Stripe checkout user=%d customer=%s subscription=%s price=%s",
			userID,
			customerID,
			subscriptionID,
			priceID,
		)

		// =================================================
		// ACTIVATE PREMIUM
		// =================================================

		_, err = tx.Exec(
			r.Context(),
			`
			UPDATE users
			SET is_premium = true
			WHERE id = $1
			`,
			userID,
		)

		if err != nil {

			log.Printf(
				"Failed to activate Premium for user %d: %v",
				userID,
				err,
			)

			http.Error(
				w,
				"failed to update user",
				http.StatusInternalServerError,
			)

			return
		}

		// =================================================
		// SAVE CUSTOMER
		// =================================================

		if customerID != "" {

			_, err = tx.Exec(
				r.Context(),
				`
				INSERT INTO billing_customers (
					user_id,
					provider,
					provider_customer_id
				)
				VALUES (
					$1,
					'stripe',
					$2
				)
				ON CONFLICT (user_id, provider)
				DO UPDATE SET
					provider_customer_id =
						EXCLUDED.provider_customer_id
				`,
				userID,
				customerID,
			)

			if err != nil {

				log.Printf(
					"Failed to save Stripe customer: %v",
					err,
				)

				http.Error(
					w,
					"failed to save billing customer",
					http.StatusInternalServerError,
				)

				return
			}
		}

		// =================================================
		// SAVE SUBSCRIPTION
		// =================================================

		if subscriptionID != "" {

			_, err = tx.Exec(
				r.Context(),
				`
				INSERT INTO subscriptions (
					user_id,
					provider,
					provider_subscription_id,
					product_id,
					status
				)
				VALUES (
					$1,
					'stripe',
					$2,
					$3,
					'active'
				)
				ON CONFLICT (
					provider,
					provider_subscription_id
				)
				DO UPDATE SET
					status = 'active',
					product_id = EXCLUDED.product_id,
					updated_at = now()
				`,
				userID,
				subscriptionID,
				priceID,
			)

			if err != nil {

				log.Printf(
					"Failed to save subscription: %v",
					err,
				)

				http.Error(
					w,
					"failed to save subscription",
					http.StatusInternalServerError,
				)

				return
			}
		}

		log.Printf(
			"Premium activated successfully for user %d",
			userID,
		)

	// =====================================================
	// SUBSCRIPTION UPDATED
	// =====================================================

	case "customer.subscription.updated":

		var subscription stripeSDK.Subscription

		if err := json.Unmarshal(
			event.Data.Raw,
			&subscription,
		); err != nil {

			log.Printf(
				"Failed to parse subscription.updated: %v",
				err,
			)

			http.Error(
				w,
				"failed to parse subscription",
				http.StatusBadRequest,
			)

			return
		}

		subscriptionID := subscription.ID
		status := string(subscription.Status)

		log.Printf(
			"Subscription updated: %s status=%s",
			subscriptionID,
			status,
		)

		// =================================================
		// GET USER
		// =================================================

		var userID int

		err = tx.QueryRow(
			r.Context(),
			`
			SELECT user_id
			FROM subscriptions
			WHERE provider = 'stripe'
			  AND provider_subscription_id = $1
			`,
			subscriptionID,
		).Scan(&userID)

		if err != nil {

			log.Printf(
				"Could not find IUNO user for subscription %s: %v",
				subscriptionID,
				err,
			)

			// We deliberately acknowledge this event.
			// The checkout event should normally have created
			// the local subscription first.
			if err := tx.Commit(r.Context()); err != nil {

				log.Printf(
					"Failed to commit ignored subscription event: %v",
					err,
				)

				http.Error(
					w,
					"failed to commit transaction",
					http.StatusInternalServerError,
				)

				return
			}

			committed = true

			w.WriteHeader(http.StatusOK)
			return
		}

		// =================================================
		// PREMIUM STATUS
		// =================================================

		isPremium :=
			status == "active" ||
				status == "trialing"

		// =================================================
		// PERIOD DATES
		// =================================================

		var periodStart *time.Time
		var periodEnd *time.Time

		if subscription.Items != nil &&
			len(subscription.Items.Data) > 0 {

			item := subscription.Items.Data[0]

			if item.CurrentPeriodStart > 0 {

				t := time.Unix(
					item.CurrentPeriodStart,
					0,
				)

				periodStart = &t
			}

			if item.CurrentPeriodEnd > 0 {

				t := time.Unix(
					item.CurrentPeriodEnd,
					0,
				)

				periodEnd = &t
			}
		}

		// =================================================
		// UPDATE SUBSCRIPTION
		// =================================================

		_, err = tx.Exec(
			r.Context(),
			`
			UPDATE subscriptions
			SET
				status = $1,
				current_period_start = $2,
				current_period_end = $3,
				cancel_at_period_end = $4,
				updated_at = now()
			WHERE provider = 'stripe'
			  AND provider_subscription_id = $5
			`,
			status,
			periodStart,
			periodEnd,
			subscription.CancelAtPeriodEnd,
			subscriptionID,
		)

		if err != nil {

			log.Printf(
				"Failed to update subscription %s: %v",
				subscriptionID,
				err,
			)

			http.Error(
				w,
				"failed to update subscription",
				http.StatusInternalServerError,
			)

			return
		}

		// =================================================
		// UPDATE USER PREMIUM STATUS
		// =================================================

		_, err = tx.Exec(
			r.Context(),
			`
			UPDATE users
			SET is_premium = $1
			WHERE id = $2
			`,
			isPremium,
			userID,
		)

		if err != nil {

			log.Printf(
				"Failed to update Premium status for user %d: %v",
				userID,
				err,
			)

			http.Error(
				w,
				"failed to update premium status",
				http.StatusInternalServerError,
			)

			return
		}

		log.Printf(
			"Premium status for user %d: %v",
			userID,
			isPremium,
		)

	// =====================================================
	// SUBSCRIPTION DELETED
	// =====================================================

	case "customer.subscription.deleted":

		var subscription stripeSDK.Subscription

		if err := json.Unmarshal(
			event.Data.Raw,
			&subscription,
		); err != nil {

			log.Printf(
				"Failed to parse subscription.deleted: %v",
				err,
			)

			http.Error(
				w,
				"failed to parse subscription",
				http.StatusBadRequest,
			)

			return
		}

		subscriptionID := subscription.ID

		log.Printf(
			"Subscription deleted: %s",
			subscriptionID,
		)

		// =================================================
		// GET USER
		// =================================================

		var userID int

		err = tx.QueryRow(
			r.Context(),
			`
			SELECT user_id
			FROM subscriptions
			WHERE provider = 'stripe'
			  AND provider_subscription_id = $1
			`,
			subscriptionID,
		).Scan(&userID)

		if err != nil {

			log.Printf(
				"Could not find user for deleted subscription %s: %v",
				subscriptionID,
				err,
			)

			if err := tx.Commit(r.Context()); err != nil {

				log.Printf(
					"Failed to commit ignored deleted subscription event: %v",
					err,
				)

				http.Error(
					w,
					"failed to commit transaction",
					http.StatusInternalServerError,
				)

				return
			}

			committed = true

			w.WriteHeader(http.StatusOK)
			return
		}

		// =================================================
		// MARK SUBSCRIPTION CANCELED
		// =================================================

		_, err = tx.Exec(
			r.Context(),
			`
			UPDATE subscriptions
			SET
				status = 'canceled',
				cancel_at_period_end = false,
				updated_at = now()
			WHERE provider = 'stripe'
			  AND provider_subscription_id = $1
			`,
			subscriptionID,
		)

		if err != nil {

			log.Printf(
				"Failed to cancel subscription %s: %v",
				subscriptionID,
				err,
			)

			http.Error(
				w,
				"failed to cancel subscription",
				http.StatusInternalServerError,
			)

			return
		}

		// =================================================
		// REMOVE PREMIUM
		// =================================================

		_, err = tx.Exec(
			r.Context(),
			`
			UPDATE users
			SET is_premium = false
			WHERE id = $1
			`,
			userID,
		)

		if err != nil {

			log.Printf(
				"Failed to remove Premium from user %d: %v",
				userID,
				err,
			)

			http.Error(
				w,
				"failed to update premium status",
				http.StatusInternalServerError,
			)

			return
		}

		log.Printf(
			"Premium removed from user %d",
			userID,
		)

	default:

		log.Printf(
			"Ignoring Stripe event: %s",
			event.Type,
		)
	}

	// =====================================================
	// COMMIT EVERYTHING
	// =====================================================

	if err := tx.Commit(r.Context()); err != nil {

		log.Printf(
			"Failed to commit Stripe webhook transaction: %v",
			err,
		)

		http.Error(
			w,
			"failed to commit webhook transaction",
			http.StatusInternalServerError,
		)

		return
	}

	committed = true

	log.Printf(
		"Stripe webhook processed successfully: %s",
		event.ID,
	)

	// =====================================================
	// ACKNOWLEDGE WEBHOOK
	// =====================================================

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	w.Write(
		[]byte(`{"received":true}`),
	)
}
