package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"iuno-api/utils"
)

type anonymousLookupContextKey string

const AnonymousLookupIDKey = anonymousLookupContextKey("anonymous_lookup_id")

func AnonymousLookupMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// =====================================================
		// LOGGED-IN USERS
		// =====================================================

		authHeader := r.Header.Get("Authorization")

		if authHeader != "" {

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			claims := &utils.Claims{}

			token, err := jwt.ParseWithClaims(
				tokenString,
				claims,
				func(token *jwt.Token) (interface{}, error) {
					return utils.JwtSecret, nil
				},
			)

			if err == nil && token != nil && token.Valid {

				ctx := context.WithValue(
					r.Context(),
					UserContextKey,
					claims,
				)

				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		// =====================================================
		// ANONYMOUS USERS
		// =====================================================

		anonymousID := r.Header.Get("X-Anonymous-ID")

		if anonymousID == "" {
			http.Error(
				w,
				"missing anonymous id",
				http.StatusBadRequest,
			)
			return
		}

		// Validate that the ID is actually a UUID.
		if _, err := uuid.Parse(anonymousID); err != nil {
			http.Error(
				w,
				"invalid anonymous id",
				http.StatusBadRequest,
			)
			return
		}

		ctx := context.WithValue(
			r.Context(),
			AnonymousLookupIDKey,
			anonymousID,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
