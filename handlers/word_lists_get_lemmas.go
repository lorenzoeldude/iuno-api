package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"iuno-api/db"
	"iuno-api/models"
	"iuno-api/utils"
)

func GetWordListLemmasHandler(w http.ResponseWriter, r *http.Request) {

	utils.EnableCORS(w)

	// =====================================================
	// METHOD CHECK
	// =====================================================
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// =====================================================
	// PARSE LIST ID
	// =====================================================
	listIDParam := r.URL.Query().Get("list_id")

	listID, err := strconv.Atoi(listIDParam)
	if err != nil || listID <= 0 {

		http.Error(
			w,
			"invalid list_id",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================================
	// QUERY LEMMAS
	// =====================================================
	rows, err := db.Pool.Query(r.Context(), `
		SELECT
			l.id,
			l.slug,
			l.lemma,
			l.lemma_display,
			l.type,
			l.definition,
			l.gender,
			l.declension,
			l.conjugation,
			l.stem,
			l.perfect,
			l.supine,
			l.is_irregular
		FROM word_list_lemmas wll
		JOIN lemmas l
			ON l.id = wll.lemma_id
		WHERE wll.list_id = $1
		ORDER BY l.lemma
	`, listID)

	if err != nil {

		log.Println("GET LIST LEMMAS ERROR:", err)

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
	lemmas := []models.Word{}

	for rows.Next() {

		var word models.Word

		err := rows.Scan(
			&word.ID,
			&word.Slug,
			&word.Lemma,
			&word.LemmaDisplay,
			&word.Type,
			&word.Definition,
			&word.Gender,
			&word.Declension,
			&word.Conjugation,
			&word.Stem,
			&word.Perfect,
			&word.Supine,
			&word.Irregular,
		)

		if err != nil {
			continue
		}

		lemmas = append(lemmas, word)
	}

	// =====================================================
	// RESPONSE
	// =====================================================
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(lemmas)
}