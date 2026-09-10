package handlers

import (
	"encoding/json"
	"net/http"

	"iuno-api/db"
	"iuno-api/middleware"
	"iuno-api/utils"
)

type WordLookupRequest struct {
	LemmaID int `json:"lemmaId"`
}

func RecordWordLookupHandler(w http.ResponseWriter, r *http.Request) {

	var req WordLookupRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.LemmaID <= 0 {
		http.Error(w, "invalid lemma id", http.StatusBadRequest)
		return
	}

	// =====================================================
	// LOGGED-IN USER
	// =====================================================

	claims, ok := r.Context().Value(
		middleware.UserContextKey,
	).(*utils.Claims)

	if ok && claims != nil {

		_, err := db.Pool.Exec(
			r.Context(),
			`INSERT INTO word_lookup_daily (
				user_id,
				lemma_id,
				lookup_date
			)
			VALUES ($1, $2, CURRENT_DATE)
			ON CONFLICT DO NOTHING`,
			claims.UserID,
			req.LemmaID,
		)

		if err != nil {
			http.Error(
				w,
				"failed to record word lookup",
				http.StatusInternalServerError,
			)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	// =====================================================
	// ANONYMOUS USER
	// =====================================================

	anonymousID, ok := r.Context().Value(
		middleware.AnonymousLookupIDKey,
	).(string)

	if !ok || anonymousID == "" {
		http.Error(
			w,
			"missing anonymous id",
			http.StatusBadRequest,
		)
		return
	}

	_, err := db.Pool.Exec(
		r.Context(),
		`INSERT INTO word_lookup_daily (
			anonymous_id,
			lemma_id,
			lookup_date
		)
		VALUES ($1::uuid, $2, CURRENT_DATE)
		ON CONFLICT DO NOTHING`,
		anonymousID,
		req.LemmaID,
	)

	if err != nil {
		http.Error(
			w,
			"failed to record word lookup",
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
