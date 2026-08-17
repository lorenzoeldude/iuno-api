package wordoftheday

import (
	"context"

	"iuno-api/db"
)

type WordOfTheDay struct {
	Lemma           string   `json:"lemma"`
	LemmaNormalized string   `json:"lemma_normalized"`
	Meanings        []string `json:"meanings"`
}

func GetToday(ctx context.Context) (*WordOfTheDay, error) {

	var word WordOfTheDay

	err := db.Pool.QueryRow(
		ctx,
		`
		SELECT
			l.lemma,
			l.lemma_normalized,
			COALESCE(
				ARRAY(
					SELECT m.meaning
					FROM meanings m
					WHERE m.lemma_id = l.id
					ORDER BY m.sort_order ASC, m.id ASC
				),
				ARRAY[]::text[]
			) AS meanings
		FROM word_of_the_day w
		JOIN lemmas l
			ON l.id = w.lemma_id
		WHERE w.date = CURRENT_DATE
		`,
	).Scan(
		&word.Lemma,
		&word.LemmaNormalized,
		&word.Meanings,
	)

	if err != nil {
		return nil, err
	}

	return &word, nil
}