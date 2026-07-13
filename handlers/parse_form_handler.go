package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"iuno-api/db"
	"iuno-api/services/morphology"
)

func ParseFormHandler(w http.ResponseWriter, r *http.Request) {

	// =====================================================
	// GET QUERY
	// =====================================================

	query := r.URL.Query().Get("form")
	query = strings.TrimSpace(strings.ToLower(query))

	if query == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]SearchFormResult{})
		return
	}

	query = morphology.NormalizeLatin(query)

	// =====================================================
	// DB QUERY
	// =====================================================

	rows, err := db.Pool.Query(r.Context(), `
		SELECT
			f.form,
			f.part_of_speech,
			f.grammatical_case,
			f.number,
			f.gender,
			f.tense,
			f.mood,
			f.voice,
			f.person,
			l.lemma,
			COALESCE(
				(
					SELECT array_agg(meaning ORDER BY id)
					FROM (
						SELECT meaning, id
						FROM meanings
						WHERE lemma_id = l.id
						ORDER BY id
						LIMIT 3
					) m
				),
				ARRAY[]::text[]
			) AS meanings,
			l.lemma_normalized
		FROM forms f
		JOIN lemmas l
			ON l.id = f.lemma_id
		WHERE f.form_normalized = $1
		GROUP BY
			f.form,
			f.part_of_speech,
			f.grammatical_case,
			f.number,
			f.gender,
			f.tense,
			f.mood,
			f.voice,
			f.person,
			l.id,
			l.lemma,
			l.lemma_normalized
		ORDER BY
			l.lemma,
			f.part_of_speech,
			f.tense NULLS FIRST,
			f.grammatical_case NULLS FIRST;
	`, query)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("error parsing form:", err)
		return
	}
	defer rows.Close()

	// =====================================================
	// BUILD RESULTS
	// =====================================================

	results := []SearchFormResult{}

	for rows.Next() {

		var res SearchFormResult

		err := rows.Scan(
			&res.Form,
			&res.PartOfSpeech,
			&res.GrammaticalCase,
			&res.Number,
			&res.Gender,
			&res.Tense,
			&res.Mood,
			&res.Voice,
			&res.Person,
			&res.Lemma,
			&res.Meanings,
			&res.LemmaNormalized,
		)

		if err != nil {
			log.Println("scan error:", err)
			continue
		}

		results = append(results, res)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("rows error:", err)
		return
	}

	// =====================================================
	// RESPONSE
	// =====================================================

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}