package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"iuno-api/db"
	"iuno-api/services/morphology"
)

type SearchResult struct {
	Lemma   string `json:"lemma"`
	Meaning string `json:"meaning"`
	Slug    string `json:"slug"`
}

type SearchFormResult struct {
	Form            string   `json:"form"`
	PartOfSpeech    string   `json:"part_of_speech"`
	Lemma           string   `json:"lemma"`
	Meanings        []string `json:"meanings"`
	LemmaNormalized string   `json:"lemma_normalized"`

	GrammaticalCase *string `json:"grammatical_case"`
	Number          *string `json:"number"`
	Gender          *string `json:"gender"`

	Tense  *string `json:"tense"`
	Mood   *string `json:"mood"`
	Voice  *string `json:"voice"`
	Person *int    `json:"person"`
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {

	// =====================================================
	// GET QUERY
	// =====================================================

	query := r.URL.Query().Get("q")
	query = strings.TrimSpace(strings.ToLower(query))

	if query == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]SearchResult{})
		return
	}

	// =====================================================
	// DB QUERY
	// =====================================================

	rows, err := db.Pool.Query(r.Context(), `
		SELECT
			l.lemma,
			COALESCE(MIN(m.meaning), '') AS meaning,
			l.slug
		FROM lemmas l
		LEFT JOIN meanings m
			ON m.lemma_id = l.id
		WHERE
			LOWER(l.lemma) LIKE $1
			OR EXISTS (
				SELECT 1
				FROM meanings m2
				WHERE
					m2.lemma_id = l.id
					AND LOWER(m2.meaning) LIKE '%' || $1 || '%'
			)
		GROUP BY l.id
		ORDER BY
			CASE
				WHEN LOWER(l.lemma) LIKE $1 THEN 0
				ELSE 1
			END,
			l.lemma ASC
		LIMIT 10
	`, query+"%")

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("error searching lemma:", err)
		return
	}
	defer rows.Close()

	// =====================================================
	// BUILD RESULTS
	// =====================================================

	results := []SearchResult{}

	for rows.Next() {

		var res SearchResult

		err := rows.Scan(
			&res.Lemma,
			&res.Meaning,
			&res.Slug,
		)

		if err != nil {
			continue
		}

		results = append(results, res)
	}

	// =====================================================
	// RESPONSE
	// =====================================================

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func SearchFormHandler(w http.ResponseWriter, r *http.Request) {

	// =====================================================
	// GET QUERY
	// =====================================================

	query := r.URL.Query().Get("q")
	query = strings.TrimSpace(strings.ToLower(query))

	if query == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]SearchFormResult{})
		return
	}

	normalizedQuery := morphology.NormalizeLatin(query)

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
		WHERE
			LOWER(f.form_normalized) LIKE LOWER($1)
			OR EXISTS (
				SELECT 1
				FROM meanings m2
				WHERE
					m2.lemma_id = l.id
					AND LOWER(m2.meaning) LIKE '%' || LOWER($2) || '%'
			)
		GROUP BY
			f.form,
			f.form_normalized,
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
			CASE
				WHEN EXISTS (
					SELECT 1
					FROM meanings m3
					WHERE
						m3.lemma_id = l.id
						AND LOWER(TRIM(m3.meaning)) = LOWER($2)
				) THEN 0

				WHEN EXISTS (
					SELECT 1
					FROM meanings m3
					WHERE
						m3.lemma_id = l.id
						AND LOWER(m3.meaning) LIKE LOWER($2) || '%'
				) THEN 1

				WHEN EXISTS (
					SELECT 1
					FROM meanings m3
					WHERE
						m3.lemma_id = l.id
						AND LOWER(m3.meaning) LIKE '%' || LOWER($2) || '%'
				) THEN 2

				WHEN LOWER(f.form_normalized) LIKE LOWER($1) THEN 3

				ELSE 4
			END,
			f.form ASC
		LIMIT 5;
	`, normalizedQuery+"%", query)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("error searching form:", err)
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
			continue
		}

		results = append(results, res)
	}

	// =====================================================
	// RESPONSE
	// =====================================================

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
