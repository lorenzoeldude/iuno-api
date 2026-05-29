package services

import (
	"context"
	"log"

	"iuno-api/db"
	"iuno-api/models"
	"iuno-api/services/morphology"
)

func WriteWord(body models.WriteRequest) error {

	// log.Println("lemma: ", body.Lemma)

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
			irregular
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11
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
	).Scan(&lemma.ID)

	if err != nil {
		log.Println("inserting error", err)
		return err
	}

	forms := morphology.Generate(lemma)

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
				person
			)
			VALUES (
				$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11
			)
		`,
			lemma.ID,
			form.Form,
			form.FormNormalized,
			form.PartOfSpeech,
			form.GrammaticalCase,
			form.Number,
			form.Gender,
			form.Tense,
			form.Mood,
			form.Voice,
			form.Person,
		)

		if err != nil {
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
			// http.Error(w, "failed to insert meanings", http.StatusInternalServerError)
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