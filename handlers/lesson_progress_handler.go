package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"iuno-api/db"
	"iuno-api/middleware"
	"iuno-api/models"
	"iuno-api/utils"

	"github.com/jackc/pgx/v5"
)

// =====================================================
// GET LESSON PROGRESS
// GET /api/lessons/{id}/progress
// =====================================================

func GetLessonProgressHandler(w http.ResponseWriter, r *http.Request) {

	// =====================================================
	// METHOD CHECK
	// =====================================================

	if r.Method != http.MethodGet {
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
	// PARSE LESSON ID
	// =====================================================

	idString := strings.TrimPrefix(
		r.URL.Path,
		"/api/lessons/",
	)

	idString = strings.TrimSuffix(
		idString,
		"/progress",
	)

	lessonID, err := strconv.Atoi(idString)

	if err != nil {
		http.Error(
			w,
			"invalid lesson id",
			http.StatusBadRequest,
		)
		return
	}

	// =====================================================
	// GET PROGRESS
	// =====================================================

	var progress models.UserLessonProgress

	err = db.Pool.QueryRow(
		r.Context(),
		`
		SELECT
			user_id,
			lesson_id,
			text_completed,
			vocabulary_completed,
			grammar_completed,
			examinatio_completed,
			score,
			started_at,
			completed_at
		FROM user_lesson_progress
		WHERE user_id = $1
		AND lesson_id = $2
		`,
		userID,
		lessonID,
	).Scan(
		&progress.UserID,
		&progress.LessonID,
		&progress.TextCompleted,
		&progress.VocabularyCompleted,
		&progress.GrammarCompleted,
		&progress.ExaminatioCompleted,
		&progress.Score,
		&progress.StartedAt,
		&progress.CompletedAt,
	)

	// =====================================================
	// NO PROGRESS YET
	// =====================================================

	if err == pgx.ErrNoRows {

		progress = models.UserLessonProgress{
			UserID:              userID,
			LessonID:            lessonID,
			TextCompleted:       false,
			VocabularyCompleted: false,
			GrammarCompleted:    false,
			ExaminatioCompleted: false,
			Score:               nil,
		}

		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		json.NewEncoder(w).Encode(progress)

		return
	}

	if err != nil {

		log.Println(
			"GET LESSON PROGRESS ERROR:",
			err,
		)

		http.Error(
			w,
			"failed to get lesson progress",
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

	json.NewEncoder(w).Encode(progress)
}

// =====================================================
// UPDATE LESSON PROGRESS
// PUT /api/lessons/{id}/progress
// =====================================================

type UpdateLessonProgressRequest struct {
	Section string `json:"section"`
	Score   *int   `json:"score,omitempty"`
}

func UpdateLessonProgressHandler(w http.ResponseWriter, r *http.Request) {

	// =====================================================
	// METHOD CHECK
	// =====================================================

	if r.Method != http.MethodPut {
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
	// PARSE LESSON ID
	// =====================================================

	idString := strings.TrimPrefix(
		r.URL.Path,
		"/api/lessons/",
	)

	idString = strings.TrimSuffix(
		idString,
		"/progress",
	)

	lessonID, err := strconv.Atoi(idString)

	if err != nil {
		http.Error(
			w,
			"invalid lesson id",
			http.StatusBadRequest,
		)
		return
	}

	// =====================================================
	// PARSE BODY
	// =====================================================

	var req UpdateLessonProgressRequest

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(
			w,
			"invalid json",
			http.StatusBadRequest,
		)
		return
	}

	// =====================================================
	// VALIDATE SECTION
	// =====================================================

	switch req.Section {

	case "text":
	case "vocabulary":
	case "grammar":
	case "examinatio":

	default:
		http.Error(
			w,
			"invalid section",
			http.StatusBadRequest,
		)
		return
	}

	// =====================================================
	// VALIDATE EXAMINATIO SCORE
	// =====================================================

	if req.Section == "examinatio" {

		if req.Score == nil {
			http.Error(
				w,
				"score is required for examinatio",
				http.StatusBadRequest,
			)
			return
		}

		if *req.Score < 0 || *req.Score > 100 {
			http.Error(
				w,
				"score must be between 0 and 100",
				http.StatusBadRequest,
			)
			return
		}
	}

	// =====================================================
	// DATABASE TRANSACTION
	// =====================================================

	tx, err := db.Pool.Begin(r.Context())

	if err != nil {

		log.Println(
			"LESSON PROGRESS TRANSACTION ERROR:",
			err,
		)

		http.Error(
			w,
			"failed to update lesson progress",
			http.StatusInternalServerError,
		)

		return
	}

	defer tx.Rollback(r.Context())

	// =====================================================
	// MAKE SURE PROGRESS ROW EXISTS
	// =====================================================

	_, err = tx.Exec(
		r.Context(),
		`
		INSERT INTO user_lesson_progress (
			user_id,
			lesson_id,
			started_at
		)
		VALUES (
			$1,
			$2,
			NOW()
		)
		ON CONFLICT (user_id, lesson_id)
		DO NOTHING
		`,
		userID,
		lessonID,
	)

	if err != nil {

		log.Println(
			"LESSON PROGRESS CREATE ERROR:",
			err,
		)

		http.Error(
			w,
			"failed to create lesson progress",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// CHECK WHETHER EXAMINATIO WAS ALREADY COMPLETED
	// =====================================================

	var examinatioAlreadyCompleted bool

	err = tx.QueryRow(
		r.Context(),
		`
		SELECT examinatio_completed
		FROM user_lesson_progress
		WHERE user_id = $1
		AND lesson_id = $2
		FOR UPDATE
		`,
		userID,
		lessonID,
	).Scan(
		&examinatioAlreadyCompleted,
	)

	if err != nil {

		log.Println(
			"LESSON PROGRESS LOCK ERROR:",
			err,
		)

		http.Error(
			w,
			"failed to read lesson progress",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// UPDATE SECTION
	// =====================================================

	if req.Section == "examinatio" {

		// =================================================
		// EXAMINATIO
		// =================================================

		_, err = tx.Exec(
			r.Context(),
			`
			UPDATE user_lesson_progress
			SET
				examinatio_completed = TRUE,
				score = GREATEST(
					COALESCE(score, 0),
					$3
				),
				completed_at = COALESCE(
					completed_at,
					NOW()
				)
			WHERE user_id = $1
			AND lesson_id = $2
			`,
			userID,
			lessonID,
			*req.Score,
		)

	} else {

		// =================================================
		// OTHER SECTIONS
		// =================================================

		var updateQuery string

		switch req.Section {

		case "text":

			updateQuery = `
				UPDATE user_lesson_progress
				SET text_completed = TRUE
				WHERE user_id = $1
				AND lesson_id = $2
			`

		case "vocabulary":

			updateQuery = `
				UPDATE user_lesson_progress
				SET vocabulary_completed = TRUE
				WHERE user_id = $1
				AND lesson_id = $2
			`

		case "grammar":

			updateQuery = `
				UPDATE user_lesson_progress
				SET grammar_completed = TRUE
				WHERE user_id = $1
				AND lesson_id = $2
			`
		}

		_, err = tx.Exec(
			r.Context(),
			updateQuery,
			userID,
			lessonID,
		)
	}

	// =====================================================
	// CHECK UPDATE ERROR
	// =====================================================

	if err != nil {

		log.Println(
			"LESSON PROGRESS UPDATE ERROR:",
			err,
		)

		http.Error(
			w,
			"failed to update lesson progress",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// AWARD 100 SESTERTII FOR LESSON COMPLETION
	// =====================================================

	if req.Section == "examinatio" &&
		!examinatioAlreadyCompleted {

		_, err = tx.Exec(
			r.Context(),
			`
			UPDATE users
			SET sestertii = sestertii + 100
			WHERE id = $1
			`,
			userID,
		)

		if err != nil {

			log.Println(
				"LESSON SESTERTII UPDATE ERROR:",
				err,
			)

			http.Error(
				w,
				"failed to award lesson sestertii",
				http.StatusInternalServerError,
			)

			return
		}
	}

	// =====================================================
	// AWARD 25 BONUS SESTERTII FOR PERFECT EXAMINATIO
	// =====================================================

	if req.Section == "examinatio" &&
		!examinatioAlreadyCompleted &&
		*req.Score == 100 {

		_, err = tx.Exec(
			r.Context(),
			`
			UPDATE users
			SET sestertii = sestertii + 25
			WHERE id = $1
			`,
			userID,
		)

		if err != nil {

			log.Println(
				"PERFECT LESSON SESTERTII UPDATE ERROR:",
				err,
			)

			http.Error(
				w,
				"failed to award perfect lesson sestertii",
				http.StatusInternalServerError,
			)

			return
		}
	}

	// =====================================================
	// COMMIT
	// =====================================================

	err = tx.Commit(r.Context())

	if err != nil {

		log.Println(
			"LESSON PROGRESS COMMIT ERROR:",
			err,
		)

		http.Error(
			w,
			"failed to save lesson progress",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// RESPONSE
	// =====================================================

	w.WriteHeader(http.StatusNoContent)
}