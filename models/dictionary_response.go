package models

type Example struct {
	ID    int    `json:"id"`
	Latin string `json:"latin"`
}

type Meaning struct {
	ID    int    `json:"id"`
	English string `json:"english"`
}

type DictionaryResponse struct {
	Word     Word      `json:"word"`
	Forms    []Form    `json:"forms"`
	Examples []Example `json:"examples"`
	Meanings []Meaning `json:"meanings"`
}