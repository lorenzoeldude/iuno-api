package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"iuno-api/db"
	"iuno-api/middleware"
	"iuno-api/utils"

	"github.com/jackc/pgx/v5"
)

// =====================================================
// GET TRAINING STATS
// GET /api/training/stats
// =====================================================

func GetTrainingStatsHandler(w http.ResponseWriter, r *http.Request) {

	// =====================================================
	// METHOD CHECK
	// =====================================================

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// =====================================================
	// AUTH
	// =====================================================

	claimsRaw := r.Context().Value(middleware.UserContextKey)

	claims, ok := claimsRaw.(*utils.Claims)

	if !ok || claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID

	// =====================================================
	// GET STATS
	// =====================================================

	var stats struct {
		UserID            int     `json:"userId"`
		QuestionsAnswered int     `json:"questionsAnswered"`
		CorrectAnswers    int     `json:"correctAnswers"`
		IncorrectAnswers  int     `json:"incorrectAnswers"`
		TrainingSessions  int     `json:"trainingSessions"`
		Sestertii         int     `json:"sestertii"`
		LastTrainedAt     *string `json:"lastTrainedAt"`
	}

	err := db.Pool.QueryRow(
		r.Context(),
		`
		SELECT
			uts.user_id,
			uts.questions_answered,
			uts.correct_answers,
			uts.incorrect_answers,
			uts.training_sessions,
			u.sestertii,
			TO_CHAR(
				uts.last_trained_at,
				'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"'
			)
		FROM user_training_stats uts
		JOIN users u ON u.id = uts.user_id
		WHERE uts.user_id = $1
		`,
		userID,
	).Scan(
		&stats.UserID,
		&stats.QuestionsAnswered,
		&stats.CorrectAnswers,
		&stats.IncorrectAnswers,
		&stats.TrainingSessions,
		&stats.Sestertii,
		&stats.LastTrainedAt,
	)

	// =====================================================
	// NO STATS YET
	// =====================================================

	if err != nil {

		if err == pgx.ErrNoRows {

			stats.UserID = userID

			// Get sestertii directly from users.
			err = db.Pool.QueryRow(
				r.Context(),
				`
				SELECT sestertii
				FROM users
				WHERE id = $1
				`,
				userID,
			).Scan(
				&stats.Sestertii,
			)

			if err != nil && err != pgx.ErrNoRows {

				log.Println("GET USER SESTERTII ERROR:", err)

				http.Error(
					w,
					"failed to get user sestertii",
					http.StatusInternalServerError,
				)

				return
			}

			w.Header().Set("Content-Type", "application/json")

			json.NewEncoder(w).Encode(stats)

			return
		}

		log.Println("GET TRAINING STATS ERROR:", err)

		http.Error(
			w,
			"failed to get training stats",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// RESPONSE
	// =====================================================

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(stats)
}

// =====================================================
// RECORD TRAINING ATTEMPT
// POST /api/training/attempt
// =====================================================

type RecordTrainingAttemptRequest struct {
	LemmaID *int `json:"lemmaId"`
	Correct bool `json:"correct"`
}

func RecordTrainingAttemptHandler(w http.ResponseWriter, r *http.Request) {

	// =====================================================
	// METHOD CHECK
	// =====================================================

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// =====================================================
	// AUTH
	// =====================================================

	claimsRaw := r.Context().Value(middleware.UserContextKey)

	claims, ok := claimsRaw.(*utils.Claims)

	if !ok || claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID

	// =====================================================
	// PARSE BODY
	// =====================================================

	var req RecordTrainingAttemptRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// =====================================================
	// DATABASE TRANSACTION
	// =====================================================

	tx, err := db.Pool.Begin(r.Context())

	if err != nil {

		log.Println("TRAINING TRANSACTION ERROR:", err)

		http.Error(
			w,
			"failed to record training attempt",
			http.StatusInternalServerError,
		)

		return
	}

	defer tx.Rollback(r.Context())

	// =====================================================
	// RECORD ATTEMPT
	// =====================================================

	_, err = tx.Exec(
		r.Context(),
		`
		INSERT INTO training_attempts (
			user_id,
			lemma_id,
			correct
		)
		VALUES ($1, $2, $3)
		`,
		userID,
		req.LemmaID,
		req.Correct,
	)

	if err != nil {

		log.Println("TRAINING ATTEMPT INSERT ERROR:", err)

		http.Error(
			w,
			"failed to record training attempt",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// CREATE STATS ROW IF NEEDED
	// =====================================================

	_, err = tx.Exec(
		r.Context(),
		`
		INSERT INTO user_training_stats (
			user_id
		)
		VALUES ($1)
		ON CONFLICT (user_id) DO NOTHING
		`,
		userID,
	)

	if err != nil {

		log.Println("TRAINING STATS CREATE ERROR:", err)

		http.Error(
			w,
			"failed to update training stats",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// UPDATE AGGREGATE STATS
	// =====================================================

	if req.Correct {

		_, err = tx.Exec(
			r.Context(),
			`
			UPDATE user_training_stats
			SET
				questions_answered = questions_answered + 1,
				correct_answers = correct_answers + 1,
				last_trained_at = NOW(),
				updated_at = NOW()
			WHERE user_id = $1
			`,
			userID,
		)

	} else {

		_, err = tx.Exec(
			r.Context(),
			`
			UPDATE user_training_stats
			SET
				questions_answered = questions_answered + 1,
				incorrect_answers = incorrect_answers + 1,
				last_trained_at = NOW(),
				updated_at = NOW()
			WHERE user_id = $1
			`,
			userID,
		)
	}

	if err != nil {

		log.Println("TRAINING STATS UPDATE ERROR:", err)

		http.Error(
			w,
			"failed to update training stats",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// AWARD SESTERTIUS
	// =====================================================

	_, err = tx.Exec(
		r.Context(),
		`
		UPDATE users
		SET
			sestertii = sestertii + 1
		WHERE id = $1
		`,
		userID,
	)

	if err != nil {

		log.Println("SESTERTII UPDATE ERROR:", err)

		http.Error(
			w,
			"failed to update sestertii",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// COMMIT
	// =====================================================

	err = tx.Commit(r.Context())

	if err != nil {

		log.Println("TRAINING TRANSACTION COMMIT ERROR:", err)

		http.Error(
			w,
			"failed to save training attempt",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// RESPONSE
	// =====================================================

	w.WriteHeader(http.StatusNoContent)
}