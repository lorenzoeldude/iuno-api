package middleware

import (
	"context"
	"net/http"
	"strings"
	"log"

	"github.com/golang-jwt/jwt/v5"

	"iuno-api/utils"
)

type contextKey string

const UserContextKey = contextKey("user")

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &utils.Claims{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				return utils.JwtSecret, nil
			},
		)

		if err != nil || !token.Valid {
			log.Println("TOKEN ERROR:", err)
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// log.Println("AUTH SUCCESS - USER ID:", claims.UserID)
		// log.Println("TOKEN:", tokenString)
		// log.Println("TOKEN VALID:", token.Valid)
		// log.Printf("CLAIMS: %+v\n", claims)

		ctx := context.WithValue(
			r.Context(),
			UserContextKey,
			claims,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}