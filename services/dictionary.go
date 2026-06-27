package services

import (
	"context"
	"log"

	"iuno-api/db"
	"iuno-api/models"
)


// =====================================================
// GET WORD BY NORMALIZED LEMMA
// =====================================================

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
			feminine,
			neuter,
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
		&lemma.Feminine,
		&lemma.Neuter,
		&lemma.Irregular,
	)

	if err != nil {
		return models.DictionaryResponse{}, err
	}


	return getWordData(lemma)
}



// =====================================================
// GET WORD BY ID (ADMIN EDITOR)
// =====================================================

func GetWordByID(id int) (models.DictionaryResponse, error) {

	var lemma models.Lemma

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
			feminine,
			neuter,
			irregular
		FROM lemmas
		WHERE id = $1
	`, id).Scan(
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
		&lemma.Feminine,
		&lemma.Neuter,
		&lemma.Irregular,
	)

	if err != nil {
		return models.DictionaryResponse{}, err
	}


	return getWordData(lemma)
}



// =====================================================
// COMMON DATA LOADING
// =====================================================

func getWordData(lemma models.Lemma) (models.DictionaryResponse, error) {


	// ================= FORMS =================

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

		forms = append(forms, form)
	}



	// ================= EXAMPLES =================

	exampleRows, err := db.Pool.Query(context.Background(), `
		SELECT id, example
		FROM examples
		WHERE lemma_id = $1
		ORDER BY id ASC
	`, lemma.ID)

	if err != nil {
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
			continue
		}

		examples = append(examples, ex)
	}



	// ================= MEANINGS =================

	meaningRows, err := db.Pool.Query(context.Background(), `
		SELECT id, meaning, governs_case
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
			&m.GovernsCase,
		)

		if err != nil {
			continue
		}

		meanings = append(meanings, m)
	}



	// ================= DEFINITIONS =================

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

		var d models.Definition

		err := definitionRows.Scan(
			&d.ID,
			&d.Definition,
		)

		if err != nil {
			continue
		}

		definitions = append(definitions, d)
	}



	// ================= DERIVATIVES =================

	derivativeRows, err := db.Pool.Query(context.Background(), `
		SELECT id, derivative
		FROM derivatives
		WHERE lemma_id = $1
		ORDER BY sort_order ASC
	`, lemma.ID)

	if err != nil {
		return models.DictionaryResponse{}, err
	}

	defer derivativeRows.Close()


	var derivatives []models.Derivative

	for derivativeRows.Next() {

		var d models.Derivative

		err := derivativeRows.Scan(
			&d.ID,
			&d.Derivative,
		)

		if err != nil {
			continue
		}

		derivatives = append(derivatives, d)
	}



	return models.DictionaryResponse{
		Lemma:       lemma,
		Forms:       forms,
		Examples:    examples,
		Meanings:    meanings,
		Definitions: definitions,
		Derivatives: derivatives,
	}, nil
}