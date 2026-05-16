package handlers

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"iuno-api/db"
	"iuno-api/utils"
)

type TrainerQuestion struct {
	Lemma   string   `json:"lemma"`
	Correct string   `json:"correct"`
	Answers []string `json:"answers"`
}

func RandomTrainerHandler(w http.ResponseWriter, r *http.Request) {

	utils.EnableCORS(w)

	w.Header().Set("Content-Type", "application/json")

	rand.Seed(time.Now().UnixNano())

	var lemmaID int
	var lemma string
	var correct string

	// =====================================================
	// RANDOM LEMMA + CORRECT MEANING
	// =====================================================
	err := db.Pool.QueryRow(context.Background(), `
		SELECT
			l.id,
			l.lemma,
			m.meaning
		FROM lemmas l
		JOIN meanings m
			ON m.lemma_id = l.id
		ORDER BY RANDOM()
		LIMIT 1
	`).Scan(&lemmaID, &lemma, &correct)

	if err != nil {

		http.Error(
			w,
			"failed to fetch trainer question",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// RANDOM WRONG ANSWERS
	// exclude ALL meanings from the same lemma
	// =====================================================
	rows, err := db.Pool.Query(context.Background(), `
		SELECT meaning
		FROM meanings
		WHERE lemma_id != $1
		ORDER BY RANDOM()
		LIMIT 10
	`, lemmaID)

	if err != nil {

		http.Error(
			w,
			"failed to fetch wrong answers",
			http.StatusInternalServerError,
		)

		return
	}

	defer rows.Close()

	answers := []string{correct}

	used := map[string]bool{
		correct: true,
	}

	for rows.Next() {

		var wrong string

		err := rows.Scan(&wrong)
		if err != nil {
			continue
		}

		// avoid duplicates
		if used[wrong] {
			continue
		}

		used[wrong] = true
		answers = append(answers, wrong)

		if len(answers) == 3 {
			break
		}
	}

	// =====================================================
	// SAFETY CHECK
	// =====================================================
	if len(answers) < 3 {

		http.Error(
			w,
			"not enough unique meanings in database",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// SHUFFLE ANSWERS
	// =====================================================
	rand.Shuffle(len(answers), func(i, j int) {

		answers[i], answers[j] = answers[j], answers[i]

	})

	question := TrainerQuestion{
		Lemma:   lemma,
		Correct: correct,
		Answers: answers,
	}

	// =====================================================
	// RESPONSE
	// =====================================================
	json.NewEncoder(w).Encode(question)
}