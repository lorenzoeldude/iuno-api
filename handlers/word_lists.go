package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	// "log"
	"iuno-api/utils"

	"context"
	"iuno-api/db"
	"iuno-api/middleware"
)

func getDefaultListID(userID int) (int, error) {

	var listID int

	err := db.Pool.QueryRow(context.Background(), `
		SELECT id
		FROM word_lists
		WHERE user_id = $1
		LIMIT 1
	`, userID).Scan(&listID)

	return listID, err
}

func CheckLemmaSavedHandler(w http.ResponseWriter, r *http.Request) {
	

	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
	if !ok || claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID

	prefix := "/api/word-lists/lemma/"

	if len(r.URL.Path) <= len(prefix) {
		http.Error(w, "missing lemma id", http.StatusBadRequest)
		return
	}

	lemmaID, err := strconv.Atoi(r.URL.Path[len(prefix):])
	if err != nil {
		http.Error(w, "invalid lemma id", http.StatusBadRequest)
		return
	}

	listID, err := getDefaultListID(userID)
	if err != nil {
		http.Error(w, "no word list found", http.StatusBadRequest)
		return
	}

	var exists bool

	err = db.Pool.QueryRow(r.Context(), `
		SELECT EXISTS (
			SELECT 1
			FROM word_list_lemmas
			WHERE list_id = $1 AND lemma_id = $2
		)
	`, listID, lemmaID).Scan(&exists)

	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{
		"saved": exists,
	})
}