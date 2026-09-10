package handlers

import (
	"encoding/json"
	"net/http"

	"iuno-api/db"
)

type TopWordLookup struct {
	ID          int    `json:"id"`
	Lemma       string `json:"lemma"`
	Meaning     string `json:"meaning"`
	LookupCount int    `json:"lookup_count"`
}

func TopWordLookupsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	rows, err := db.Pool.Query(
		r.Context(),
		`
		SELECT
			l.id,
			l.lemma,
			COALESCE(MIN(m.meaning), '') AS meaning,
			COUNT(*)::int AS lookup_count
		FROM word_lookup_daily w
		JOIN lemmas l
			ON l.id = w.lemma_id
		LEFT JOIN meanings m
			ON m.lemma_id = l.id
		WHERE w.lookup_date >= CURRENT_DATE - INTERVAL '6 days'
		GROUP BY
			l.id,
			l.lemma
		ORDER BY
			COUNT(*) DESC,
			l.lemma ASC
		LIMIT 5
		`,
	)

	if err != nil {
		http.Error(
			w,
			"failed to fetch top lookups",
			http.StatusInternalServerError,
		)
		return
	}

	defer rows.Close()

	results := make(
		[]TopWordLookup,
		0,
	)

	for rows.Next() {

		var result TopWordLookup

		err := rows.Scan(
			&result.ID,
			&result.Lemma,
			&result.Meaning,
			&result.LookupCount,
		)

		if err != nil {
			http.Error(
				w,
				"failed to read top lookups",
				http.StatusInternalServerError,
			)
			return
		}

		results = append(
			results,
			result,
		)
	}

	if err := rows.Err(); err != nil {
		http.Error(
			w,
			"failed to read top lookups",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(results)
}
