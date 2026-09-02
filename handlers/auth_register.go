//
//  AuthHandler.go
//  IUNO API
//

package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"iuno-api/db"
	"iuno-api/email"
	"iuno-api/models"
	"iuno-api/utils"
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)

// =====================================================
// MARK: - REGISTER
// =====================================================

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

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

	body.Email = strings.TrimSpace(
		strings.ToLower(body.Email),
	)

	// Don't TrimSpace the username.
	// Spaces should be rejected, not silently removed.
	body.Username = strings.ToLower(body.Username)

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
	// USERNAME VALIDATION
	// =====================================================

	if strings.ContainsAny(
		body.Username,
		" \t\n\r",
	) {

		http.Error(
			w,
			"no spaces allowed",
			http.StatusBadRequest,
		)

		return
	}

	if !usernameRegex.MatchString(
		body.Username,
	) {

		http.Error(
			w,
			"username must be 3-20 characters and contain only letters, numbers, and underscores",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================================
	// CHECK IF EMAIL ALREADY EXISTS
	// =====================================================

	var emailExists bool

	err = db.Pool.QueryRow(
		context.Background(),
		`
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE email = $1
		)
		`,
		body.Email,
	).Scan(&emailExists)

	if err != nil {

		log.Println(
			"EMAIL CHECK ERROR:",
			err,
		)

		http.Error(
			w,
			"failed to create account",
			http.StatusInternalServerError,
		)

		return
	}

	if emailExists {

		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		w.WriteHeader(
			http.StatusConflict,
		)

		json.NewEncoder(w).Encode(
			map[string]string{
				"error":
					"email already exists",
			},
		)

		return
	}

	// =====================================================
	// CHECK IF USERNAME ALREADY EXISTS
	// =====================================================

	var usernameExists bool

	err = db.Pool.QueryRow(
		context.Background(),
		`
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE username = $1
		)
		`,
		body.Username,
	).Scan(&usernameExists)

	if err != nil {

		log.Println(
			"USERNAME CHECK ERROR:",
			err,
		)

		http.Error(
			w,
			"failed to create account",
			http.StatusInternalServerError,
		)

		return
	}

	if usernameExists {

		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		w.WriteHeader(
			http.StatusConflict,
		)

		json.NewEncoder(w).Encode(
			map[string]string{
				"error":
					"username already exists",
			},
		)

		return
	}

	// =====================================================
	// HASH PASSWORD
	// =====================================================

	passwordHash, err :=
		bcrypt.GenerateFromPassword(
			[]byte(body.Password),
			bcrypt.DefaultCost,
		)

	if err != nil {

		log.Println(
			"BCRYPT ERROR:",
			err,
		)

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

	verificationToken, err :=
		utils.GenerateVerificationToken()

	if err != nil {

		log.Println(
			"TOKEN ERROR:",
			err,
		)

		http.Error(
			w,
			"failed to create account",
			http.StatusInternalServerError,
		)

		return
	}

	verificationHash :=
		utils.HashVerificationToken(
			verificationToken,
		)

	// =====================================================
	// INSERT USER
	// =====================================================

	var userID int

	err = db.Pool.QueryRow(
		context.Background(),
		`
		INSERT INTO users (
			email,
			username,
			password_hash,
			email_verified,
			email_verification_hash,
			email_verification_expires_at,
			email_verification_last_sent_at
		)
		VALUES (
			$1,
			$2,
			$3,
			FALSE,
			$4,
			NOW() + INTERVAL '24 hours',
			NOW()
		)
		RETURNING id
		`,
		body.Email,
		body.Username,
		string(passwordHash),
		verificationHash,
	).Scan(&userID)

	if err != nil {

		log.Println(
			"REGISTER ERROR:",
			err,
		)

		// This can still happen if two requests
		// try to register the same email/username
		// simultaneously and the database has UNIQUE
		// constraints on those columns.
		http.Error(
			w,
			"email or username already exists",
			http.StatusConflict,
		)

		return
	}

	// =====================================================
	// CREATE DEFAULT WORD LIST
	// =====================================================

	_, err = db.Pool.Exec(
		context.Background(),
		`
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

		log.Println(
			"WORDLIST CREATE ERROR:",
			err,
		)

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

		log.Println(
			"EMAIL ERROR:",
			err,
		)

		// Don't fail the registration if the email
		// couldn't be sent.
		// The user can request another verification
		// email later.
	}

	// =====================================================
	// RESPONSE
	// =====================================================

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"status":   "ok",
			"user_id":  userID,
			"username": body.Username,
		},
	)
}


// =====================================================
// MARK: - RESEND VERIFICATION EMAIL
// =====================================================

func ResendVerificationHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

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

	var body struct {
		Email string `json:"email"`
	}

	err := json.NewDecoder(
		r.Body,
	).Decode(&body)

	if err != nil {

		log.Println(
			"RESEND VERIFICATION JSON ERROR:",
			err,
		)

		http.Error(
			w,
			"invalid json",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================================
	// NORMALIZE EMAIL
	// =====================================================

	body.Email = strings.TrimSpace(
		strings.ToLower(body.Email),
	)

	if body.Email == "" {

		http.Error(
			w,
			"missing email",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================================
	// FIND USER
	// =====================================================

	var (
		userID     int
		verified   bool
		lastSentAt *time.Time
	)

	err = db.Pool.QueryRow(
		context.Background(),
		`
		SELECT
			id,
			email_verified,
			email_verification_last_sent_at
		FROM users
		WHERE email = $1
		`,
		body.Email,
	).Scan(
		&userID,
		&verified,
		&lastSentAt,
	)

	// =====================================================
	// GENERIC RESPONSE FOR UNKNOWN EMAIL
	// =====================================================

	if err != nil {

		// Don't reveal whether an account exists.
		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		json.NewEncoder(w).Encode(
			map[string]string{
				"status": "ok",
			},
		)

		return
	}

	// =====================================================
	// ALREADY VERIFIED
	// =====================================================

	if verified {

		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		json.NewEncoder(w).Encode(
			map[string]string{
				"status": "ok",
			},
		)

		return
	}

	// =====================================================
	// RATE LIMIT
	// =====================================================

	if lastSentAt != nil {

		if time.Since(*lastSentAt) <
			time.Minute {

			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			w.WriteHeader(
				http.StatusTooManyRequests,
			)

			json.NewEncoder(w).Encode(
				map[string]string{
					"error":
						"please wait before requesting another verification email",
				},
			)

			return
		}
	}

	// =====================================================
	// GENERATE NEW VERIFICATION TOKEN
	// =====================================================

	verificationToken, err :=
		utils.GenerateVerificationToken()

	if err != nil {

		log.Println(
			"RESEND VERIFICATION TOKEN ERROR:",
			err,
		)

		http.Error(
			w,
			"failed to resend verification email",
			http.StatusInternalServerError,
		)

		return
	}

	verificationHash :=
		utils.HashVerificationToken(
			verificationToken,
		)

	// =====================================================
	// UPDATE VERIFICATION TOKEN
	// =====================================================

	_, err = db.Pool.Exec(
		context.Background(),
		`
		UPDATE users
		SET
			email_verification_hash = $1,
			email_verification_expires_at =
				NOW() + INTERVAL '24 hours',
			email_verification_last_sent_at = NOW()
		WHERE id = $2
		`,
		verificationHash,
		userID,
	)

	if err != nil {

		log.Println(
			"RESEND VERIFICATION UPDATE ERROR:",
			err,
		)

		http.Error(
			w,
			"failed to resend verification email",
			http.StatusInternalServerError,
		)

		return
	}

	// =====================================================
	// SEND EMAIL
	// =====================================================

	err = email.SendVerificationEmail(
		body.Email,
		verificationToken,
	)

	if err != nil {

		log.Println(
			"RESEND VERIFICATION EMAIL ERROR:",
			err,
		)

		http.Error(
			w,
			"failed to send verification email",
			http.StatusInternalServerError,
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

	json.NewEncoder(w).Encode(
		map[string]string{
			"status": "ok",
		},
	)
}