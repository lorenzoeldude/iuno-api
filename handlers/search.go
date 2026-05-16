package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"iuno-api/db"
	"iuno-api/utils"
)

type SearchResult struct {
	Lemma   string `json:"lemma"`
	Meaning string `json:"meaning"`
	Slug    string `json:"slug"`
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {

	utils.EnableCORS(w)

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