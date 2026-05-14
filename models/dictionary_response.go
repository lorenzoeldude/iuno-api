package models

type DictionaryResponse struct {
	Word  Word  `json:"word"`
	Forms []Form `json:"forms"`
}