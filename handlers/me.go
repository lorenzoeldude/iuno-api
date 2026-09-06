package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"iuno-api/db"
	"iuno-api/middleware"
	"iuno-api/models"
	"iuno-api/utils"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {

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
	// GET AUTHENTICATED USER FROM CONTEXT
	// =====================================================

	claims, ok := r.Context().Value(
		middleware.UserContextKey,
	).(*utils.Claims)

	if !ok || claims == nil {

		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)

		return
	}

	// =====================================================
	// FIND USER
	// =====================================================

	var user models.User

	err := db.Pool.QueryRow(context.Background(), `
		SELECT
			id,
			email,
			username,
			is_premium,
			is_admin,
			email_verified,
			created_at
		FROM users
		WHERE id = $1
	`,
		claims.UserID,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.IsPremium,
		&user.IsAdmin,
		&user.EmailVerified,
		&user.CreatedAt,
	)

	if err != nil {

		http.Error(
			w,
			"user not found",
			http.StatusUnauthorized,
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

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":             user.ID,
		"email":          user.Email,
		"username":       user.Username,
		"is_premium":     user.IsPremium,
		"is_admin":       user.IsAdmin,
		"email_verified": user.EmailVerified,
	})
}
