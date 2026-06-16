package models

type Example struct {
	ID    int    `json:"id"`
	Latin string `json:"latin"`
}

type Meaning struct {
	ID    int    `json:"id"`
	Meaning string `json:"meaning"`
	GovernsCase *string `json:"governs_case"`
}

type Definition struct {
	ID    int    `json:"id"`
	Definition string `json:"definition"`
}

type Derivative struct {
	ID    int    `json:"id"`
	Derivative string `json:"derivative"`
}

type DictionaryResponse struct {
	Lemma     Lemma      `json:"lemma"`
	Forms    []Form    `json:"forms"`
	Examples []Example `json:"examples"`
	Meanings []Meaning `json:"meanings"`
	Definitions []Definition `json:"definitions"`
	Derivatives []Derivative `json:"derivatives"`
}