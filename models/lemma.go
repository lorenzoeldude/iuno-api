package models

type Lemma struct {
	ID int `json:"id"`
	Slug string `json:"slug"`
	Lemma string `json:"lemma"`
	Type string `json:"type"`
	Definition string `json:"definition"`

	// NOUN / ADJECTIVE
	Gender *string `json:"gender"`
	Declension *int `json:"declension"`
	Genitive *string `json:"genitive"`

	// VERBS
	Conjugation *int `json:"conjugation"`
	Perfect *string `json:"perfect"`
	Supine *string `json:"supine"`
	Infinitive *string `json:"infinitive"`

	// ENGINE FLAGS
	Irregular bool `json:"irregular"`
}