package handlers

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"iuno-api/db"
)

type TrainerQuestion struct {
	Lemma      string   `json:"lemma"`
	LemmaNormalized string   `json:"lemma_normalized"`
	Infinitive *string `json:"infinitive"`
	Correct    string   `json:"correct"`
	Answers    []string `json:"answers"`
	Definition string   `json:"definition"`
	Examples   []string `json:"examples"`
}

func RandomTrainerHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	rand.Seed(time.Now().UnixNano())

	var lemmaID int
	var lemma string
	var lemmaNormalized string
	var infinitive *string
	var partOfSpeech string
	var correct string

	// =====================================================
	// RANDOM LEMMA
	// =====================================================
	err := db.Pool.QueryRow(context.Background(), `
		SELECT
			id,
			lemma,
			lemma_normalized,
			infinitive,
			part_of_speech
		FROM lemmas
		ORDER BY RANDOM()
		LIMIT 1
	`).Scan(
		&lemmaID,
		&lemma,
		&lemmaNormalized,
		&infinitive,
		&partOfSpeech,
	)

	if err != nil {

		http.Error(
			w,
			"failed to fetch trainer lemma",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// FIRST MEANING OF LEMMA
	// =====================================================
	err = db.Pool.QueryRow(context.Background(), `
		SELECT meaning
		FROM meanings
		WHERE lemma_id = $1
		ORDER BY id
		LIMIT 1
	`, lemmaID).Scan(&correct)

	if err != nil {

		http.Error(
			w,
			"failed to fetch meaning",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// FIRST DEFINITION OF LEMMA
	// =====================================================
	var definition string

	err = db.Pool.QueryRow(context.Background(), `
		SELECT definition
		FROM definitions
		WHERE lemma_id = $1
		ORDER BY id
		LIMIT 1
	`, lemmaID).Scan(&definition)

	if err != nil {
		definition = ""
	}

	// =====================================================
	// FIRST 3 EXAMPLES
	// =====================================================
	examples := []string{}

	exampleRows, err := db.Pool.Query(context.Background(), `
		SELECT example
		FROM examples
		WHERE lemma_id = $1
		ORDER BY id
		LIMIT 3
	`, lemmaID)

	if err == nil {

		defer exampleRows.Close()

		for exampleRows.Next() {

			var example string

			if err := exampleRows.Scan(&example); err == nil {
				examples = append(examples, example)
			}
		}
	}

	// =====================================================
	// RANDOM WRONG ANSWERS
	// =====================================================
	rows, err := db.Pool.Query(context.Background(), `
		SELECT m.meaning
		FROM meanings m
		JOIN lemmas l ON l.id = m.lemma_id
		WHERE
			m.lemma_id != $1
			AND l.part_of_speech = $2
		ORDER BY RANDOM()
		LIMIT 20;
	`, lemmaID, partOfSpeech)

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

		if err := rows.Scan(&wrong); err != nil {
			continue
		}

		if used[wrong] {
			continue
		}

		used[wrong] = true
		answers = append(answers, wrong)

		if len(answers) == 4 {
			break
		}
	}

	// =====================================================
	// SAFETY CHECK
	// =====================================================
	if len(answers) < 4 {

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
		Lemma:      lemma,
		LemmaNormalized: lemmaNormalized,
		Infinitive: infinitive,
		Correct:    correct,
		Answers:    answers,
		Definition: definition,
		Examples:   examples,
	}

	// =====================================================
	// RESPONSE
	// =====================================================
	json.NewEncoder(w).Encode(question)
}
