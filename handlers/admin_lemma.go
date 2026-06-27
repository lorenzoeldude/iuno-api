package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"


	"iuno-api/db"
	"iuno-api/models"
	"iuno-api/utils"
	"iuno-api/services"
)

type AdminLemmaRequest struct {
	Lemma     models.Lemma `json:"lemma"`
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

	if body.Lemma.Lemma == "" {
		http.Error(w, "lemma is required", http.StatusBadRequest)
		return
	}

	// =====================================================
	// FALLBACK LemmaNormalized
	// =====================================================
	if body.Lemma.LemmaNormalized == "" {
		body.Lemma.LemmaNormalized = body.Lemma.Lemma
	}

	// =====================================================
	// UPSERT LEMMA
	// =====================================================
	var id int

	err = db.Pool.QueryRow(context.Background(), `
		INSERT INTO lemmas (
			lemma,
			lemma_normalized,
			part_of_speech,
			gender,
			declension,
			conjugation,
			perfect,
			supine,
			irregular,
			genitive
		)
		VALUES (
			$1,$2,$3,$4,$5,
			$6,$7,$8,$9,$10,$11
		)
		ON CONFLICT (lemma)
		DO UPDATE SET
			lemma_normalized = EXCLUDED.lemma_normalized,
			type = EXCLUDED.type,
			gender = EXCLUDED.gender,
			declension = EXCLUDED.declension,
			conjugation = EXCLUDED.conjugation,
			perfect = EXCLUDED.perfect,
			supine = EXCLUDED.supine,
			irregular = EXCLUDED.irregular,
			genitive = EXCLUDED.genitive
		RETURNING id
	`,
		body.Lemma.Lemma,
		body.Lemma.LemmaNormalized,
		body.Lemma.PartOfSpeech,
		body.Lemma.Gender,
		body.Lemma.Declension,
		body.Lemma.Conjugation,
		body.Lemma.Perfect,
		body.Lemma.Supine,
		body.Lemma.Irregular,
		body.Lemma.Genitive,
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

func GetLemmaByIDHandler(w http.ResponseWriter, r *http.Request) {

	// URL:
	// /api/admin/lemma/123

	idString := strings.TrimPrefix(
		r.URL.Path,
		"/api/admin/lemma/",
	)


	id, err := strconv.Atoi(idString)

	if err != nil {

		http.Error(
			w,
			"invalid lemma id",
			http.StatusBadRequest,
		)

		return
	}


	word, err := services.GetWordByID(id)

	if err != nil {

		http.Error(
			w,
			"lemma not found",
			http.StatusNotFound,
		)

		return
	}


	w.Header().Set(
		"Content-Type",
		"application/json",
	)


	json.NewEncoder(w).Encode(word)
}