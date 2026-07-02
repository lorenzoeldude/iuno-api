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
	"iuno-api/email"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

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
	// CREATE EMAIL VERIFICATION TOKEN
	// =====================================================
	verificationToken, err := utils.GenerateVerificationToken()
	if err != nil {

		log.Println("TOKEN ERROR:", err)

		http.Error(
			w,
			"failed to create account",
			http.StatusInternalServerError,
		)

		return
	}

	verificationHash := utils.HashVerificationToken(verificationToken)

	// =====================================================
	// INSERT USER
	// =====================================================
	var userID int

	err = db.Pool.QueryRow(context.Background(), `
		INSERT INTO users (
			email,
			username,
			password_hash,
			email_verified,
			email_verification_hash,
			email_verification_expires_at
		)
		VALUES (
			$1,
			$2,
			$3,
			FALSE,
			$4,
			NOW() + INTERVAL '24 hours'
		)
		RETURNING id
	`,
		body.Email,
		body.Username,
		string(passwordHash),
		verificationHash,
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
	// SEND VERIFICATION EMAIL
	// =====================================================
	err = email.SendVerificationEmail(
		body.Email,
		verificationToken,
	)

	if err != nil {

		log.Println("EMAIL ERROR:", err)

		// Don't fail the registration if the email couldn't be sent.
		// The user can request another verification email later.
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