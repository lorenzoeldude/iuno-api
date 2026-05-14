package models

type Word struct {
    ID      int    `json:"id"`
    Slug    string `json:"slug"`
    Lemma   string `json:"lemma"`
    Type    string `json:"type"`

    Meaning string `json:"meaning"`
    Definition string `json:"definition"`
}