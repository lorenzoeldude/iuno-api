package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"iuno-api/db"
	"iuno-api/models"
	"iuno-api/utils"
)

func CreateWordListHandler(w http.ResponseWriter, r *http.Request) {

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
	// PARSE REQUEST
	// =====================================================
	var body models.CreateWordListRequest

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {

		log.Println("JSON ERROR:", err)

		http.Error(
			w,
			"invalid json",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================================
	// VALIDATION
	// =====================================================
	body.Name = strings.TrimSpace(body.Name)

	if body.Name == "" {

		http.Error(
			w,
			"list name is required",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================================
	// TEMP USER ID
	// later comes from auth middleware
	// =====================================================
	userID := 1

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
		body.Name,
	).Scan(&listID)

	if err != nil {

		log.Println("CREATE LIST ERROR:", err)

		http.Error(
			w,
			"failed to create list",
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
		"name":   body.Name,
	})
}