package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"iuno-api/middleware"
	"iuno-api/utils"
)

type AppleStoreKitTransactionRequest struct {
	SignedTransaction string `json:"signed_transaction"`
}

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

	var req AppleStoreKitTransactionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if req.SignedTransaction == "" {
		http.Error(
			w,
			"signed_transaction is required",
			http.StatusBadRequest,
		)
		return
	}

	// =====================================================
	// TEMPORARY LOGGING
	// =====================================================

	log.Printf(
		"Apple StoreKit transaction received for user %d",
		claims.UserID,
	)

	log.Printf(
		"Apple signed transaction length: %d",
		len(req.SignedTransaction),
	)

	// =====================================================
	// RESPONSE
	// =====================================================

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"success": true,
			"user_id": claims.UserID,
		},
	)
}