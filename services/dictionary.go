package services

import (
	"context"
	// "log"

	"iuno-api/db"
	"iuno-api/models"
)

func GetWord(slug string) (models.Word, error) {

    var word models.Word

    err := db.Pool.QueryRow(context.Background(), `
        SELECT id, slug, lemma, type, meaning, definition
        FROM lemmas
        WHERE slug = $1
    `, slug).Scan(
        &word.ID,
        &word.Slug,
        &word.Lemma,
        &word.Type,
        &word.Meaning,
        &word.Definition,
    )

    return word, err
}

func GetForms(lemmaID int) ([]models.Form, error) {

    rows, err := db.Pool.Query(context.Background(), `
        SELECT form, part, grammatical_case, number, gender, tense, mood, voice, person
        FROM forms
        WHERE lemma_id = $1
    `, lemmaID)

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var result []models.Form

    for rows.Next() {
        var f models.Form

        rows.Scan(
            &f.Form,
            &f.Part,
            &f.Case,
            &f.Number,
            &f.Gender,
            &f.Tense,
            &f.Mood,
            &f.Voice,
            &f.Person,
        )

        result = append(result, f)
    }

    return result, nil
}