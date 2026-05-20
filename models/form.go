package models

type Form struct {
	Form string `json:"form"`
	Part string `json:"part"`

	// shared
	Number string `json:"number"`

	// nouns/adjectives
	Case string `json:"case"`
	Gender string `json:"gender"`

	// verbs
	Person int `json:"person"`
	Tense string `json:"tense"`
	Mood string `json:"mood"`
	Voice string `json:"voice"`
	NonFinite string `json:"non_finite"`
}