package models

import "time"

type WordList struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type WordListLemma struct {
	ListID    int       `json:"list_id"`
	LemmaID   int       `json:"lemma_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateWordListRequest struct {
	Name string `json:"name"`
}

type AddLemmaToListRequest struct {
	ListID  int `json:"list_id"`
	LemmaID int `json:"lemma_id"`
}