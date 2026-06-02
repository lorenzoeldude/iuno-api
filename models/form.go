package models

type Form struct {
	ID                int `json:"id"`
	LemmaID           int `json:"lemma_id"`
	Form              string `json:"form"`
	FormNormalized    string `json:"form_normalized"`

	PartOfSpeech      string `json:"part_of_speech"`
	Number            string `json:"number"`

	// nouns/adjectives
	GrammaticalCase   *string `json:"grammatical_case"`
	Gender            *string `json:"gender"`

	// adjectives
	Degree			  *string `json:"degree"`

	// verbs
	Tense             *string `json:"tense"`
	Mood              *string `json:"mood"`
	Voice             *string `json:"voice"`
	Person            *int `json:"person"`
}