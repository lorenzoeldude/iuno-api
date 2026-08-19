package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"

	"iuno-api/middleware"
	"iuno-api/storekit"
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
	// VERIFY APPLE TRANSACTION
	// =====================================================

	transaction, err :=
		storekit.VerifyTransaction(
			req.SignedTransaction,
		)

	if err != nil {

		log.Printf(
			"Apple StoreKit: transaction verification failed: %v",
			err,
		)

		http.Error(
			w,
			"invalid signed transaction",
			http.StatusBadRequest,
		)

		return
	}

	decodedPayload :=
		transaction.Payload

	// =====================================================
	// VALIDATE PRODUCT
	// =====================================================

	if !storekit.IsValidProduct(
		decodedPayload.ProductID,
	) {

		log.Printf(
			"Apple StoreKit: unknown product ID: %s",
			decodedPayload.ProductID,
		)

		http.Error(
			w,
			"unknown product",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================================
	// PREMIUM STATUS
	// =====================================================

	isPremium :=
		storekit.IsPremium(
			decodedPayload,
		)

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
		decodedPayload.TransactionID,
	)

	log.Printf(
		"Original Transaction ID: %s",
		decodedPayload.OriginalTransactionID,
	)

	log.Printf(
		"Product ID: %s",
		decodedPayload.ProductID,
	)

	log.Printf(
		"Transaction Type: %s",
		decodedPayload.Type,
	)

	log.Printf(
		"Transaction Reason: %s",
		decodedPayload.TransactionReason,
	)

	log.Printf(
		"Environment: %s",
		decodedPayload.Environment,
	)

	log.Printf(
		"Purchase Date: %s",
		formatAppleDate(
			decodedPayload.PurchaseDate,
		),
	)

	log.Printf(
		"Expiration Date: %s",
		formatAppleDate(
			decodedPayload.ExpiresDate,
		),
	)

	if decodedPayload.RevocationDate != nil {

		log.Printf(
			"Revocation Date: %s",
			formatAppleDate(
				*decodedPayload.RevocationDate,
			),
		)

	} else {

		log.Printf(
			"Revocation Date: none",
		)
	}

	if decodedPayload.Price != nil {

		log.Printf(
			"Apple Price: %d milliunits",
			*decodedPayload.Price,
		)

	}

	if decodedPayload.Currency != nil {

		log.Printf(
			"Currency: %s",
			*decodedPayload.Currency,
		)

	}

	log.Printf(
		"Storefront: %s",
		decodedPayload.Storefront,
	)

	log.Printf(
		"Storefront ID: %s",
		decodedPayload.StorefrontID,
	)

	log.Printf(
		"Premium Status: %v",
		isPremium,
	)

	log.Printf(
		"========================================",
	)

	// =====================================================
	// DATABASE TRANSACTION
	// =====================================================

	tx, err :=
		storekit.BeginTransaction(
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
	// CHECK SUBSCRIPTION OWNERSHIP
	// =====================================================

	ownerID, err :=
		storekit.GetSubscriptionOwner(
			r.Context(),
			tx,
			decodedPayload.OriginalTransactionID,
		)

	if err != nil &&
		!errors.Is(err, pgx.ErrNoRows) {

		log.Printf(
			"Apple StoreKit: failed to check subscription ownership: %v",
			err,
		)

		http.Error(
			w,
			"failed to check subscription ownership",
			http.StatusInternalServerError,
		)

		return
	}

	if err == nil {

		// Subscription already exists.

		if ownerID != claims.UserID {

			log.Printf(
				"Apple StoreKit: subscription %s already belongs to user %d; user %d rejected",
				decodedPayload.OriginalTransactionID,
				ownerID,
				claims.UserID,
			)

			http.Error(
				w,
				"subscription belongs to another user",
				http.StatusForbidden,
			)

			return
		}

		log.Printf(
			"Apple StoreKit: subscription ownership confirmed for user %d",
			claims.UserID,
		)

	} else {

		// No subscription exists yet.

		log.Printf(
			"Apple StoreKit: subscription %s has no existing owner",
			decodedPayload.OriginalTransactionID,
		)
	}

	// =====================================================
	// SAVE SUBSCRIPTION
	// =====================================================

	if err :=
		storekit.SaveSubscription(
			r.Context(),
			tx,
			claims.UserID,
			transaction,
			isPremium,
		); err != nil {

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
	// SAVE PAYMENT TRANSACTION
	// =====================================================

	if err :=
		storekit.SavePaymentTransaction(
			r.Context(),
			tx,
			claims.UserID,
			transaction,
			isPremium,
		); err != nil {

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

	if storekit.IsXcodeTransaction(
		decodedPayload,
	) {

		log.Printf(
			"Apple StoreKit: Xcode transaction detected - payment transaction not persisted",
		)

	} else {

		log.Printf(
			"Apple StoreKit: payment transaction saved for user %d, transaction %s",
			claims.UserID,
			decodedPayload.TransactionID,
		)
	}

	// =====================================================
	// UPDATE USER PREMIUM
	// =====================================================

	if err :=
		storekit.UpdateUserPremium(
			r.Context(),
			tx,
			claims.UserID,
			isPremium,
		); err != nil {

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

	if err :=
		tx.Commit(
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
			"product_id":  decodedPayload.ProductID,
			"is_premium":  isPremium,
			"environment": decodedPayload.Environment,
		},
	)
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
