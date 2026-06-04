package services

import (
	"context"
	"log"

	"iuno-api/db"
	"iuno-api/models"
	"iuno-api/services/morphology"
)

func WriteWord(body models.WriteRequest) error {

	log.Println("lemma: ", body.Lemma.Lemma)
	log.Println("lemma: ", *body.Lemma.Supine)

	lemma := body.Lemma

	ctx := context.Background()

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

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
			neuter
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13
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
	).Scan(&lemma.ID)

	if err != nil {
		log.Println("inserting error", err)
		return err
	}

	log.Println("forms ->")

	forms := morphology.Generate(lemma)

	log.Println("generated forms:", len(forms))

	for _, form := range forms {

		// form.FormNormalized := morphology.NormalizeLatin(form.Form)

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
			form.PartOfSpeech,
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
			log.Println("error inserting forms", err)
			return err
		}
	}

	// DELETE OLD MEANINGS
	_, err = tx.Exec(ctx, `
		DELETE FROM meanings
			WHERE lemma_id = $1
		`, lemma.ID)

	if err != nil {
		log.Println("MEANINGS DELETE ERROR:", err)
		return err
	}

	// INSERT NEW MEANINGS
	for _, m := range body.Meanings {

		if m == "" {
			continue
		}

		_, err := tx.Exec(ctx, `
			INSERT INTO meanings (lemma_id, meaning)
			VALUES ($1, $2)
		`, lemma.ID, m)

		if err != nil {
			log.Println("MEANING INSERT ERROR:", err)
			return err
		}
	}

	// INSERT EXAMPLES
	for _, m := range body.Examples {

		if m == "" {
			continue
		}

		_, err := tx.Exec(ctx, `
			INSERT INTO examples (lemma_id, example)
			VALUES ($1, $2)
		`, lemma.ID, m)

		if err != nil {
			log.Println("EXAMPLE INSERT ERROR:", err)
			return err
		}
	}

	// INSERT DEFINITIONS
	for _, m := range body.Definitions {

		if m == "" {
			continue
		}

		_, err := tx.Exec(ctx, `
			INSERT INTO definitions (lemma_id, definition)
			VALUES ($1, $2)
		`, lemma.ID, m)

		if err != nil {
			log.Println("DEFINITIONS INSERT ERROR:", err)
			return err
		}
	}

	// INSERT DERIVATIVES
	for _, m := range body.Derivatives {

		if m == "" {
			continue
		}

		_, err := tx.Exec(ctx, `
			INSERT INTO derivatives (lemma_id, derivative)
			VALUES ($1, $2)
		`, lemma.ID, m)

		if err != nil {
			log.Println("DERIVATIVES INSERT ERROR:", err)
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