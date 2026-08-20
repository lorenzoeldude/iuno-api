package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"iuno-api/db"
)

func LessonPreviousTrainerHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	// =====================================================
	// GET LESSON ID FROM URL
	//
	// /api/lessons/{id}/trainer/previous/random
	// =====================================================

	idString := strings.TrimPrefix(
		r.URL.Path,
		"/api/lessons/",
	)

	idString = strings.TrimSuffix(
		idString,
		"/trainer/previous/random",
	)

	lessonID, err := strconv.Atoi(idString)

	if err != nil {
		http.Error(
			w,
			"invalid lesson id",
			http.StatusBadRequest,
		)
		return
	}

	// =====================================================
	// PICK RANDOM LEMMA FROM ALL PREVIOUS LESSONS
	// =====================================================

	var lemmaID int
	var lemma string
	var lemmaNormalized string
	var infinitive *string
	var partOfSpeech string

	err = db.Pool.QueryRow(
		r.Context(),
		`
		SELECT
			l.id,
			l.lemma,
			l.lemma_normalized,
			l.infinitive,
			l.part_of_speech
		FROM lesson_vocabulary lv
		JOIN lemmas l
			ON l.id = lv.lemma_id
		WHERE lv.lesson_id < $1
		ORDER BY RANDOM()
		LIMIT 1
		`,
		lessonID,
	).Scan(
		&lemmaID,
		&lemma,
		&lemmaNormalized,
		&infinitive,
		&partOfSpeech,
	)

	if err != nil {
		http.Error(
			w,
			"no previous lesson vocabulary found",
			http.StatusBadRequest,
		)
		return
	}

	// =====================================================
	// BUILD TRAINER QUESTION
	// =====================================================

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

	// =====================================================
	// RESPONSE
	// =====================================================

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(question)
}