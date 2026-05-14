package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"iuno-api/db"
	"iuno-api/utils"
)

type SearchResult struct {
	Latin       string `json:"latin"`
	Translation string `json:"translation"`
	Slug        string `json:"slug"`
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {

	utils.EnableCORS(w)

	query := r.URL.Query().Get("q")
	query = strings.TrimSpace(strings.ToLower(query))

	if query == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]SearchResult{})
		return
	}

	rows, err := db.Pool.Query(r.Context(), `
		SELECT 
			latin, 
			translation_1, 
			slug
		FROM words
		WHERE LOWER(latin) LIKE $1
		ORDER BY latin
		LIMIT 10
	`, query+"%")

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	results := []SearchResult{}

	for rows.Next() {

		var res SearchResult

		err := rows.Scan(
			&res.Latin,
			&res.Translation,
			&res.Slug,
		)

		if err != nil {
			continue
		}

		results = append(results, res)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}