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

	utils.EnableCORS(w)

	// =====================================================
	// CORS PRE-FLIGHT
	// =====================================================
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

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
	// GENERATE JWT
	// =====================================================
	jwtToken, err := utils.GenerateJWT(
		user.ID,
		user.Username,
		user.IsPremium,
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
		},
	})
}