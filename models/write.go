package models

type WriteRequest struct {
	Lemma Lemma				`json:"lemma"`
	Definitions []string	`json:"definitions"`
	Meanings []string		`json:"meanings"`
	Examples []string		`json:"examples"` 	
	Derivatives []string	`json:"derivatives"`
}
