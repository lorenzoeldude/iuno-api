package wordoftheday

import (
	"context"

	"iuno-api/db"
)

type WordOfTheDay struct {
	Lemma           string `json:"lemma"`
	LemmaNormalized string `json:"lemma_normalized"`
	Meaning         string `json:"meaning"`
}

func GetToday(ctx context.Context) (*WordOfTheDay, error) {

	var word WordOfTheDay

	err := db.Pool.QueryRow(
		ctx,
		`
		SELECT
			l.lemma,
			l.lemma_normalized,
			(
				SELECT d.definition
				FROM definitions d
				WHERE d.lemma_id = l.id
				ORDER BY d.sort_order ASC, d.id ASC
				LIMIT 1
			) AS meaning
		FROM word_of_the_day w
		JOIN lemmas l
			ON l.id = w.lemma_id
		WHERE w.date = CURRENT_DATE
		`,
	).Scan(
		&word.Lemma,
		&word.LemmaNormalized,
		&word.Meaning,
	)

	if err != nil {
		return nil, err
	}

	return &word, nil
}