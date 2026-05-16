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

type AdminLemmaRequest struct {
	Word     models.Word `json:"word"`
	Meanings []string    `json:"meanings"`
}

func UpsertLemmaHandler(w http.ResponseWriter, r *http.Request) {

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

	log.Println("UpsertLemmaHandler HIT")

	// =====================================================
	// PARSE REQUEST
	// =====================================================
	var body AdminLemmaRequest

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println("JSON ERROR:", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if body.Word.Lemma == "" {
		http.Error(w, "lemma is required", http.StatusBadRequest)
		return
	}

	// =====================================================
	// FALLBACK SLUG
	// =====================================================
	if body.Word.Slug == "" {
		body.Word.Slug = body.Word.Lemma
	}

	// =====================================================
	// UPSERT LEMMA
	// =====================================================
	var id int

	err = db.Pool.QueryRow(context.Background(), `
		INSERT INTO lemmas (
			lemma,
			slug,
			lemma_display,
			type,
			definition,
			gender,
			declension,
			conjugation,
			stem,
			perfect,
			supine,
			is_irregular
		)
		VALUES (
			$1,$2,$3,$4,$5,
			$6,$7,$8,$9,$10,$11,$12
		)
		ON CONFLICT (lemma)
		DO UPDATE SET
			slug = EXCLUDED.slug,
			lemma_display = EXCLUDED.lemma_display,
			type = EXCLUDED.type,
			definition = EXCLUDED.definition,
			gender = EXCLUDED.gender,
			declension = EXCLUDED.declension,
			conjugation = EXCLUDED.conjugation,
			stem = EXCLUDED.stem,
			perfect = EXCLUDED.perfect,
			supine = EXCLUDED.supine,
			is_irregular = EXCLUDED.is_irregular
		RETURNING id
	`,
		body.Word.Lemma,
		body.Word.Slug,
		body.Word.LemmaDisplay,
		body.Word.Type,
		body.Word.Definition,
		body.Word.Gender,
		body.Word.Declension,
		body.Word.Conjugation,
		body.Word.Stem,
		body.Word.Perfect,
		body.Word.Supine,
		body.Word.Irregular,
	).Scan(&id)

	if err != nil {
		log.Println("UPSERT ERROR:", err)
		http.Error(w, "failed to save lemma", http.StatusInternalServerError)
		return
	}

	// =====================================================
	// DELETE OLD MEANINGS
	// =====================================================
	_, err = db.Pool.Exec(context.Background(), `
		DELETE FROM meanings
		WHERE lemma_id = $1
	`, id)

	if err != nil {
		log.Println("MEANINGS DELETE ERROR:", err)
		http.Error(w, "failed to update meanings", http.StatusInternalServerError)
		return
	}

	// =====================================================
	// INSERT NEW MEANINGS
	// =====================================================
	for _, m := range body.Meanings {

		if m == "" {
			continue
		}

		_, err := db.Pool.Exec(context.Background(), `
			INSERT INTO meanings (lemma_id, meaning)
			VALUES ($1, $2)
		`, id, m)

		if err != nil {
			log.Println("MEANING INSERT ERROR:", err)
			http.Error(w, "failed to insert meanings", http.StatusInternalServerError)
			return
		}
	}

	// =====================================================
	// RESPONSE
	// =====================================================
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"id":     id,
	})
}