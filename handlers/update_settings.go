package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"iuno-api/db"
	"iuno-api/middleware"
	"iuno-api/utils"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UpdateSettingsRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func UpdateSettingsHandler(w http.ResponseWriter, r *http.Request) {

	// =====================================================
	// METHOD CHECK
	// =====================================================
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// =====================================================
	// AUTH
	// =====================================================
	claimsRaw := r.Context().Value(middleware.UserContextKey)

	claims, ok := claimsRaw.(*utils.Claims)

	if !ok || claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID

	// =====================================================
	// PARSE BODY
	// =====================================================
	var req UpdateSettingsRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// =====================================================
	// VALIDATION
	// =====================================================
	if req.Username == "" {
		http.Error(w, "username required", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "email required", http.StatusBadRequest)
		return
	}

	// =====================================================
	// HASH PASSWORD (ONLY IF PROVIDED)
	// =====================================================
	var hashedPassword string

	if req.Password != "" {

		hash, err := bcrypt.GenerateFromPassword(
			[]byte(req.Password),
			bcrypt.DefaultCost,
		)

		if err != nil {

			log.Println("HASH ERROR:", err)

			http.Error(
				w,
				"failed to hash password",
				http.StatusInternalServerError,
			)

			return
		}

		hashedPassword = string(hash)
	}

	// =====================================================
	// UPDATE USER
	// =====================================================
	if req.Password != "" {

		_, err = db.Pool.Exec(r.Context(), `
			UPDATE users
			SET
				username = $1,
				email = $2,
				password_hash = $3
			WHERE id = $4
		`,
			req.Username,
			req.Email,
			hashedPassword,
			userID,
		)

	} else {

		_, err = db.Pool.Exec(r.Context(), `
			UPDATE users
			SET
				username = $1,
				email = $2
			WHERE id = $3
		`,
			req.Username,
			req.Email,
			userID,
		)
	}

	if err != nil {

		log.Println("UPDATE SETTINGS ERROR:", err)

		http.Error(
			w,
			"failed to update settings",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// CREATE NEW TOKEN
	// =====================================================
	newClaims := utils.Claims{
		UserID:    userID,
		Username:  req.Username,
		IsPremium: claims.IsPremium,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(7 * 24 * time.Hour),
			),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		newClaims,
	)

	tokenString, err := token.SignedString(utils.JwtSecret)

	if err != nil {

		log.Println("TOKEN SIGN ERROR:", err)

		http.Error(
			w,
			"failed to create token",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// RESPONSE
	// =====================================================
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": tokenString,
		"user": map[string]interface{}{
			"id":       userID,
			"username": req.Username,
			"email":    req.Email,
		},
	})
}