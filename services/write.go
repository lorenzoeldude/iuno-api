package services

import (
	"context"
	"log"

	"iuno-api/db"
	"iuno-api/models"
	"iuno-api/services/morphology"
)

func WriteWord(body models.WriteRequest) error {

	log.Println("lemma:", body.Lemma.Lemma)

	lemma := body.Lemma
	ctx := context.Background()

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// =====================================================
	// CREATE OR UPDATE LEMMA
	// =====================================================

	var existingID int

	if lemma.PartOfSpeech == "verb" {

		err = tx.QueryRow(ctx, `
			SELECT id
			FROM lemmas
			WHERE lemma_normalized = $1
			AND part_of_speech = 'verb'
			AND infinitive = $2
		`,
			lemma.LemmaNormalized,
			nullString(lemma.Infinitive),
		).Scan(&existingID)

	} else {

		err = tx.QueryRow(ctx, `
			SELECT id
			FROM lemmas
			WHERE lemma_normalized = $1
			AND part_of_speech = $2
		`,
			lemma.LemmaNormalized,
			lemma.PartOfSpeech,
		).Scan(&existingID)

	}

	if err == nil {

		// UPDATE EXISTING
		lemma.ID = existingID

		_, err = tx.Exec(ctx, `
			UPDATE lemmas
			SET
				lemma = $2,
				part_of_speech = $3,
				gender = $4,
				declension = $5,
				conjugation = $6,
				perfect = $7,
				supine = $8,
				genitive = $9,
				infinitive = $10,
				irregular = $11,
				feminine = $12,
				neuter = $13,
				pronoun_type = $14
			WHERE id = $1
		`,
			lemma.ID,
			lemma.Lemma,
			lemma.PartOfSpeech,
			lemma.Gender,
			nullInt(lemma.Declension),
			nullInt(lemma.Conjugation),
			nullString(lemma.Perfect),
			nullString(lemma.Supine),
			nullString(lemma.Genitive),
			nullString(lemma.Infinitive),
			lemma.Irregular,
			nullString(lemma.Feminine),
			nullString(lemma.Neuter),
			nullString(lemma.PronounType),
		)

		if err != nil {
			return err
		}

	} else {

		// INSERT NEW
		err = tx.QueryRow(ctx, `
			INSERT INTO lemmas (
				lemma,
				lemma_normalized,
				part_of_speech,
				gender,
				declension,
				conjugation,
				perfect,
				supine,
				genitive,
				infinitive,
				irregular,
				feminine,
				neuter,
				pronoun_type
			)
			VALUES (
				$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14
			)
			RETURNING id
		`,
			lemma.Lemma,
			lemma.LemmaNormalized,
			lemma.PartOfSpeech,
			lemma.Gender,
			nullInt(lemma.Declension),
			nullInt(lemma.Conjugation),
			nullString(lemma.Perfect),
			nullString(lemma.Supine),
			nullString(lemma.Genitive),
			nullString(lemma.Infinitive),
			lemma.Irregular,
			nullString(lemma.Feminine),
			nullString(lemma.Neuter),
			nullString(lemma.PronounType),
		).Scan(&lemma.ID)

		if err != nil {
			return err
		}
	}

	// =====================================================
	// CLEAR CHILD TABLES
	// =====================================================

	_, err = tx.Exec(ctx, `DELETE FROM forms WHERE lemma_id = $1`, lemma.ID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM meanings WHERE lemma_id = $1`, lemma.ID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM examples WHERE lemma_id = $1`, lemma.ID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM definitions WHERE lemma_id = $1`, lemma.ID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM derivatives WHERE lemma_id = $1`, lemma.ID)
	if err != nil {
		return err
	}

	// =====================================================
	// GENERATE FORMS
	// =====================================================

	var forms []models.Form

	if lemma.PartOfSpeech == "preposition" ||
		lemma.PartOfSpeech == "conjunction" ||
		lemma.PartOfSpeech == "interjection" ||
		lemma.PartOfSpeech == "adverb" {

		forms = []models.Form{
			{
				LemmaID:       lemma.ID,
				Form:          lemma.Lemma,
				FormNormalized: morphology.NormalizeLatin(lemma.Lemma),
				PartOfSpeech:  lemma.PartOfSpeech,
			},
		}

	} else {

		if lemma.Irregular && len(body.ManualForms) > 0 {
			forms = body.ManualForms
		} else {
			forms = morphology.Generate(lemma)
		}
	}

	// INSERT FORMS
	for _, form := range forms {

		_, err := tx.Exec(ctx, `
			INSERT INTO forms (
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
			)
			VALUES (
				$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13
			)
		`,
			lemma.ID,
			form.Form,
			morphology.NormalizeLatin(form.Form),
			lemma.PartOfSpeech,
			form.GrammaticalCase,
			form.Number,
			form.Gender,
			form.Tense,
			form.Mood,
			form.Voice,
			form.Person,
			form.Degree,
			form.FormType,
		)

		if err != nil {
			return err
		}
	}

	// INSERT MEANINGS
	for _, m := range body.Meanings {

		if m.Meaning == "" {
			continue
		}

		_, err := tx.Exec(ctx, `
			INSERT INTO meanings (
				lemma_id,
				meaning,
				governs_case
			)
			VALUES ($1, $2, $3)
		`,
			lemma.ID,
			m.Meaning,
			m.GovernsCase,
		)

		if err != nil {
			return err
		}
	}

	// INSERT EXAMPLES
	for _, ex := range body.Examples {

		if ex == "" {
			continue
		}

		_, err := tx.Exec(ctx, `
			INSERT INTO examples (lemma_id, example)
			VALUES ($1, $2)
		`, lemma.ID, ex)

		if err != nil {
			return err
		}
	}

	// INSERT DEFINITIONS
	for _, def := range body.Definitions {

		if def == "" {
			continue
		}

		_, err := tx.Exec(ctx, `
			INSERT INTO definitions (lemma_id, definition)
			VALUES ($1, $2)
		`, lemma.ID, def)

		if err != nil {
			return err
		}
	}

	// INSERT DERIVATIVES
	for _, d := range body.Derivatives {

		if d == "" {
			continue
		}

		_, err := tx.Exec(ctx, `
			INSERT INTO derivatives (lemma_id, derivative)
			VALUES ($1, $2)
		`, lemma.ID, d)

		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func nullString(s *string) any {
	if s == nil {
		return nil
	}
	return *s
}

func nullInt(i *int) any {
	if i == nil {
		return nil
	}
	return *i
}