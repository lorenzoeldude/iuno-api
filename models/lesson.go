package models

import (
	"time"
)

type Lesson struct {
	ID           int              `json:"id"`
	Title        string           `json:"title"`
	Image        string           `json:"image"`
	Introduction string           `json:"introduction"`
	Text         []TextPage       `json:"text"`
	Grammar      []GrammarSlide   `json:"grammar"`
	Exam         []ExamQuestion   `json:"exam"`
	IsPublished  bool             `json:"is_published"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

type TextPage struct {
	Text string `json:"text"`
}

type GrammarSlide struct {
	Type           string   `json:"type"`
	Title          string   `json:"title,omitempty"`
	Text           []string `json:"text,omitempty"`
	SentenceBefore string   `json:"sentenceBefore,omitempty"`
	Correct        string   `json:"correct,omitempty"`
	Options        []string `json:"options,omitempty"`
	Ending         string   `json:"ending,omitempty"`
}

type ExamQuestion struct {
	Type     string   `json:"type"`
	Question string   `json:"question,omitempty"`
	Before   string   `json:"before,omitempty"`
	After    string   `json:"after,omitempty"`
	Correct  string   `json:"correct"`
	Options  []string `json:"options"`
}