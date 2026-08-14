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
	// STATS STRUCT
	// =====================================================

	var stats struct {
		UserID            int     `json:"userId"`
		QuestionsAnswered int     `json:"questionsAnswered"`
		CorrectAnswers    int     `json:"correctAnswers"`
		IncorrectAnswers  int     `json:"incorrectAnswers"`
		TrainingSessions  int     `json:"trainingSessions"`
		QuestionsToday    int     `json:"questionsToday"`
		CurrentStreak     int     `json:"currentStreak"`
		LongestStreak     int     `json:"longestStreak"`
		Sestertii         int     `json:"sestertii"`
		LastTrainedAt     *string `json:"lastTrainedAt"`
	}

	stats.UserID = userID

	// =====================================================
	// GET AGGREGATE STATS
	// =====================================================

	err := db.Pool.QueryRow(
		r.Context(),
		`
		SELECT
			questions_answered,
			correct_answers,
			incorrect_answers,
			training_sessions,
			TO_CHAR(
				last_trained_at,
				'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"'
			)
		FROM user_training_stats
		WHERE user_id = $1
		`,
		userID,
	).Scan(
		&stats.QuestionsAnswered,
		&stats.CorrectAnswers,
		&stats.IncorrectAnswers,
		&stats.TrainingSessions,
		&stats.LastTrainedAt,
	)

	// =====================================================
	// NO STATS YET
	// =====================================================

	if err != nil && err != pgx.ErrNoRows {

		log.Println("GET TRAINING STATS ERROR:", err)

		http.Error(
			w,
			"failed to get training stats",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// GET SESTERTII
	// =====================================================

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

	if err != nil {

		log.Println("GET SESTERTII ERROR:", err)

		http.Error(
			w,
			"failed to get user sestertii",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// QUESTIONS ANSWERED TODAY
	//
	// Uses Vietnam local time.
	// =====================================================

	err = db.Pool.QueryRow(
		r.Context(),
		`
		SELECT COUNT(*)
		FROM training_attempts
		WHERE user_id = $1
		  AND (
			created_at AT TIME ZONE 'Asia/Ho_Chi_Minh'
		  )::date = (
			CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Ho_Chi_Minh'
		  )::date
		`,
		userID,
	).Scan(
		&stats.QuestionsToday,
	)

	if err != nil {

		log.Println("GET QUESTIONS TODAY ERROR:", err)

		http.Error(
			w,
			"failed to get today's training stats",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// CURRENT STREAK
	//
	// A streak is consecutive days on which the user
	// completed at least one training question.
	//
	// The streak remains active if the user trained:
	//
	// - today
	// OR
	// - yesterday
	//
	// If neither happened, the streak is 0.
	// =====================================================

	err = db.Pool.QueryRow(
		r.Context(),
		`
		WITH training_days AS (

			SELECT DISTINCT
				(
					created_at
					AT TIME ZONE 'Asia/Ho_Chi_Minh'
				)::date AS training_date

			FROM training_attempts

			WHERE user_id = $1
		),

		ordered_days AS (

			SELECT
				training_date,

				training_date
				- ROW_NUMBER() OVER (
					ORDER BY training_date
				)::int AS streak_group

			FROM training_days
		),

		streaks AS (

			SELECT
				streak_group,
				COUNT(*) AS streak_length,
				MAX(training_date) AS last_day

			FROM ordered_days

			GROUP BY streak_group
		)

		SELECT COALESCE(
			(
				SELECT streak_length

				FROM streaks

				WHERE last_day >= (
					CURRENT_TIMESTAMP
					AT TIME ZONE 'Asia/Ho_Chi_Minh'
				)::date - 1

				ORDER BY last_day DESC

				LIMIT 1
			),
			0
		)
		`,
		userID,
	).Scan(
		&stats.CurrentStreak,
	)

	if err != nil {

		log.Println("GET CURRENT STREAK ERROR:", err)

		http.Error(
			w,
			"failed to read training streak",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// LONGEST STREAK
	// =====================================================

	err = db.Pool.QueryRow(
		r.Context(),
		`
		WITH training_days AS (

			SELECT DISTINCT
				(
					created_at
					AT TIME ZONE 'Asia/Ho_Chi_Minh'
				)::date AS training_date

			FROM training_attempts

			WHERE user_id = $1
		),

		ordered_days AS (

			SELECT
				training_date,

				training_date
				- ROW_NUMBER() OVER (
					ORDER BY training_date
				)::int AS streak_group

			FROM training_days
		),

		streaks AS (

			SELECT
				streak_group,
				COUNT(*) AS streak_length

			FROM ordered_days

			GROUP BY streak_group
		)

		SELECT COALESCE(
			MAX(streak_length),
			0
		)

		FROM streaks
		`,
		userID,
	).Scan(
		&stats.LongestStreak,
	)

	if err != nil {

		log.Println("GET LONGEST STREAK ERROR:", err)

		http.Error(
			w,
			"failed to read longest training streak",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// RESPONSE
	// =====================================================

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

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

		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	// =====================================================
	// AUTH
	// =====================================================

	claimsRaw := r.Context().Value(middleware.UserContextKey)

	claims, ok := claimsRaw.(*utils.Claims)

	if !ok || claims == nil {

		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)

		return
	}

	userID := claims.UserID

	// =====================================================
	// PARSE BODY
	// =====================================================

	var req RecordTrainingAttemptRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {

		http.Error(
			w,
			"invalid json",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================================
	// DATABASE TRANSACTION
	// =====================================================

	tx, err := db.Pool.Begin(r.Context())

	if err != nil {

		log.Println(
			"TRAINING TRANSACTION ERROR:",
			err,
		)

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

		log.Println(
			"TRAINING ATTEMPT INSERT ERROR:",
			err,
		)

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

		log.Println(
			"TRAINING STATS CREATE ERROR:",
			err,
		)

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
				questions_answered =
					questions_answered + 1,

				correct_answers =
					correct_answers + 1,

				last_trained_at =
					NOW(),

				updated_at =
					NOW()

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
				questions_answered =
					questions_answered + 1,

				incorrect_answers =
					incorrect_answers + 1,

				last_trained_at =
					NOW(),

				updated_at =
					NOW()

			WHERE user_id = $1
			`,
			userID,
		)
	}

	if err != nil {

		log.Println(
			"TRAINING STATS UPDATE ERROR:",
			err,
		)

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

		log.Println(
			"SESTERTII UPDATE ERROR:",
			err,
		)

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

		log.Println(
			"TRAINING TRANSACTION COMMIT ERROR:",
			err,
		)

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