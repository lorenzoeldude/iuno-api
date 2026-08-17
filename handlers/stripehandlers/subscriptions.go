package stripehandlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"log"
	"os"
	"time"

	"iuno-api/db"
	"iuno-api/middleware"
	"iuno-api/stripe"
	"iuno-api/utils"

	"github.com/jackc/pgx/v5"
	stripeSDK "github.com/stripe/stripe-go/v86"
)

// =====================================================
// CREATE STRIPE CUSTOMER PORTAL SESSION
// =====================================================

func CreateStripePortalSessionHandler(
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

	claims, ok := r.Context().Value(
		middleware.UserContextKey,
	).(*utils.Claims)

	if !ok || claims == nil {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	var customerID string

	err := db.Pool.QueryRow(
		r.Context(),
		`
		SELECT provider_customer_id
		FROM billing_customers
		WHERE user_id = $1
		  AND provider = 'stripe'
		`,
		claims.UserID,
	).Scan(&customerID)

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(
			w,
			"Stripe customer not found",
			http.StatusNotFound,
		)
		return
	}

	if err != nil {
		http.Error(
			w,
			"failed to find Stripe customer",
			http.StatusInternalServerError,
		)
		return
	}

	if customerID == "" {
		http.Error(
			w,
			"Stripe customer not found",
			http.StatusNotFound,
		)
		return
	}

	returnURL := os.Getenv("STRIPE_PORTAL_RETURN_URL")

	if returnURL == "" {
		returnURL = "https://iunoni.com/account"
	}

	params := &stripeSDK.BillingPortalSessionCreateParams{
		Customer:  stripeSDK.String(customerID),
		ReturnURL: stripeSDK.String(returnURL),
	}

	session, err := stripe.Client.V1BillingPortalSessions.Create(
		r.Context(),
		params,
	)

	if err != nil {
		http.Error(
			w,
			"failed to create Stripe portal session",
			http.StatusInternalServerError,
		)
		return
	}

	if session == nil || session.URL == "" {
		http.Error(
			w,
			"Stripe portal session has no URL",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	_ = json.NewEncoder(w).Encode(
		map[string]string{
			"url": session.URL,
		},
	)
}

// =====================================================
// GET BILLING STATUS
// =====================================================

func GetBillingStatusHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	claims, ok := r.Context().Value(
		middleware.UserContextKey,
	).(*utils.Claims)

	if !ok || claims == nil {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	var status struct {
		IsPremium          bool       `json:"is_premium"`
		Provider           *string    `json:"provider,omitempty"`
		SubscriptionID     *string    `json:"subscription_id,omitempty"`
		SubscriptionStatus *string    `json:"subscription_status,omitempty"`
		ProductID          *string    `json:"product_id,omitempty"`
		CurrentPeriodStart *time.Time `json:"current_period_start,omitempty"`
		CurrentPeriodEnd   *time.Time `json:"current_period_end,omitempty"`
		StripeCustomerID   *string    `json:"stripe_customer_id,omitempty"`
	}

	err := db.Pool.QueryRow(
		r.Context(),
		`
		SELECT
			u.is_premium,
			s.provider,
			s.provider_subscription_id,
			s.status,
			s.product_id,
			s.current_period_start,
			s.current_period_end,
			bc.provider_customer_id
		FROM users u
		LEFT JOIN LATERAL (
			SELECT
				provider,
				provider_subscription_id,
				status,
				product_id,
				current_period_start,
				current_period_end
			FROM subscriptions
			WHERE user_id = u.id
			  AND provider = 'stripe'
			ORDER BY created_at DESC
			LIMIT 1
		) s ON true
		LEFT JOIN LATERAL (
			SELECT provider_customer_id
			FROM billing_customers
			WHERE user_id = u.id
			  AND provider = 'stripe'
			ORDER BY created_at DESC
			LIMIT 1
		) bc ON true
		WHERE u.id = $1
		`,
		claims.UserID,
	).Scan(
		&status.IsPremium,
		&status.Provider,
		&status.SubscriptionID,
		&status.SubscriptionStatus,
		&status.ProductID,
		&status.CurrentPeriodStart,
		&status.CurrentPeriodEnd,
		&status.StripeCustomerID,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(
			w,
			"user not found",
			http.StatusNotFound,
		)
		return
	}

	if err != nil {
		http.Error(
			w,
			"failed to get billing status",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	_ = json.NewEncoder(w).Encode(status)
}

// =====================================================
// CHECKOUT SESSION COMPLETED
// =====================================================

func handleCheckoutSessionCompleted(
	ctx context.Context,
	tx pgx.Tx,
	event stripeSDK.Event,
) error {

	var session struct {
		ID           string            `json:"id"`
		Customer     string            `json:"customer"`
		Mode         string            `json:"mode"`
		Metadata     map[string]string `json:"metadata"`
		Subscription string            `json:"subscription"`
	}

	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		return err
	}

	if session.ID == "" {
		return errors.New(
			"checkout.session.completed has no session ID",
		)
	}

	userID, err := FindPaymentOwner(
		ctx,
		tx,
		session.Customer,
		session.Metadata,
	)

	if err != nil {
		return err
	}

	if session.Customer != "" {
		_, err = tx.Exec(
			ctx,
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
			ON CONFLICT (
				provider,
				provider_customer_id
			)
			DO UPDATE SET
				user_id = EXCLUDED.user_id
			`,
			userID,
			session.Customer,
		)

		if err != nil {
			return err
		}
	}

	if session.Mode == "subscription" &&
		session.Subscription != "" {

		_, err = tx.Exec(
			ctx,
			`
			UPDATE users
			SET is_premium = true
			WHERE id = $1
			`,
			userID,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

// =====================================================
// SUBSCRIPTION CREATED
// =====================================================

func handleSubscriptionCreated(
	ctx context.Context,
	tx pgx.Tx,
	event stripeSDK.Event,
) error {
	return upsertStripeSubscription(
		ctx,
		tx,
		event,
	)
}

// =====================================================
// SUBSCRIPTION UPDATED
// =====================================================

func handleSubscriptionUpdated(
	ctx context.Context,
	tx pgx.Tx,
	event stripeSDK.Event,
) error {
	return upsertStripeSubscription(
		ctx,
		tx,
		event,
	)
}

// =====================================================
// SUBSCRIPTION DELETED
// =====================================================

func handleSubscriptionDeleted(
	ctx context.Context,
	tx pgx.Tx,
	event stripeSDK.Event,
) error {

	var subscription struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(
		event.Data.Raw,
		&subscription,
	); err != nil {
		return err
	}

	if subscription.ID == "" {
		return errors.New(
			"customer.subscription.deleted has no subscription ID",
		)
	}

	var userID int

	err := tx.QueryRow(
		ctx,
		`
		SELECT user_id
		FROM subscriptions
		WHERE provider = 'stripe'
		  AND provider_subscription_id = $1
		`,
		subscription.ID,
	).Scan(&userID)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}

	if err != nil {
		return err
	}

	_, err = tx.Exec(
		ctx,
		`
		UPDATE subscriptions
		SET status = 'canceled',
		    current_period_end = NOW()
		WHERE provider = 'stripe'
		  AND provider_subscription_id = $1
		`,
		subscription.ID,
	)

	if err != nil {
		return err
	}

	_, err = tx.Exec(
		ctx,
		`
		UPDATE users
		SET is_premium = false
		WHERE id = $1
		`,
		userID,
	)

	return err
}

// =====================================================
// UPSERT STRIPE SUBSCRIPTION
// =====================================================

func upsertStripeSubscription(
	ctx context.Context,
	tx pgx.Tx,
	event stripeSDK.Event,
) error {

	var subscription struct {
		ID       string `json:"id"`
		Customer string `json:"customer"`
		Status   string `json:"status"`

		CurrentPeriodStart int64 `json:"current_period_start"`
		CurrentPeriodEnd   int64 `json:"current_period_end"`

		Metadata map[string]string `json:"metadata"`

		Items struct {
			Data []struct {
				Price struct {
					ID      string `json:"id"`
					Product string `json:"product"`
				} `json:"price"`
			} `json:"data"`
		} `json:"items"`
	}

	if err := json.Unmarshal(
		event.Data.Raw,
		&subscription,
	); err != nil {
		return err
	}

	log.Printf(
		"Stripe subscription decoded: id=%s customer=%s status=%s period_start=%d period_end=%d",
		subscription.ID,
		subscription.Customer,
		subscription.Status,
		subscription.CurrentPeriodStart,
		subscription.CurrentPeriodEnd,
	)

	if subscription.ID == "" {
		return errors.New(
			"Stripe subscription has no ID",
		)
	}

	userID, err := FindPaymentOwner(
		ctx,
		tx,
		subscription.Customer,
		subscription.Metadata,
	)

	if err != nil {
		return err
	}

	productID := ""

	if len(subscription.Items.Data) > 0 {
		productID = subscription.Items.Data[0].Price.ID

		if productID == "" {
			productID = subscription.Items.Data[0].Price.Product
		}
	}

	if productID == "" && subscription.Metadata != nil {
		productID = subscription.Metadata["price_id"]
	}

	if productID == "" {
		productID = "stripe_subscription"
	}

	var periodStart *time.Time
	var periodEnd *time.Time

	if subscription.CurrentPeriodStart > 0 {
		t := time.Unix(
			subscription.CurrentPeriodStart,
			0,
		)
		periodStart = &t
	}

	if subscription.CurrentPeriodEnd > 0 {
		t := time.Unix(
			subscription.CurrentPeriodEnd,
			0,
		)
		periodEnd = &t
	}

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO subscriptions (
			user_id,
			provider,
			provider_subscription_id,
			product_id,
			status,
			current_period_start,
			current_period_end
		)
		VALUES (
			$1,
			'stripe',
			$2,
			$3,
			$4,
			$5,
			$6
		)
		ON CONFLICT (
			provider,
			provider_subscription_id
		)
		DO UPDATE SET
			product_id = EXCLUDED.product_id,
			status = EXCLUDED.status,
			current_period_start = EXCLUDED.current_period_start,
			current_period_end = EXCLUDED.current_period_end
		`,
		userID,
		subscription.ID,
		productID,
		subscription.Status,
		periodStart,
		periodEnd,
	)

	if err != nil {
		return err
	}

	switch subscription.Status {
	case "active", "trialing":
		_, err = tx.Exec(
			ctx,
			`
			UPDATE users
			SET is_premium = true
			WHERE id = $1
			`,
			userID,
		)

	default:
		_, err = tx.Exec(
			ctx,
			`
			UPDATE users
			SET is_premium = false
			WHERE id = $1
			`,
			userID,
		)
	}

	return err
}
