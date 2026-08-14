package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"iuno-api/db"
	"iuno-api/middleware"
	"iuno-api/utils"
)

type AppleStoreKitTransactionRequest struct {
	SignedTransaction string `json:"signed_transaction"`
}

// =====================================================
// APPLE TRANSACTION PAYLOAD
// =====================================================
//
// This represents the decoded payload inside Apple's
// signed transaction JWS.
//
// IMPORTANT:
// The JWS is decoded here for development/testing.
// Before using this in production, we should also
// cryptographically verify the JWS against Apple's
// certificate chain / public keys.
// =====================================================

type AppleTransactionPayload struct {
	TransactionID         string `json:"transactionId"`
	OriginalTransactionID string `json:"originalTransactionId"`

	ProductID string `json:"productId"`

	PurchaseDate int64 `json:"purchaseDate"`
	ExpiresDate  int64 `json:"expiresDate"`

	RevocationDate *int64 `json:"revocationDate"`

	Type string `json:"type"`

	InAppOwnershipType string `json:"inAppOwnershipType"`

	Environment string `json:"environment"`
}

// =====================================================
// HANDLER
// =====================================================

func AppleStoreKitTransactionHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	log.Printf(
		"APPLE STOREKIT REQUEST: %s %s",
		r.Method,
		r.URL.Path,
	)

	// =====================================================
	// METHOD
	// =====================================================

	if r.Method != http.MethodPost {

		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	// =====================================================
	// AUTHENTICATED USER
	// =====================================================

	claims, ok := r.Context().Value(
		middleware.UserContextKey,
	).(*utils.Claims)

	if !ok || claims == nil {

		log.Printf(
			"Apple StoreKit: unauthorized request",
		)

		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)

		return
	}

	log.Printf(
		"Apple StoreKit: authenticated user %d",
		claims.UserID,
	)

	// =====================================================
	// PARSE REQUEST
	// =====================================================

	var req AppleStoreKitTransactionRequest

	if err := json.NewDecoder(
		r.Body,
	).Decode(&req); err != nil {

		log.Printf(
			"Apple StoreKit: failed to decode request: %v",
			err,
		)

		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)

		return
	}

	if req.SignedTransaction == "" {

		log.Printf(
			"Apple StoreKit: signed_transaction missing",
		)

		http.Error(
			w,
			"signed_transaction is required",
			http.StatusBadRequest,
		)

		return
	}

	log.Printf(
		"Apple StoreKit: signed transaction length = %d",
		len(req.SignedTransaction),
	)

	// =====================================================
	// DECODE JWS PAYLOAD
	// =====================================================

	payload, err := decodeAppleTransaction(
		req.SignedTransaction,
	)

	if err != nil {

		log.Printf(
			"Apple StoreKit: failed to decode JWS: %v",
			err,
		)

		http.Error(
			w,
			"invalid signed transaction",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================================
	// LOG TRANSACTION
	// =====================================================

	log.Printf(
		"========== APPLE TRANSACTION ==========",
	)

	log.Printf(
		"User ID: %d",
		claims.UserID,
	)

	log.Printf(
		"Transaction ID: %s",
		payload.TransactionID,
	)

	log.Printf(
		"Original Transaction ID: %s",
		payload.OriginalTransactionID,
	)

	log.Printf(
		"Product ID: %s",
		payload.ProductID,
	)

	log.Printf(
		"Transaction Type: %s",
		payload.Type,
	)

	log.Printf(
		"Environment: %s",
		payload.Environment,
	)

	log.Printf(
		"Purchase Date: %s",
		formatAppleDate(payload.PurchaseDate),
	)

	log.Printf(
		"Expiration Date: %s",
		formatAppleDate(payload.ExpiresDate),
	)

	if payload.RevocationDate != nil {

		log.Printf(
			"Revocation Date: %s",
			formatAppleDate(*payload.RevocationDate),
		)

	} else {

		log.Printf(
			"Revocation Date: none",
		)
	}

	log.Printf(
		"========================================",
	)

	// =====================================================
	// VALIDATE PRODUCT
	// =====================================================

	if payload.ProductID !=
		"com.iunoni.premium.monthly" &&
		payload.ProductID !=
			"com.iunoni.premium.yearly" {

		log.Printf(
			"Apple StoreKit: unknown product ID: %s",
			payload.ProductID,
		)

		http.Error(
			w,
			"unknown product",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================================
	// CHECK EXPIRATION
	// =====================================================

	var expirationDate *time.Time

	if payload.ExpiresDate > 0 {

		t := time.UnixMilli(
			payload.ExpiresDate,
		)

		expirationDate = &t
	}

	// =====================================================
	// DETERMINE PREMIUM STATUS
	// =====================================================

	isPremium := true

	if expirationDate != nil &&
		expirationDate.Before(time.Now()) {

		isPremium = false

		log.Printf(
			"Apple StoreKit: transaction is already expired",
		)
	}

	if payload.RevocationDate != nil {

		isPremium = false

		log.Printf(
			"Apple StoreKit: transaction was revoked",
		)
	}

	// =====================================================
	// DATABASE TRANSACTION
	// =====================================================

	tx, err := db.Pool.Begin(
		r.Context(),
	)

	if err != nil {

		log.Printf(
			"Apple StoreKit: failed to begin DB transaction: %v",
			err,
		)

		http.Error(
			w,
			"failed to begin transaction",
			http.StatusInternalServerError,
		)

		return
	}

	defer tx.Rollback(
		r.Context(),
	)

	// =====================================================
	// SAVE / UPDATE SUBSCRIPTION
	// =====================================================

	_, err = tx.Exec(
		r.Context(),
		`
		INSERT INTO subscriptions (
			user_id,
			provider,
			provider_subscription_id,
			product_id,
			status,
			current_period_start,
			current_period_end,
			cancel_at_period_end
		)
		VALUES (
			$1,
			'apple',
			$2,
			$3,
			$4,
			$5,
			$6,
			false
		)
		ON CONFLICT (
			provider,
			provider_subscription_id
		)
		DO UPDATE SET
			product_id =
				EXCLUDED.product_id,

			status =
				EXCLUDED.status,

			current_period_start =
				EXCLUDED.current_period_start,

			current_period_end =
				EXCLUDED.current_period_end,

			updated_at =
				now()
		`,
		claims.UserID,
		payload.OriginalTransactionID,
		payload.ProductID,
		appleSubscriptionStatus(
			isPremium,
		),
		applePurchaseDate(
			payload.PurchaseDate,
		),
		expirationDate,
	)

	if err != nil {

		log.Printf(
			"Apple StoreKit: failed to save subscription: %v",
			err,
		)

		http.Error(
			w,
			"failed to save subscription",
			http.StatusInternalServerError,
		)

		return
	}

	log.Printf(
		"Apple StoreKit: subscription saved for user %d",
		claims.UserID,
	)

	// =====================================================
	// UPDATE USER PREMIUM
	// =====================================================

	_, err = tx.Exec(
		r.Context(),
		`
		UPDATE users
		SET is_premium = $1
		WHERE id = $2
		`,
		isPremium,
		claims.UserID,
	)

	if err != nil {

		log.Printf(
			"Apple StoreKit: failed to update premium status: %v",
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
		"Apple StoreKit: Premium status for user %d = %v",
		claims.UserID,
		isPremium,
	)

	// =====================================================
	// COMMIT
	// =====================================================

	if err := tx.Commit(
		r.Context(),
	); err != nil {

		log.Printf(
			"Apple StoreKit: failed to commit transaction: %v",
			err,
		)

		http.Error(
			w,
			"failed to commit transaction",
			http.StatusInternalServerError,
		)

		return
	}

	log.Printf(
		"Apple StoreKit: transaction processed successfully",
	)

	// =====================================================
	// RESPONSE
	// =====================================================

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(
		http.StatusOK,
	)

	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"success":     true,
			"user_id":     claims.UserID,
			"product_id":  payload.ProductID,
			"is_premium":  isPremium,
			"environment": payload.Environment,
		},
	)
}

// =====================================================
// DECODE APPLE JWS
// =====================================================

func decodeAppleTransaction(
	jws string,
) (*AppleTransactionPayload, error) {

	parts := splitJWS(jws)

	if len(parts) != 3 {

		return nil, fmt.Errorf(
			"invalid JWS: expected 3 parts, got %d",
			len(parts),
		)
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(
		parts[1],
	)

	if err != nil {

		return nil, fmt.Errorf(
			"failed to decode JWS payload: %w",
			err,
		)
	}

	var payload AppleTransactionPayload

	if err := json.Unmarshal(
		payloadBytes,
		&payload,
	); err != nil {

		return nil, fmt.Errorf(
			"failed to parse Apple transaction payload: %w",
			err,
		)
	}

	if payload.TransactionID == "" {

		return nil, fmt.Errorf(
			"transactionId missing",
		)
	}

	if payload.OriginalTransactionID == "" {

		return nil, fmt.Errorf(
			"originalTransactionId missing",
		)
	}

	if payload.ProductID == "" {

		return nil, fmt.Errorf(
			"productId missing",
		)
	}

	return &payload, nil
}

// =====================================================
// SPLIT JWS
// =====================================================

func splitJWS(
	jws string,
) []string {

	var parts []string

	start := 0

	for i := 0; i < len(jws); i++ {

		if jws[i] == '.' {

			parts = append(
				parts,
				jws[start:i],
			)

			start = i + 1
		}
	}

	parts = append(
		parts,
		jws[start:],
	)

	return parts
}

// =====================================================
// APPLE PURCHASE DATE
// =====================================================

func applePurchaseDate(
	milliseconds int64,
) *time.Time {

	if milliseconds <= 0 {
		return nil
	}

	t := time.UnixMilli(
		milliseconds,
	)

	return &t
}

// =====================================================
// APPLE SUBSCRIPTION STATUS
// =====================================================

func appleSubscriptionStatus(
	isPremium bool,
) string {

	if isPremium {
		return "active"
	}

	return "expired"
}

// =====================================================
// FORMAT APPLE DATE
// =====================================================

func formatAppleDate(
	milliseconds int64,
) string {

	if milliseconds <= 0 {
		return "unknown"
	}

	return time.UnixMilli(
		milliseconds,
	).UTC().Format(
		time.RFC3339,
	)
}