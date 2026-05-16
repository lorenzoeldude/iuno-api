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

func AddLemmaToListHandler(
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

	_ = userID // (kept for future ownership validation if needed)

	// =====================================================
	// PARSE REQUEST
	// =====================================================
	var req models.AddLemmaToListRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {

		http.Error(
			w,
			"invalid json",
			http.StatusBadRequest,
		)

		return
	}

	if req.ListID == 0 || req.LemmaID == 0 {

		http.Error(
			w,
			"list_id and lemma_id required",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================================
	// INSERT RELATION
	// =====================================================
	_, err = db.Pool.Exec(context.Background(), `
		INSERT INTO word_list_lemmas (
			list_id,
			lemma_id
		)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`,
		req.ListID,
		req.LemmaID,
	)

	if err != nil {

		http.Error(
			w,
			"failed to add lemma to list",
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
	})
}