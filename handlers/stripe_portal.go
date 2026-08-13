package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"iuno-api/db"
	"iuno-api/middleware"
	"iuno-api/stripe"
	"iuno-api/utils"

	stripeSDK "github.com/stripe/stripe-go/v86"
)

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

	// =====================================================
	// GET AUTHENTICATED USER
	// =====================================================

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

	// =====================================================
	// GET STRIPE CUSTOMER
	// =====================================================

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

	if err != nil {
		http.Error(
			w,
			"Stripe customer not found",
			http.StatusNotFound,
		)
		return
	}

	// =====================================================
	// RETURN URL
	// =====================================================

	returnURL := os.Getenv("STRIPE_PORTAL_RETURN_URL")

	if returnURL == "" {
		returnURL = "https://iunoni.com/account"
	}

	// =====================================================
	// CREATE PORTAL SESSION
	// =====================================================

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

	// =====================================================
	// RESPONSE
	// =====================================================

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string]string{
			"url": session.URL,
		},
	)
}
