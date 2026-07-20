package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"iuno-api/db"
	"iuno-api/middleware"
	"iuno-api/utils"
)

func ListTrainerHandler(w http.ResponseWriter, r *http.Request) {


	if r.Method != http.MethodGet {

		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}



	claimsRaw := r.Context().Value(middleware.UserContextKey)

	claims, ok := claimsRaw.(*utils.Claims)

	if !ok || claims == nil {

		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)

		return
	}



	userID := claims.UserID

	listID := r.URL.Query().Get("list_id")


	// =====================================================
	// 1. PICK RANDOM LEMMA FROM USER LIST
	// =====================================================

	var lemmaID int
	var lemma string
	var lemmaNormalized string
	var infinitive *string
	var partOfSpeech string


	var err error

	if listID != "" {

		err = db.Pool.QueryRow(r.Context(), `
			SELECT
				l.id,
				l.lemma,
				l.lemma_normalized,
				l.infinitive,
				l.part_of_speech
			FROM word_list_lemmas wll
			JOIN word_lists wl
				ON wl.id = wll.list_id
			JOIN lemmas l
				ON l.id = wll.lemma_id
			WHERE wl.id = $1
			ORDER BY RANDOM()
			LIMIT 1
		`,
			listID,
		).Scan(
			&lemmaID,
			&lemma,
			&lemmaNormalized,
			&infinitive,
			&partOfSpeech,
		)

	} else {

		err = db.Pool.QueryRow(r.Context(), `
			SELECT
				l.id,
				l.lemma,
				l.lemma_normalized,
				l.infinitive,
				l.part_of_speech
			FROM word_list_lemmas wll
			JOIN word_lists wl
				ON wl.id = wll.list_id
			JOIN lemmas l
				ON l.id = wll.lemma_id
			WHERE wl.user_id = $1
			ORDER BY RANDOM()
			LIMIT 1
		`,
			userID,
		).Scan(
			&lemmaID,
			&lemma,
			&lemmaNormalized,
			&infinitive,
			&partOfSpeech,
		)

	}

	if err != nil {

		log.Println("NO WORDS ERROR:", err)

		http.Error(
			w,
			"no words in list",
			http.StatusBadRequest,
		)

		return
	}

	question, err := BuildTrainerQuestion(
		r.Context(),
		lemmaID,
		lemma,
		lemmaNormalized,
		infinitive,
		partOfSpeech,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(question)
}