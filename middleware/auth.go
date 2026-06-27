package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"iuno-api/utils"
)

type contextKey string

const UserContextKey = contextKey("user")

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// =====================================================
		// CORS PREFLIGHT
		// =====================================================
		// OPTIONS requests do not contain Authorization headers.
		// They must pass through so CORSMiddleware can respond.
		// =====================================================
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		// =====================================================
		// GET AUTHORIZATION HEADER
		// =====================================================
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			log.Println("AUTH FAILED: missing Authorization header")
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// =====================================================
		// PARSE JWT
		// =====================================================
		claims := &utils.Claims{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				return utils.JwtSecret, nil
			},
		)

		if err != nil || token == nil || !token.Valid {
			log.Println("AUTH FAILED: invalid token")
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// =====================================================
		// STORE USER IN CONTEXT
		// =====================================================
		ctx := context.WithValue(
			r.Context(),
			UserContextKey,
			claims,
		)

		// =====================================================
		// CONTINUE
		// =====================================================
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}