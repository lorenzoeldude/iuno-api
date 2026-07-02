package handlers

import (
	"context"
	"net/http"

	"iuno-api/db"
	"iuno-api/utils"
)

func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")

	if token == "" {
		http.Error(w, "missing token", http.StatusBadRequest)
		return
	}

	hash := utils.HashVerificationToken(token)

	commandTag, err := db.Pool.Exec(context.Background(), `
		UPDATE users
		SET
			email_verified = TRUE,
			email_verification_hash = NULL,
			email_verification_expires_at = NULL
		WHERE
			email_verification_hash = $1
			AND email_verification_expires_at > NOW()
			AND email_verified = FALSE
	`,
		hash,
	)

	if err != nil {
		http.Error(w, "verification failed", http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Redirect(
			w,
			r,
			"http://www.iunoni.com/verify-email?status=invalid",
			http.StatusSeeOther,
		)
		return
	}

	http.Redirect(
		w,
		r,
		"https://www.iunoni.com/login?verified=true",
		http.StatusSeeOther,
	)
}