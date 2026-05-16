package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

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
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// =====================================================
	// PARSE REQUEST
	// =====================================================
	var body models.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {

		log.Println("JSON ERROR:", err)

		http.Error(
			w,
			"invalid json",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================================
	// NORMALIZE INPUT
	// =====================================================
	body.Email = strings.TrimSpace(strings.ToLower(body.Email))

	// =====================================================
	// VALIDATION
	// =====================================================
	if body.Email == "" || body.Password == "" {

		http.Error(
			w,
			"email and password required",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================================
	// FETCH USER
	// =====================================================
	var userID int
	var username string
	var passwordHash string
	var isPremium bool

	err = db.Pool.QueryRow(context.Background(), `
		SELECT
			id,
			username,
			password_hash,
			is_premium
		FROM users
		WHERE email = $1
	`,
		body.Email,
	).Scan(
		&userID,
		&username,
		&passwordHash,
		&isPremium,
	)

	if err != nil {

		log.Println("LOGIN QUERY ERROR:", err)

		http.Error(
			w,
			"invalid credentials",
			http.StatusUnauthorized,
		)

		return
	}

	// =====================================================
	// VERIFY PASSWORD
	// =====================================================
	err = bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(body.Password),
	)

	if err != nil {

		http.Error(
			w,
			"invalid credentials",
			http.StatusUnauthorized,
		)

		return
	}

	// =====================================================
	// RESPONSE
	// =====================================================
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "ok",
		"user_id":     userID,
		"username":    username,
		"is_premium":  isPremium,
	})
}