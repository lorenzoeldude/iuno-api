package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"iuno-api/db"
	"iuno-api/middleware"
	"iuno-api/models"
	"iuno-api/utils"
)

func CreateWordListHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	utils.EnableCORS(w)

	// =====================================================
	// CORS PRE-FLIGHT
	// =====================================================
	if r.Method == http.MethodOptions {

		w.WriteHeader(http.StatusOK)
		return
	}

	// =====================================================
	// METHOD CHECK
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
	// AUTH USER
	// =====================================================
	claims := r.Context().Value(
		middleware.UserContextKey,
	).(*utils.Claims)

	userID := claims.UserID

	// =====================================================
	// PARSE REQUEST
	// =====================================================
	var req models.CreateWordListRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {

		http.Error(
			w,
			"invalid json",
			http.StatusBadRequest,
		)

		return
	}

	if req.Name == "" {

		http.Error(
			w,
			"name is required",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================================
	// INSERT WORD LIST
	// =====================================================
	var listID int

	err = db.Pool.QueryRow(context.Background(), `
		INSERT INTO word_lists (
			user_id,
			name
		)
		VALUES ($1, $2)
		RETURNING id
	`,
		userID,
		req.Name,
	).Scan(&listID)

	if err != nil {

		http.Error(
			w,
			"failed to create word list",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// RESPONSE
	// =====================================================
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"id":     listID,
	})
}