package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"iuno-api/middleware"
	"iuno-api/stripe"
	"iuno-api/utils"

	stripeSDK "github.com/stripe/stripe-go/v86"
)

type CreateCheckoutSessionRequest struct {
	PriceID string `json:"price_id"`
}

func CreateCheckoutSessionHandler(w http.ResponseWriter, r *http.Request) {

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
	// PARSE REQUEST
	// =====================================================

	var req CreateCheckoutSessionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if req.PriceID == "" {
		http.Error(
			w,
			"price_id is required",
			http.StatusBadRequest,
		)
		return
	}

	// =====================================================
	// CHECKOUT URLS
	// =====================================================

	successURL := os.Getenv("STRIPE_SUCCESS_URL")
	cancelURL := os.Getenv("STRIPE_CANCEL_URL")

	if successURL == "" {
		successURL = "https://iunoni.com/payment/success"
	}

	if cancelURL == "" {
		cancelURL = "https://iunoni.com/payment/cancel"
	}

	// =====================================================
	// CREATE CHECKOUT SESSION
	// =====================================================

	params := &stripeSDK.CheckoutSessionCreateParams{
		Mode: stripeSDK.String(
			string(stripeSDK.CheckoutSessionModeSubscription),
		),

		SuccessURL: stripeSDK.String(successURL),
		CancelURL:  stripeSDK.String(cancelURL),
	}

	// =====================================================
	// LINE ITEMS
	// =====================================================

	params.LineItems = []*stripeSDK.CheckoutSessionCreateLineItemParams{
		{
			Price:    stripeSDK.String(req.PriceID),
			Quantity: stripeSDK.Int64(1),
		},
	}

	// =====================================================
	// ATTACH IUNO DATA TO STRIPE SESSION
	// =====================================================

	params.Metadata = map[string]string{
		"user_id":  strconv.Itoa(claims.UserID),
		"price_id": req.PriceID,
	}

	// =====================================================
	// CREATE SESSION
	// =====================================================

	session, err := stripe.Client.V1CheckoutSessions.Create(
		context.Background(),
		params,
	)

	if err != nil {
		http.Error(
			w,
			"failed to create checkout session",
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

	json.NewEncoder(w).Encode(map[string]string{
		"checkout_url": session.URL,
		"session_id":   session.ID,
	})
}
