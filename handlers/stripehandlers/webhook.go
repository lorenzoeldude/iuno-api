package stripehandlers

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"iuno-api/db"

	"github.com/jackc/pgx/v5"
	stripeSDK "github.com/stripe/stripe-go/v86"
	"github.com/stripe/stripe-go/v86/webhook"
)

func StripeWebhookHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	// Read body
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	// Verify signature
	event, err := webhook.ConstructEvent(
		payload,
		r.Header.Get("Stripe-Signature"),
		os.Getenv("STRIPE_WEBHOOK_SECRET"),
	)
	if err != nil {
		log.Printf("Invalid Stripe webhook: %v", err)
		http.Error(w, "invalid webhook signature", http.StatusBadRequest)
		return
	}

	log.Printf(
		"Stripe webhook received: %s (%s)",
		event.Type,
		event.ID,
	)

	// Begin transaction
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		http.Error(w, "failed to begin transaction", http.StatusInternalServerError)
		return
	}

	committed := false

	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	// Idempotency
	var eventID int

	err = tx.QueryRow(
		ctx,
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
	).Scan(&eventID)

	if errors.Is(err, pgx.ErrNoRows) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"received":true,"duplicate":true}`))
		return
	}

	if err != nil {
		http.Error(w, "failed to record webhook event", http.StatusInternalServerError)
		return
	}

	// Dispatch
	if err := dispatchStripeEvent(ctx, tx, event); err != nil {

		log.Printf(
			"Failed to process Stripe event %s: %v",
			event.ID,
			err,
		)

		http.Error(
			w,
			"failed to process webhook",
			http.StatusInternalServerError,
		)

		return
	}

	// Commit
	if err := tx.Commit(ctx); err != nil {
		log.Printf(
			"Failed to commit Stripe webhook %s: %v",
			event.ID,
			err,
		)

		http.Error(
			w,
			"failed to commit webhook",
			http.StatusInternalServerError,
		)

		return
	}

	committed = true

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"received":true}`))
}

func dispatchStripeEvent(
	ctx context.Context,
	tx pgx.Tx,
	event stripeSDK.Event,
) error {

	switch event.Type {

	case "checkout.session.completed":
		return handleCheckoutSessionCompleted(ctx, tx, event)

	case "payment_intent.succeeded":
		return handlePaymentIntentSucceeded(ctx, tx, event)

	case "payment_intent.payment_failed":
		return handlePaymentIntentFailed(ctx, tx, event)

	case "invoice_payment.paid":
		return handleInvoicePaymentPaid(ctx, tx, event)

	case "invoice.paid":
		return handleInvoicePaid(ctx, tx, event)

	case "invoice.payment_failed":
		return handleInvoicePaymentFailed(ctx, tx, event)

	case "customer.subscription.created":
		return handleSubscriptionCreated(ctx, tx, event)

	case "customer.subscription.updated":
		return handleSubscriptionUpdated(ctx, tx, event)

	case "customer.subscription.deleted":
		return handleSubscriptionDeleted(ctx, tx, event)

	default:
		log.Printf("Ignoring Stripe event: %s", event.Type)
		return nil
	}
}
