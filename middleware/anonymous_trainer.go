package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"iuno-api/db"
	"iuno-api/utils"
)

type anonymousTrainerContextKey string

const AnonymousTrainerIDKey = anonymousTrainerContextKey("anonymous_trainer_id")

const AnonymousDailyLimit = 10

func AnonymousTrainerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// =====================================================
		// LOGGED-IN USERS
		// =====================================================

		authHeader := r.Header.Get("Authorization")

		if authHeader != "" {

			tokenString := strings.TrimPrefix(
				authHeader,
				"Bearer ",
			)

			claims := &utils.Claims{}

			token, err := jwt.ParseWithClaims(
				tokenString,
				claims,
				func(token *jwt.Token) (interface{}, error) {
					return utils.JwtSecret, nil
				},
			)

			// Valid authenticated user.
			// They do NOT use the anonymous quota.
			if err == nil &&
				token != nil &&
				token.Valid {

				ctx := context.WithValue(
					r.Context(),
					UserContextKey,
					claims,
				)

				next.ServeHTTP(
					w,
					r.WithContext(ctx),
				)

				return
			}
		}

		// =====================================================
		// GET ANONYMOUS ID COOKIE
		// =====================================================

		cookie, err := r.Cookie("anonymous_id")

		if err != nil || cookie.Value == "" {

			anonymousID := uuid.New().String()

			http.SetCookie(w, &http.Cookie{
				Name:     "anonymous_id",
				Value:    anonymousID,
				Path:     "/",
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
				Secure:   true,
				MaxAge:   60 * 60 * 24 * 365,
			})

			cookie = &http.Cookie{
				Value: anonymousID,
			}
		}

		anonymousID := cookie.Value

		// =====================================================
		// VALIDATE UUID
		// =====================================================

		_, err = uuid.Parse(anonymousID)

		if err != nil {

			http.Error(
				w,
				"invalid anonymous id",
				http.StatusBadRequest,
			)

			return
		}

		// =====================================================
		// TODAY
		// =====================================================

		today := time.Now().UTC().Format("2006-01-02")

		// =====================================================
		// GET DAILY USAGE
		// =====================================================

		var questionsAnswered int

		err = db.Pool.QueryRow(
			context.Background(),
			`
			SELECT questions_answered
			FROM anonymous_trainer_daily_usage
			WHERE anonymous_id = $1
			AND usage_date = $2
			`,
			anonymousID,
			today,
		).Scan(&questionsAnswered)

		// No record yet = zero usage.
		if err != nil {
			questionsAnswered = 0
		}

		// =====================================================
		// DAILY LIMIT
		// =====================================================

		if questionsAnswered >= AnonymousDailyLimit {

			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			w.WriteHeader(http.StatusTooManyRequests)

			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":       "daily anonymous training limit reached",
				"limit":       AnonymousDailyLimit,
				"used":        questionsAnswered,
				"limitReached": true,
			})

			return
		}

		// =====================================================
		// STORE ANONYMOUS ID IN CONTEXT
		// =====================================================

		ctx := context.WithValue(
			r.Context(),
			AnonymousTrainerIDKey,
			anonymousID,
		)

		// =====================================================
		// CONTINUE
		// =====================================================

		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)

		// =====================================================
		// RECORD QUESTION
		// =====================================================

		// The request successfully reached the trainer handler,
		// so count it as one anonymous training question.
		_, err = db.Pool.Exec(
			context.Background(),
			`
			INSERT INTO anonymous_trainer_daily_usage (
				anonymous_id,
				usage_date,
				questions_answered
			)
			VALUES ($1, $2, 1)
			ON CONFLICT (anonymous_id, usage_date)
			DO UPDATE SET
				questions_answered =
					anonymous_trainer_daily_usage.questions_answered + 1
			`,
			anonymousID,
			today,
		)

		if err != nil {
			// The question was already served, so don't break
			// the response. Just log the database failure.
			// You can add proper logging later.
		}
	}
}