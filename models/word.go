package models

type Word struct {

	ID int `json:"id"`

	Slug string `json:"slug"`

	Lemma string `json:"lemma"`

	Type string `json:"type"`

	Meaning string `json:"meaning"`

	Definition string `json:"definition"`

	// =====================================================
	// NOUN / ADJECTIVE
	// =====================================================

	Gender string `json:"gender"`

	Declension int `json:"declension"`

	// =====================================================
	// VERBS
	// =====================================================

	Conjugation int `json:"conjugation"`

	Stem string `json:"stem"`

	Perfect string `json:"perfect"`

	Supine string `json:"supine"`

	// =====================================================
	// ENGINE FLAGS
	// =====================================================

	Irregular bool `json:"irregular"`
}