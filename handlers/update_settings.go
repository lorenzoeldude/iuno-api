package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"iuno-api/db"
	"iuno-api/middleware"
	"iuno-api/utils"

	"golang.org/x/crypto/bcrypt"
)

type UpdateSettingsRequest struct {
	Username        string `json:"username"`
	CurrentPassword string `json:"current_password"`
	Password        string `json:"password"`
}

func UpdateSettingsHandler(w http.ResponseWriter, r *http.Request) {

	// =====================================================
	// METHOD CHECK
	// =====================================================

	if r.Method != http.MethodPut {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	// =====================================================
	// AUTH
	// =====================================================

	claimsRaw := r.Context().Value(
		middleware.UserContextKey,
	)

	claims, ok := claimsRaw.(*utils.Claims)

	if !ok || claims == nil {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	userID := claims.UserID

	// =====================================================
	// PARSE BODY
	// =====================================================

	var req UpdateSettingsRequest

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
	// VALIDATION
	// =====================================================

	req.Username = strings.TrimSpace(req.Username)

	if req.Username == "" {
		http.Error(
			w,
			"username required",
			http.StatusBadRequest,
		)
		return
	}

	// =====================================================
	// PASSWORD CHANGE
	// =====================================================

	changingPassword := req.Password != ""

	if changingPassword {

		if req.CurrentPassword == "" {
			http.Error(
				w,
				"current password required",
				http.StatusBadRequest,
			)
			return
		}

		// ---------------------------------------------
		// GET CURRENT PASSWORD HASH
		// ---------------------------------------------

		var currentPasswordHash string

		err = db.Pool.QueryRow(
			r.Context(),
			`
			SELECT password_hash
			FROM users
			WHERE id = $1
			`,
			userID,
		).Scan(&currentPasswordHash)

		if err != nil {

			log.Println(
				"GET PASSWORD HASH ERROR:",
				err,
			)

			http.Error(
				w,
				"failed to verify password",
				http.StatusInternalServerError,
			)

			return
		}

		// ---------------------------------------------
		// VERIFY CURRENT PASSWORD
		// ---------------------------------------------

		err = bcrypt.CompareHashAndPassword(
			[]byte(currentPasswordHash),
			[]byte(req.CurrentPassword),
		)

		if err != nil {

			http.Error(
				w,
				"current password is incorrect",
				http.StatusUnauthorized,
			)

			return
		}
	}

	// =====================================================
	// UPDATE USER
	// =====================================================

	if changingPassword {

		// ---------------------------------------------
		// HASH NEW PASSWORD
		// ---------------------------------------------

		hash, err := bcrypt.GenerateFromPassword(
			[]byte(req.Password),
			bcrypt.DefaultCost,
		)

		if err != nil {

			log.Println(
				"HASH ERROR:",
				err,
			)

			http.Error(
				w,
				"failed to hash password",
				http.StatusInternalServerError,
			)

			return
		}

		// ---------------------------------------------
		// UPDATE USERNAME + PASSWORD
		// ---------------------------------------------

		_, err = db.Pool.Exec(
			r.Context(),
			`
			UPDATE users
			SET
				username = $1,
				password_hash = $2
			WHERE id = $3
			`,
			req.Username,
			string(hash),
			userID,
		)

		if err != nil {

			log.Println(
				"UPDATE SETTINGS ERROR:",
				err,
			)

			http.Error(
				w,
				"failed to update settings",
				http.StatusInternalServerError,
			)

			return
		}

	} else {

		// ---------------------------------------------
		// UPDATE USERNAME ONLY
		// ---------------------------------------------

		_, err = db.Pool.Exec(
			r.Context(),
			`
			UPDATE users
			SET
				username = $1
			WHERE id = $2
			`,
			req.Username,
			userID,
		)

		if err != nil {

			log.Println(
				"UPDATE SETTINGS ERROR:",
				err,
			)

			http.Error(
				w,
				"failed to update settings",
				http.StatusInternalServerError,
			)

			return
		}
	}

	// =====================================================
	// CREATE NEW TOKEN
	// =====================================================

	tokenString, err := utils.GenerateJWT(
		userID,
		req.Username,
		claims.IsPremium,
		claims.IsAdmin,
	)

	if err != nil {

		log.Println(
			"TOKEN SIGN ERROR:",
			err,
		)

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

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"token": tokenString,
		},
	)
}