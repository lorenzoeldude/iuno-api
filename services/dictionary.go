package services

import (
	"context"
	"log"

	"iuno-api/db"
	"iuno-api/models"
	// "iuno-api/services/morphology"
)

func GetWord(lemma_normalized string) (models.DictionaryResponse, error) {

	var lemma models.Lemma

	log.Println("LOOKUP:", lemma_normalized)

	err := db.Pool.QueryRow(context.Background(), `
		SELECT
			id,
			lemma,
			lemma_normalized,
			part_of_speech,
			gender,
			declension,
			genitive,
			conjugation,
			perfect,
			supine,
			infinitive,
			irregular
		FROM lemmas
		WHERE LOWER(lemma_normalized) = LOWER($1)
	`, lemma_normalized).Scan(
		&lemma.ID,
		&lemma.Lemma,
		&lemma.LemmaNormalized,
		&lemma.PartOfSpeech,
		&lemma.Gender,
		&lemma.Declension,
		&lemma.Genitive,
		&lemma.Conjugation,
		&lemma.Perfect,
		&lemma.Supine,
		&lemma.Infinitive,
		&lemma.Irregular,
	)

	if err != nil {
		log.Println("DB ERROR:", err)
		return models.DictionaryResponse{}, err
	}

	// log.Println("FOUND:", lemma.Lemma)

	// Get Forms
	formRows, err := db.Pool.Query(context.Background(), `
		SELECT
			id,
			lemma_id,
			form,
			form_normalized,
			part_of_speech,
			grammatical_case,
			number,
			gender,
			tense,
			mood,
			voice,
			person,
			degree,
			form_type
		FROM forms
		WHERE lemma_id = $1
	`, lemma.ID)

	if err != nil {
		log.Println("FORMS ERROR:", err)
		return models.DictionaryResponse{}, err
	}

	defer formRows.Close()

	var forms []models.Form

	for formRows.Next() {

		var form models.Form

		err := formRows.Scan(
			&form.ID,
			&form.LemmaID,
			&form.Form,
			&form.FormNormalized,
			&form.PartOfSpeech,
			&form.GrammaticalCase,
			&form.Number,
			&form.Gender,
			&form.Tense,
			&form.Mood,
			&form.Voice,
			&form.Person,
			&form.Degree,
			&form.FormType,
		)

		if err != nil {
			log.Println("FORM SCAN ERROR:", err)
			continue
		}

		// log.Println("Gender: ", *form.Gender)

		forms = append(forms, form)
	}

	// Get Examples
	exampleRows, err := db.Pool.Query(context.Background(), `
		SELECT id, example
		FROM examples
		WHERE lemma_id = $1
		ORDER BY id ASC
	`, lemma.ID)

	if err != nil {
		log.Println("EXAMPLES ERROR:", err)
		return models.DictionaryResponse{}, err
	}

	defer exampleRows.Close()

	var examples []models.Example

	for exampleRows.Next() {

		var ex models.Example

		err := exampleRows.Scan(
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
	meaningRows, err := db.Pool.Query(context.Background(), `
		SELECT id, meaning
		FROM meanings
		WHERE lemma_id = $1
		ORDER BY sort_order ASC
	`, lemma.ID)

	if err != nil {
		return models.DictionaryResponse{}, err
	}

	defer meaningRows.Close()

	var meanings []models.Meaning

	for meaningRows.Next() {

		var m models.Meaning

		err := meaningRows.Scan(
			&m.ID,
			&m.Meaning,
		)

		if err != nil {
			log.Println("SCAN ERROR:", err)
			continue
		}

		meanings = append(meanings, m)
	}

	// Get Definitions
	definitionRows, err := db.Pool.Query(context.Background(), `
		SELECT id, definition
		FROM definitions
		WHERE lemma_id = $1
		ORDER BY sort_order ASC
	`, lemma.ID)

	if err != nil {
		return models.DictionaryResponse{}, err
	}

	defer definitionRows.Close()

	var definitions []models.Definition

	for definitionRows.Next() {

		var m models.Definition

		err := definitionRows.Scan(
			&m.ID,
			&m.Definition,
		)

		if err != nil {
			log.Println("SCAN ERROR:", err)
			continue
		}

		definitions = append(definitions, m)
	}

	// Get Derivatives
	derivativeRows, err := db.Pool.Query(context.Background(), `
		SELECT id, derivative
		FROM derivatives
		WHERE lemma_id = $1
		ORDER BY sort_order ASC
	`, lemma.ID)

	if err != nil {
		log.Println("derivative error: ", err)
		return models.DictionaryResponse{}, err
	}

	defer derivativeRows.Close()

	var derivatives []models.Derivative

	for derivativeRows.Next() {

		var m models.Derivative

		err := derivativeRows.Scan(
			&m.ID,
			&m.Derivative,
		)

		if err != nil {
			log.Println("SCAN ERROR:", err)
			continue
		}

		derivatives = append(derivatives, m)
	}
	
// 	if err := derivativeRows.Err(); err != nil {
//     log.Println("ROWS ERROR:", err)
// }
	log.Println("derivatives; ", derivatives)

	// RESPONSE
	response := models.DictionaryResponse{
		Lemma:     lemma,
		Forms:    forms,
		Examples: examples,
		Meanings: meanings,
		Definitions: definitions,
		Derivatives: derivatives,
	}

	return response, nil
}