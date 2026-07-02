package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"iuno-api/db"
	"iuno-api/models"
	"iuno-api/utils"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {

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
	// PARSE REQUEST
	// =====================================================
	var req models.LoginRequest

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
	// FIND USER
	// =====================================================
	var user models.User

	err = db.Pool.QueryRow(context.Background(), `
		SELECT
			id,
			email,
			username,
			password_hash,
			is_premium,
			is_admin,
			email_verified,
			created_at
		FROM users
		WHERE email = $1
	`,
		req.Email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.IsPremium,
		&user.IsAdmin,
		&user.EmailVerified,
		&user.CreatedAt,
	)

	if err != nil {

		http.Error(
			w,
			"invalid email or password",
			http.StatusUnauthorized,
		)

		return
	}

	// =====================================================
	// CHECK PASSWORD
	// =====================================================
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)

	if err != nil {

		http.Error(
			w,
			"invalid email or password",
			http.StatusUnauthorized,
		)

		return
	}

	// =====================================================
	// CHECK EMAIL VERIFIED
	// =====================================================
	if !user.EmailVerified {

		http.Error(
			w,
			"Please verify your email before logging in.",
			http.StatusForbidden,
		)

		return
	}

	// =====================================================
	// GENERATE JWT
	// =====================================================
	jwtToken, err := utils.GenerateJWT(
		user.ID,
		user.Username,
		user.IsPremium,
		user.IsAdmin,
	)

	if err != nil {

		http.Error(
			w,
			"failed to generate token",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// RESPONSE
	// =====================================================
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"token":  jwtToken,

		"user": map[string]interface{}{
			"id":         user.ID,
			"email":      user.Email,
			"username":   user.Username,
			"is_premium": user.IsPremium,
			"is_admin": user.IsAdmin,
			"email_verified": user.EmailVerified,
		},
	})
}