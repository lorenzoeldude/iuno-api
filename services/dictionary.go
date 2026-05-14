package services

import (
	"context"
	"log"

	"iuno-api/db"
	"iuno-api/models"
	"iuno-api/services/morphology"
)

func GetWord(slug string) (models.DictionaryResponse, error) {

	var word models.Word

	log.Println("LOOKUP:", slug)

	err := db.Pool.QueryRow(context.Background(), `
		SELECT
			id,
			slug,
			lemma,
			type,
			meaning,
			definition,
			gender,
			declension,
			conjugation,
			stem,
			perfect,
			supine,
			is_irregular
		FROM lemmas
		WHERE LOWER(slug) = LOWER($1)
	`, slug).Scan(
		&word.ID,
		&word.Slug,
		&word.Lemma,
		&word.Type,
		&word.Meaning,
		&word.Definition,
		&word.Gender,
		&word.Declension,
		&word.Conjugation,
		&word.Stem,
		&word.Perfect,
		&word.Supine,
		&word.Irregular,
	)

	if err != nil {
		log.Println("DB ERROR:", err)
		return models.DictionaryResponse{}, err
	}

	log.Println("FOUND:", word.Lemma)

	forms := morphology.Generate(word)

	response := models.DictionaryResponse{
		Word:  word,
		Forms: forms,
	}

	return response, nil
}