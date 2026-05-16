package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"iuno-api/db"
	"iuno-api/models"
	"iuno-api/utils"
)

func GetWordListsHandler(w http.ResponseWriter, r *http.Request) {

	utils.EnableCORS(w)

	// =====================================================
	// METHOD CHECK
	// =====================================================
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// =====================================================
	// TEMP USER ID
	// later comes from auth middleware
	// =====================================================
	userID := 1

	// =====================================================
	// QUERY WORD LISTS
	// =====================================================
	rows, err := db.Pool.Query(r.Context(), `
		SELECT
			id,
			user_id,
			name,
			created_at
		FROM word_lists
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)

	if err != nil {

		log.Println("GET WORD LISTS ERROR:", err)

		http.Error(
			w,
			"database error",
			http.StatusInternalServerError,
		)

		return
	}

	defer rows.Close()

	// =====================================================
	// BUILD RESPONSE
	// =====================================================
	lists := []models.WordList{}

	for rows.Next() {

		var list models.WordList

		err := rows.Scan(
			&list.ID,
			&list.UserID,
			&list.Name,
			&list.CreatedAt,
		)

		if err != nil {
			continue
		}

		lists = append(lists, list)
	}

	// =====================================================
	// RESPONSE
	// =====================================================
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(lists)
}