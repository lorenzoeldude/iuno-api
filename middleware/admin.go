package middleware

import (
	"net/http"

	"iuno-api/utils"
)

func AdminOnly(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		claims, ok := r.Context().Value(UserContextKey).(*utils.Claims)

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !claims.IsAdmin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}