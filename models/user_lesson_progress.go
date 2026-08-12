package models

import "time"

type UserLessonProgress struct {
	UserID               int        `json:"userId"`
	LessonID             int        `json:"lessonId"`
	TextCompleted        bool       `json:"textCompleted"`
	VocabularyCompleted  bool       `json:"vocabularyCompleted"`
	GrammarCompleted     bool       `json:"grammarCompleted"`
	ExaminatioCompleted  bool       `json:"examinatioCompleted"`
	Score                *int       `json:"score,omitempty"`
	StartedAt            *time.Time `json:"startedAt,omitempty"`
	CompletedAt          *time.Time `json:"completedAt,omitempty"`
}