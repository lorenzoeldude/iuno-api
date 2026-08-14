package handlers

import (
	"encoding/json"
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
	// DECODE APPLE JWS
	// =====================================================
	//
	// Apple JWS handling lives in:
	//
	// utils/apple_jws.go
	//
	// =====================================================

	payload, err := utils.DecodeAppleTransaction(
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
		"Transaction Reason: %s",
		payload.TransactionReason,
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

	if payload.Price != nil {

		log.Printf(
			"Apple Price: %d milliunits",
			*payload.Price,
		)

	}

	if payload.Currency != nil {

		log.Printf(
			"Currency: %s",
			*payload.Currency,
		)

	}

	log.Printf(
		"Storefront: %s",
		payload.Storefront,
	)

	log.Printf(
		"Storefront ID: %s",
		payload.StorefrontID,
	)

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
	// CHECK ENVIRONMENT
	// =====================================================

	switch payload.Environment {

	case "Xcode":

		log.Printf(
			"Apple StoreKit: Xcode transaction",
		)

	case "Sandbox":

		log.Printf(
			"Apple StoreKit: Sandbox transaction",
		)

	case "Production":

		log.Printf(
			"Apple StoreKit: Production transaction",
		)

	default:

		log.Printf(
			"Apple StoreKit: unknown environment: %s",
			payload.Environment,
		)
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
			user_id =
				EXCLUDED.user_id,

			product_id =
				EXCLUDED.product_id,

			status =
				EXCLUDED.status,

			current_period_start =
				EXCLUDED.current_period_start,

			current_period_end =
				EXCLUDED.current_period_end,

			cancel_at_period_end =
				EXCLUDED.cancel_at_period_end,

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
	// PAYMENT AMOUNT
	// =====================================================
	//
	// Apple provides price in milliunits.
	//
	// Example:
	//
	// 4990 milliunits = $4.99
	//
	// Our payment_transactions.amount is stored in cents:
	//
	// 4990 / 10 = 499 cents
	//
	// =====================================================

	amount := 0
	currency := "USD"

	if payload.Price != nil {

		amount = int(
			*payload.Price / 10,
		)
	}

	if payload.Currency != nil &&
		*payload.Currency != "" {

		currency = *payload.Currency
	}

	log.Printf(
		"Apple StoreKit: payment amount = %d %s cents",
		amount,
		currency,
	)

	// =====================================================
	// SAVE PAYMENT TRANSACTION
	// =====================================================
	//
	// Xcode StoreKit transactions use:
	//
	// transactionId = "0"
	//
	// These are local test transactions and should NOT
	// be treated as real payment transactions.
	//
	// Sandbox and Production transactions have real
	// transaction IDs and are persisted.
	//
	// =====================================================

	if payload.TransactionID != "" &&
		payload.TransactionID != "0" {

		_, err = tx.Exec(
			r.Context(),
			`
			INSERT INTO payment_transactions (
				user_id,
				provider,
				provider_transaction_id,
				product_id,
				amount,
				currency,
				status
			)
			VALUES (
				$1,
				'apple',
				$2,
				$3,
				$4,
				$5,
				$6
			)
			ON CONFLICT (
				provider,
				provider_transaction_id
			)
			DO UPDATE SET
				user_id =
					EXCLUDED.user_id,

				product_id =
					EXCLUDED.product_id,

				amount =
					EXCLUDED.amount,

				currency =
					EXCLUDED.currency,

				status =
					EXCLUDED.status
			`,
			claims.UserID,
			payload.TransactionID,
			payload.ProductID,
			amount,
			currency,
			appleSubscriptionStatus(
				isPremium,
			),
		)

		if err != nil {

			log.Printf(
				"Apple StoreKit: failed to save payment transaction: %v",
				err,
			)

			http.Error(
				w,
				"failed to save payment transaction",
				http.StatusInternalServerError,
			)

			return
		}

		log.Printf(
			"Apple StoreKit: payment transaction saved for user %d, transaction %s",
			claims.UserID,
			payload.TransactionID,
		)

	} else {

		log.Printf(
			"Apple StoreKit: Xcode transaction detected - payment transaction not persisted",
		)
	}

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
