package handlers

import (
	"context"
	"math/rand"

	"iuno-api/db"
)

func BuildTrainerQuestion(
	ctx context.Context,
	lemmaID int,
	lemma string,
	lemmaNormalized string,
	infinitive *string,
	partOfSpeech string,
) (*TrainerQuestion, error) {

	// =====================================================
	// GET CORRECT MEANING
	// =====================================================

	var correct string

	err := db.Pool.QueryRow(ctx, `
		SELECT meaning
		FROM meanings
		WHERE lemma_id = $1
		ORDER BY id
		LIMIT 1
	`,
		lemmaID,
	).Scan(&correct)

	if err != nil {
		return nil, err
	}

	// =====================================================
	// GET DEFINITION
	// =====================================================

	definition := ""

	_ = db.Pool.QueryRow(ctx, `
		SELECT definition
		FROM definitions
		WHERE lemma_id = $1
		ORDER BY id
		LIMIT 1
	`,
		lemmaID,
	).Scan(&definition)

	// =====================================================
	// GET EXAMPLES
	// =====================================================

	examples := []string{}

	rows, err := db.Pool.Query(ctx, `
		SELECT example
		FROM examples
		WHERE lemma_id = $1
		ORDER BY id
		LIMIT 3
	`,
		lemmaID,
	)

	if err == nil {

		defer rows.Close()

		for rows.Next() {

			var example string

			if err := rows.Scan(&example); err == nil {
				examples = append(examples, example)
			}
		}
	}

	// =====================================================
	// GET WRONG ANSWERS
	// =====================================================

	rows, err = db.Pool.Query(ctx, `
		SELECT m.meaning
		FROM meanings m
		JOIN lemmas l
			ON l.id = m.lemma_id
		WHERE
			m.lemma_id != $1
			AND l.part_of_speech = $2
		ORDER BY RANDOM()
		LIMIT 20
	`,
		lemmaID,
		partOfSpeech,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	answers := []string{correct}

	used := map[string]bool{
		correct: true,
	}

	for rows.Next() {

		var meaning string

		if err := rows.Scan(&meaning); err != nil {
			continue
		}

		if used[meaning] {
			continue
		}

		used[meaning] = true
		answers = append(answers, meaning)

		if len(answers) == 4 {
			break
		}
	}

	for len(answers) < 4 {
		answers = append(answers, "—")
	}

	rand.Shuffle(len(answers), func(i, j int) {
		answers[i], answers[j] = answers[j], answers[i]
	})

	return &TrainerQuestion{
		Lemma:            lemma,
		LemmaNormalized:  lemmaNormalized,
		Infinitive:       infinitive,
		Correct:          correct,
		Answers:          answers,
		Definition:       definition,
		Examples:         examples,
	}, nil
}