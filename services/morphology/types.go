package morphology

import "iuno-api/models"

//
// FEATURE BUNDLE
//
// Internal generation state.
//
type Features struct {
	Case   string
	Number string
	Gender string

	Tense string
	Mood  string
	Voice string

	Person int

	NonFinite string
}

//
// VERB STEM SYSTEM
//
// Latin verbs operate on multiple stems.
//
type VerbStems struct {
	Present string
	Perfect string
	Supine  string
}

//
// NOUN STEM SYSTEM
//
type NounStem struct {
	Base string
}

//
// ADJECTIVE STEM SYSTEM
//
type AdjectiveStem struct {
	Base string
}

//
// GENERATION RESULT
//
// Full morphology payload.
//
type GenerationResult struct {
	Word  Word   `json:"word"`
	Forms []models.Form `json:"forms"`
}