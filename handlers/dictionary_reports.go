package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"iuno-api/db"
	"iuno-api/middleware"
	"iuno-api/utils"
)

// =====================================================
// REPORT DICTIONARY ISSUE
// POST /api/dictionary-reports
// =====================================================

type DictionaryReportRequest struct {
	LemmaID  int    `json:"lemmaId"`
	Message  string `json:"message"`
	Platform string `json:"platform"`
}

func ReportDictionaryIssueHandler(w http.ResponseWriter, r *http.Request) {

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
	// AUTH
	// =====================================================

	claimsRaw := r.Context().Value(middleware.UserContextKey)

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

	var req DictionaryReportRequest

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
	// VALIDATE LEMMA ID
	// =====================================================

	if req.LemmaID <= 0 {

		http.Error(
			w,
			"invalid lemma id",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================================
	// VALIDATE MESSAGE
	// =====================================================

	req.Message = strings.TrimSpace(req.Message)

	if req.Message == "" {

		http.Error(
			w,
			"message is required",
			http.StatusBadRequest,
		)

		return
	}

	if len(req.Message) > 1000 {

		http.Error(
			w,
			"message is too long",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================================
	// VALIDATE PLATFORM
	// =====================================================

	if req.Platform != "web" &&
		req.Platform != "ios" {

		http.Error(
			w,
			"invalid platform",
			http.StatusBadRequest,
		)

		return
	}

	// =====================================================
	// VERIFY LEMMA EXISTS
	// =====================================================

	var exists bool

	err = db.Pool.QueryRow(
		r.Context(),
		`
		SELECT EXISTS (
			SELECT 1
			FROM lemmas
			WHERE id = $1
		)
		`,
		req.LemmaID,
	).Scan(&exists)

	if err != nil {

		log.Println(
			"DICTIONARY REPORT LEMMA CHECK ERROR:",
			err,
		)

		http.Error(
			w,
			"failed to verify dictionary entry",
			http.StatusInternalServerError,
		)

		return
	}

	if !exists {

		http.Error(
			w,
			"dictionary entry not found",
			http.StatusNotFound,
		)

		return
	}

	// =====================================================
	// INSERT REPORT
	// =====================================================

	_, err = db.Pool.Exec(
		r.Context(),
		`
		INSERT INTO dictionary_reports (
			user_id,
			lemma_id,
			message,
			platform
		)
		VALUES ($1, $2, $3, $4)
		`,
		userID,
		req.LemmaID,
		req.Message,
		req.Platform,
	)

	if err != nil {

		log.Println(
			"DICTIONARY REPORT INSERT ERROR:",
			err,
		)

		http.Error(
			w,
			"failed to submit dictionary report",
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

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(
		map[string]string{
			"message": "report submitted",
		},
	)
}
