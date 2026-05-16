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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

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
	var body models.RegisterRequest

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
	body.Username = strings.TrimSpace(strings.ToLower(body.Username))

	// =====================================================
	// VALIDATION
	// =====================================================
	if body.Email == "" ||
		body.Username == "" ||
		body.Password == "" {

		http.Error(
			w,
			"missing required fields",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================================
	// HASH PASSWORD
	// =====================================================
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(body.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {

		log.Println("BCRYPT ERROR:", err)

		http.Error(
			w,
			"failed to create account",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// INSERT USER
	// =====================================================
	var userID int

	err = db.Pool.QueryRow(context.Background(), `
		INSERT INTO users (
			email,
			username,
			password_hash
		)
		VALUES ($1, $2, $3)
		RETURNING id
	`,
		body.Email,
		body.Username,
		string(passwordHash),
	).Scan(&userID)

	if err != nil {

		log.Println("REGISTER ERROR:", err)

		http.Error(
			w,
			"email or username already exists",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================================
	// CREATE DEFAULT WORD LIST
	// =====================================================
	_, err = db.Pool.Exec(context.Background(), `
		INSERT INTO word_lists (
			user_id,
			name
		)
		VALUES ($1, $2)
	`,
		userID,
		"My Vocabulary",
	)

	if err != nil {

		log.Println("WORDLIST CREATE ERROR:", err)

		http.Error(
			w,
			"failed to create default word list",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// RESPONSE
	// =====================================================
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "ok",
		"user_id":  userID,
		"username": body.Username,
	})
}