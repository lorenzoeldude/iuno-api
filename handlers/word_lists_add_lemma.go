package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"iuno-api/db"
	"iuno-api/models"
	"iuno-api/utils"
)

func AddLemmaToListHandler(w http.ResponseWriter, r *http.Request) {

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
	var body models.AddLemmaToListRequest

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
	if body.ListID == 0 || body.LemmaID == 0 {

		http.Error(
			w,
			"list_id and lemma_id are required",
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
		ON CONFLICT (list_id, lemma_id)
		DO NOTHING
	`,
		body.ListID,
		body.LemmaID,
	)

	if err != nil {

		log.Println("ADD LEMMA ERROR:", err)

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
		"status":   "ok",
		"list_id":  body.ListID,
		"lemma_id": body.LemmaID,
	})
}