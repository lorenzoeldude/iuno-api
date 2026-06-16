package models

type Lemma struct {
	ID int `json:"id"`
	Lemma string `json:"lemma"`
	LemmaNormalized string `json:"lemma_normalized"`
	PartOfSpeech string `json:"part_of_speech"`

	// NOUN / ADJECTIVE
	Gender *string `json:"gender"`
	Declension *int `json:"declension"`
	Genitive *string `json:"genitive"`

	// ADJECTIVES
	Feminine *string `json:"feminine"`
	Neuter *string `json:"neuter"`

	// VERBS
	Conjugation *int `json:"conjugation"`
	Perfect *string `json:"perfect"`
	Supine *string `json:"supine"`
	Infinitive *string `json:"infinitive"`
	
	// PRONOUN
	PronounType *string `json:"pronoun_type"`

	// PREPOSITION
	GovernsCase *string `json:"governs_case"`

	// ENGINE FLAGS
	Irregular bool `json:"irregular"`
}