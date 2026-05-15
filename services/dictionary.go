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

	// =====================================================
	// GENERATED MORPHOLOGY
	// =====================================================

	forms := morphology.Generate(word)

	// =====================================================
	// EXAMPLES
	// =====================================================

	rows, err := db.Pool.Query(context.Background(), `
		SELECT id, latin
		FROM examples
		WHERE lemma_id = $1
		ORDER BY id ASC
	`, word.ID)

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


		// =====================================================
	// MEANINGS
	// =====================================================

	rows, err = db.Pool.Query(context.Background(), `
		SELECT id, meaning
		FROM meanings
		WHERE lemma_id = $1
		ORDER BY sort_order ASC
	`, word.ID)

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



	// =====================================================
	// RESPONSE
	// =====================================================

	response := models.DictionaryResponse{
		Word:     word,
		Forms:    forms,
		Examples: examples,
		Meanings: meanings,
	}

	return response, nil
}