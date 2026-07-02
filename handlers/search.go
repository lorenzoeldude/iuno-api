package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"log"
	
	"iuno-api/services/morphology"
	"iuno-api/db"
)

type SearchResult struct {
	Lemma   string `json:"lemma"`
	Meaning string `json:"meaning"`
	Slug    string `json:"slug"`
}

type SearchFormResult struct {
	Form 				string `json:"form"`
	PartOfSpeech 		string `json:"part_of_speech"`
	Lemma   			string `json:"lemma"`
	Meaning   			string `json:"meaning"`
	LemmaNormalized   	string `json:"lemma_normalized"`
	GrammaticalCase 	*string `json:"grammatical_case"`
	Tense 				*string `json:"tense"`
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {

	// =====================================================
	// GET QUERY
	// =====================================================
	query := r.URL.Query().Get("q")
	query = strings.TrimSpace(strings.ToLower(query))

	// empty query → return empty list
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
		LEFT JOIN meanings m ON m.lemma_id = l.id
		WHERE LOWER(l.lemma) LIKE $1
		GROUP BY l.id
		ORDER BY l.lemma ASC
		LIMIT 10
	`, query+"%")

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("error searching form: ", err)
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

	// Get Query Word
	query := r.URL.Query().Get("q")
	query = strings.TrimSpace(strings.ToLower(query))

	// empty query → return empty list
	if query == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]SearchResult{})
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
			f.tense,
			l.lemma,
			COALESCE(
				(
					SELECT meaning
					FROM meanings
					WHERE lemma_id = l.id
					ORDER BY id
					LIMIT 1
				),
				''
			) AS meaning,
			l.lemma_normalized
		FROM forms f
		JOIN lemmas l
			ON l.id = f.lemma_id
		WHERE LOWER(f.form_normalized) LIKE LOWER($1)
		GROUP BY
			f.form,
			f.part_of_speech,
			f.grammatical_case,
			f.tense,
			l.id,
			l.lemma,
			l.lemma_normalized
		ORDER BY f.form ASC
		LIMIT 5;
	`, query+"%")

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("error searching form: ", err)
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
			&res.Tense,
			&res.Lemma,
			&res.Meaning,
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