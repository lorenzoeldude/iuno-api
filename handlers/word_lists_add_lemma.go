package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"iuno-api/db"
	"iuno-api/middleware"
	"iuno-api/models"
	"iuno-api/utils"
)

func AddLemmaToUserListHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("hello there")

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
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// =====================================================
	// AUTH USER
	// =====================================================
	claimsRaw := r.Context().Value(middleware.UserContextKey)
	if claimsRaw == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	claims, ok := claimsRaw.(*utils.Claims)
	if !ok {
		http.Error(w, "invalid auth context", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID

	log.Println("ADD LEMMA REQUEST - USER:", userID)

	// =====================================================
	// PARSE REQUEST
	// =====================================================
	var req models.AddLemmaToListRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("JSON ERROR:", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.LemmaID == 0 {
		http.Error(w, "lemma_id required", http.StatusBadRequest)
		return
	}

	log.Println("LEMMA ID:", req.LemmaID)

	// =====================================================
	// GET USER WORD LIST
	// =====================================================
	var listID int

	err = db.Pool.QueryRow(context.Background(), `
		SELECT id
		FROM word_lists
		WHERE user_id = $1
		LIMIT 1
	`,
		userID,
	).Scan(&listID)

	if err != nil {
		log.Println("WORD LIST NOT FOUND:", err)
		http.Error(w, "word list not found", http.StatusBadRequest)
		return
	}

	log.Println("LIST ID:", listID)

	// =====================================================
	// INSERT LEMMA
	// =====================================================
	res, err := db.Pool.Exec(context.Background(), `
		INSERT INTO word_list_lemmas (
			list_id,
			lemma_id
		)
		VALUES ($1, $2)
	`,
		listID,
		req.LemmaID,
	)

	if err != nil {
		log.Println("DB INSERT ERROR:", err)
		http.Error(w, "failed to add lemma", http.StatusInternalServerError)
		return
	}

	rows := res.RowsAffected()
	log.Println("ROWS AFFECTED:", rows)

	// =====================================================
	// RESPONSE
	// =====================================================
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}