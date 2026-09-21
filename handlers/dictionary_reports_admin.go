package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"iuno-api/db"
)

// =====================================================
// GET DICTIONARY REPORTS
// GET /api/admin/dictionary-reports
// =====================================================

type AdminDictionaryReport struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Username  string    `json:"username"`
	LemmaID   int       `json:"lemmaId"`
	Lemma     string    `json:"lemma"`
	Message   string    `json:"message"`
	Platform  string    `json:"platform"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func GetDictionaryReportsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Pool.Query(
		r.Context(),
		`
		SELECT
			dr.id,
			dr.user_id,
			u.username,
			dr.lemma_id,
			l.lemma,
			dr.message,
			dr.platform,
			dr.status,
			dr.created_at
		FROM dictionary_reports dr
		JOIN lemmas l
			ON l.id = dr.lemma_id
		JOIN users u
			ON u.id = dr.user_id
		ORDER BY dr.created_at DESC
		`,
	)

	if err != nil {
		log.Println("GET DICTIONARY REPORTS ERROR:", err)
		http.Error(
			w,
			"failed to fetch dictionary reports",
			http.StatusInternalServerError,
		)
		return
	}

	defer rows.Close()

	reports := make([]AdminDictionaryReport, 0)

	for rows.Next() {
		var report AdminDictionaryReport

		err := rows.Scan(
			&report.ID,
			&report.UserID,
			&report.Username,
			&report.LemmaID,
			&report.Lemma,
			&report.Message,
			&report.Platform,
			&report.Status,
			&report.CreatedAt,
		)

		if err != nil {
			log.Println(
				"GET DICTIONARY REPORTS SCAN ERROR:",
				err,
			)

			http.Error(
				w,
				"failed to read dictionary reports",
				http.StatusInternalServerError,
			)
			return
		}

		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		log.Println(
			"GET DICTIONARY REPORTS ROW ERROR:",
			err,
		)

		http.Error(
			w,
			"failed to read dictionary reports",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(reports)
}
