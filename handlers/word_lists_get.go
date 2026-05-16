package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"iuno-api/db"
	"iuno-api/middleware"
	// "iuno-api/models"
	"iuno-api/utils"
)

type WordListWithCount struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	LemmaCount int   `json:"lemma_count"`
}

func GetWordListsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

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
	if r.Method != http.MethodGet {

		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	// =====================================================
	// AUTH USER (FROM JWT)
	// =====================================================
	claims := r.Context().Value(
		middleware.UserContextKey,
	).(*utils.Claims)

	userID := claims.UserID

	// =====================================================
	// QUERY WORD LISTS
	// =====================================================
	rows, err := db.Pool.Query(context.Background(), `
		SELECT 
			wl.id,
			wl.name,
			COUNT(wll.lemma_id) AS lemma_count
		FROM word_lists wl
		LEFT JOIN word_list_lemmas wll
			ON wl.id = wll.list_id
		WHERE wl.user_id = $1
		GROUP BY wl.id
		ORDER BY wl.created_at DESC
	`,
		userID,
	)

	if err != nil {

		http.Error(
			w,
			"failed to fetch word lists",
			http.StatusInternalServerError,
		)

		return
	}
	defer rows.Close()

	lists := []WordListWithCount{}

	for rows.Next() {

		var wl WordListWithCount

		err := rows.Scan(
			&wl.ID,
			&wl.Name,
			&wl.LemmaCount,
		)

		if err != nil {
			continue
		}

		lists = append(lists, wl)
	}

	// =====================================================
	// RESPONSE
	// =====================================================
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(lists)
}