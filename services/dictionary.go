package services

import (
	"context"
	"log"

	"iuno-api/db"
	"iuno-api/models"
	"iuno-api/services/morphology"
)

func GetWord(slug string) (models.DictionaryResponse, error) {

	var lemma models.Lemma

	log.Println("LOOKUP:", slug)

	err := db.Pool.QueryRow(context.Background(), `
		SELECT
			id,
			slug,
			lemma,
			type,
			definition,
			gender,
			declension,
			conjugation,
			perfect,
			supine,
			irregular,
			genitive
		FROM lemmas
		WHERE LOWER(slug) = LOWER($1)
	`, slug).Scan(
		&lemma.ID,
		&lemma.Slug,
		&lemma.Lemma,
		&lemma.Type,
		&lemma.Definition,
		&lemma.Gender,
		&lemma.Declension,
		&lemma.Conjugation,
		&lemma.Perfect,
		&lemma.Supine,
		&lemma.Irregular,
		&lemma.Genitive,
	)

	if err != nil {
		log.Println("DB ERROR:", err)
		return models.DictionaryResponse{}, err
	}

	log.Println("FOUND:", lemma.Lemma)

	// Generate Morphology
	forms := morphology.Generate(lemma)

	// Get Examples
	rows, err := db.Pool.Query(context.Background(), `
		SELECT id, latin
		FROM examples
		WHERE lemma_id = $1
		ORDER BY id ASC
	`, lemma.ID)

	if err != nil {
		log.Println("EXAMPLES ERROR:", err)
		return models.DictionaryResponse{}, err
	}

	defer rows.Close()

	var examples []models.Example

	for rows.Next() {

		var ex models.Example

		err := rows.Scan(
			&ex.ID,
			&ex.Latin,
		)

		if err != nil {
			log.Println("SCAN ERROR:", err)
			continue
		}

		examples = append(examples, ex)
	}


	// Get Meanings
	rows, err = db.Pool.Query(context.Background(), `
		SELECT id, meaning
		FROM meanings
		WHERE lemma_id = $1
		ORDER BY sort_order ASC
	`, lemma.ID)

	if err != nil {
		return models.DictionaryResponse{}, err
	}

	defer rows.Close()

	var meanings []models.Meaning

	for rows.Next() {

		var m models.Meaning

		err := rows.Scan(
			&m.ID,
			&m.English,
		)

		if err != nil {
			log.Println("SCAN ERROR:", err)
			continue
		}

		meanings = append(meanings, m)
	}

	// RESPONSE
	response := models.DictionaryResponse{
		Lemma:     lemma,
		Forms:    forms,
		Examples: examples,
		Meanings: meanings,
	}

	return response, nil
}