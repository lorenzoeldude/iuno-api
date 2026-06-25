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

		authHeader := r.Header.Get("Authorization")

		// log.Println("AUTH HEADER:", authHeader)

		if authHeader == "" {
			log.Println("AUTH FAILED: missing Authorization header")
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// log.Println("TOKEN STRING:", tokenString)

		claims := &utils.Claims{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				return utils.JwtSecret, nil
			},
		)

		// log.Println("PARSE ERROR:", err)

		// if token != nil {
		// 	log.Println("TOKEN VALID:", token.Valid)
		// } else {
		// 	log.Println("TOKEN IS NIL")
		// }

		// log.Printf("CLAIMS AFTER PARSE: %+v\n", claims)

		if err != nil || token == nil || !token.Valid {
			log.Println("AUTH FAILED: invalid token")
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// log.Println("AUTH SUCCESS - USER ID:", claims.UserID)

		ctx := context.WithValue(
			r.Context(),
			UserContextKey,
			claims,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}