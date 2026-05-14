package models

type Example struct {
	ID    int    `json:"id"`
	Latin string `json:"latin"`
}

type DictionaryResponse struct {
	Word     Word      `json:"word"`
	Forms    []Form    `json:"forms"`
	Examples []Example `json:"examples"`
}