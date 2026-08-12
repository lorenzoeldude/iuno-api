package models

type UserTrainingStats struct {
	UserID            int     `json:"userId"`
	QuestionsAnswered int     `json:"questionsAnswered"`
	CorrectAnswers    int     `json:"correctAnswers"`
	IncorrectAnswers  int     `json:"incorrectAnswers"`
	TrainingSessions  int     `json:"trainingSessions"`
	Sestertii         int     `json:"sestertii"`
	LastTrainedAt     *string `json:"lastTrainedAt,omitempty"`
}