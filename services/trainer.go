package services

import (
	"context"
	"math/rand"

	"iuno-api/db"
)

type TrainerWord struct {
	Latin   string   `json:"latin"`
	Correct string   `json:"correct"`
	Answers []string `json:"answers"`
}

func GetTrainerWords() ([]TrainerWord, error) {

	rows, err := db.Pool.Query(context.Background(), `
		SELECT latin, translation_1
		FROM words
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type W struct {
		Latin string
		Trans string
	}

	var words []W

	for rows.Next() {
		var w W
		rows.Scan(&w.Latin, &w.Trans)
		words = append(words, w)
	}

	if len(words) < 3 {
		return nil, err
	}

	rand.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
	})

	var result []TrainerWord

	for i := 0; i < len(words) && i < 10; i++ {

		cur := words[i]

		answers := []string{cur.Trans}

		for len(answers) < 3 {
			r := words[rand.Intn(len(words))].Trans

			exists := false
			for _, a := range answers {
				if a == r {
					exists = true
				}
			}

			if !exists {
				answers = append(answers, r)
			}
		}

		rand.Shuffle(len(answers), func(i, j int) {
			answers[i], answers[j] = answers[j], answers[i]
		})

		result = append(result, TrainerWord{
			Latin:   cur.Latin,
			Correct: cur.Trans,
			Answers: answers,
		})
	}

	return result, nil
}