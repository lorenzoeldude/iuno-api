package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"iuno-api/db"
	"iuno-api/middleware"
	"iuno-api/models"
	"iuno-api/utils"
)

func AddLemmaToUserListHandler(w http.ResponseWriter, r *http.Request) {

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
	if !ok || claims == nil {
		http.Error(w, "invalid auth context", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID
	log.Println("ADD LEMMA REQUEST - USER:", userID)

	// =====================================================
	// PARSE REQUEST
	// =====================================================
	var req models.AddLemmaToListRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
	// GET USER DEFAULT WORD LIST
	// =====================================================
	var listID int

	err := db.Pool.QueryRow(r.Context(), `
		SELECT id
		FROM word_lists
		WHERE user_id = $1
		LIMIT 1
	`, userID).Scan(&listID)

	if err != nil {
		log.Println("WORD LIST NOT FOUND:", err)
		http.Error(w, "word list not found", http.StatusBadRequest)
		return
	}

	log.Println("LIST ID:", listID)

	// =====================================================
	// INSERT LEMMA (SAFE)
	// =====================================================
	res, err := db.Pool.Exec(r.Context(), `
		INSERT INTO word_list_lemmas (
			list_id,
			lemma_id
		)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, listID, req.LemmaID)

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
		"added":  rows > 0,
	})
}