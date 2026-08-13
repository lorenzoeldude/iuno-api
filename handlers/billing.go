package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"iuno-api/db"
	"iuno-api/middleware"
	"iuno-api/utils"
)

type BillingStatusResponse struct {
	IsPremium         bool       `json:"is_premium"`
	Status            *string    `json:"status"`
	CurrentPeriodEnd  *time.Time `json:"current_period_end"`
	CancelAtPeriodEnd bool       `json:"cancel_at_period_end"`
	HasStripeCustomer bool       `json:"has_stripe_customer"`
}

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
	// GET BILLING STATUS
	// =====================================================

	var response BillingStatusResponse

	err := db.Pool.QueryRow(
		r.Context(),
		`
		SELECT
			u.is_premium,
			s.status,
			s.current_period_end,
			COALESCE(s.cancel_at_period_end, false),

			EXISTS (
				SELECT 1
				FROM billing_customers bc
				WHERE bc.user_id = u.id
				  AND bc.provider = 'stripe'
			)

		FROM users u

		LEFT JOIN subscriptions s
			ON s.user_id = u.id
			AND s.provider = 'stripe'

		WHERE u.id = $1
		`,
		claims.UserID,
	).Scan(
		&response.IsPremium,
		&response.Status,
		&response.CurrentPeriodEnd,
		&response.CancelAtPeriodEnd,
		&response.HasStripeCustomer,
	)

	if err != nil {
		// Log the actual database error so we can diagnose
		// future billing problems instead of hiding it.
		http.Error(
			w,
			"failed to get billing status: "+err.Error(),
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

	json.NewEncoder(w).Encode(response)
}
