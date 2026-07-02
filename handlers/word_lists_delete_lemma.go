package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"iuno-api/db"
	"iuno-api/middleware"
	"iuno-api/utils"
)

func DeleteLemmaFromUserListHandler(w http.ResponseWriter, r *http.Request) {

	// =====================================================
	// METHOD CHECK
	// =====================================================
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// =====================================================
	// AUTH
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

	// =====================================================
	// GET lemma id from URL
	// /api/word-lists/lemma/:id
	// =====================================================
	lemmaIDStr := r.URL.Path[len("/api/word-lists/lemma/"):]

	lemmaID, err := strconv.Atoi(lemmaIDStr)
	if err != nil || lemmaID <= 0 {
		http.Error(w, "invalid lemma id", http.StatusBadRequest)
		return
	}

	log.Println("DELETE LEMMA:", lemmaID, "USER:", userID)

	// =====================================================
	// GET USER LIST
	// =====================================================
	var listID int
	err = db.Pool.QueryRow(r.Context(), `
		SELECT id
		FROM word_lists
		WHERE user_id = $1
		LIMIT 1
	`, userID).Scan(&listID)

	if err != nil {
		log.Println("LIST NOT FOUND:", err)
		http.Error(w, "word list not found", http.StatusBadRequest)
		return
	}

	// =====================================================
	// DELETE LEMMA
	// =====================================================
	res, err := db.Pool.Exec(r.Context(), `
		DELETE FROM word_list_lemmas
		WHERE list_id = $1 AND lemma_id = $2
	`, listID, lemmaID)

	if err != nil {
		log.Println("DELETE ERROR:", err)
		http.Error(w, "failed to delete lemma", http.StatusInternalServerError)
		return
	}

	rows := res.RowsAffected()

	// =====================================================
	// RESPONSE
	// =====================================================
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"removed": rows > 0,
	})
}