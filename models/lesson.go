package models

import "time"

type Lesson struct {
	ID           int            `json:"id"`
	Title        string         `json:"title"`
	Image        string         `json:"image"`
	Introduction string         `json:"introduction"`
	Text         []TextPage     `json:"text"`
	Grammar      []GrammarPage  `json:"grammar"`
	Exam         []ExamQuestion `json:"exam"`
	IsPublished  bool           `json:"is_published"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type TextPage struct {
	Text string `json:"text"`
}

// =====================================================
// GRAMMAR
// =====================================================

type GrammarPage struct {
	Title  string         `json:"title"`
	Blocks []GrammarBlock `json:"blocks"`
}

type GrammarBlock struct {
	Type string `json:"type"`

	// Paragraph
	Text string `json:"text,omitempty"`

	// Emphasis
	// Uses Text

	// Grammar Diagram
	Words        []GrammarDiagramWord `json:"words,omitempty"`
	Explanations []string             `json:"explanations,omitempty"`

	// Sentence
	// Uses Text

	// Question
	Question    string   `json:"question,omitempty"`
	Correct     string   `json:"correct,omitempty"`
	Options     []string `json:"options,omitempty"`
	Explanation string   `json:"explanation,omitempty"`

	// Sentence Question
	Sentence string `json:"sentence,omitempty"`

	// Quiz blocks
	SentenceBefore string `json:"sentenceBefore,omitempty"`
	Ending          string `json:"ending,omitempty"`
}

type GrammarDiagramWord struct {
	Word  string `json:"word"`
	Case  string `json:"case"`
	Color string `json:"color"`
}

// =====================================================
// EXAM
// =====================================================

type ExamQuestion struct {
	Type     string   `json:"type"`
	Sentence string   `json:"sentence,omitempty"`
	Question string   `json:"question,omitempty"`
	Before   string   `json:"before,omitempty"`
	After    string   `json:"after,omitempty"`
	Correct  string   `json:"correct"`
	Options  []string `json:"options"`
}