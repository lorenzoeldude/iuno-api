package morphology

import "iuno-api/models"

type Features struct {
	Case   string
	Number string
	Gender string

	Tense string
	Mood  string
	Voice string

	Person int

	VerbForm string
}

type GenerationResult struct {
	Lemma models.Lemma `json:"lemma"`
	Forms []models.Form `json:"forms"`
}