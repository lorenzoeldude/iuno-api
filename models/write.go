package models

type WriteRequest struct {
	Lemma Lemma				`json:"lemma"`
	Definitions []string	`json:"definitions"`
	Meanings []MeaningInput	`json:"meanings"`
	Examples []string		`json:"examples"` 	
	Derivatives []string	`json:"derivatives"`
	ManualForms []Form 		`json:"manual_forms"`
}

type MeaningInput struct {
    Meaning string `json:"meaning"`
    GovernsCase *string `json:"governs_case"`
}