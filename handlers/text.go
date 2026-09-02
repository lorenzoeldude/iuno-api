package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"iuno-api/db"
	"iuno-api/middleware"
	"iuno-api/utils"
)

type Text struct {
	ID          int64         `json:"id"`
	Title       string        `json:"title"`
	Author      string        `json:"author"`
	Description sql.NullString `json:"-"`
	Difficulty  sql.NullString `json:"-"`
	Sections    []TextSection `json:"sections"`
}

type TextSection struct {
	ID         int64  `json:"id"`
	Position   int    `json:"position"`
	Title      string `json:"title"`
	WordListID *int   `json:"word_list_id"`
}

// =========================================================
// TEXT
// =========================================================

func TextHandler(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/api/text/")
	parts := strings.Split(path, "/")

	if len(parts) != 2 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	author := parts[0]
	title := parts[1]

	var text struct {
		ID          int64  `json:"id"`
		Title       string `json:"title"`
		Author      string `json:"author"`
		Description string `json:"description"`
		Difficulty  string `json:"difficulty"`
	}

	err := db.Pool.QueryRow(
		r.Context(),
		`
		SELECT
			id,
			title,
			author,
			COALESCE(description, ''),
			COALESCE(difficulty, '')
		FROM texts
		WHERE
			author = $1
			AND title = $2
		`,
		author,
		title,
	).Scan(
		&text.ID,
		&text.Title,
		&text.Author,
		&text.Description,
		&text.Difficulty,
	)

	if err != nil {
		http.Error(w, "text not found", http.StatusNotFound)
		return
	}

	rows, err := db.Pool.Query(
		r.Context(),
		`
		SELECT
			id,
			position,
			title,
			word_list_id
		FROM text_sections
		WHERE text_id = $1
		ORDER BY position
		`,
		text.ID,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var sections []TextSection

	for rows.Next() {

		var section TextSection

		err := rows.Scan(
			&section.ID,
			&section.Position,
			&section.Title,
			&section.WordListID,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sections = append(sections, section)
	}

	response := struct {
		ID          int64         `json:"id"`
		Title       string        `json:"title"`
		Author      string        `json:"author"`
		Description string        `json:"description"`
		Difficulty  string        `json:"difficulty"`
		Sections    []TextSection `json:"sections"`
	}{
		ID:          text.ID,
		Title:        text.Title,
		Author:      text.Author,
		Description: text.Description,
		Difficulty:  text.Difficulty,
		Sections:    sections,
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

// =========================================================
// TEXT SECTION
// =========================================================

func TextSectionHandler(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/api/text-section/")
	parts := strings.Split(path, "/")

	if len(parts) != 3 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	author := parts[0]
	title := parts[1]

	position, err := strconv.Atoi(parts[2])

	if err != nil {
		http.Error(w, "invalid position", http.StatusBadRequest)
		return
	}

	var response struct {
		TextTitle    string `json:"text_title"`
		Author       string `json:"author"`
		SectionTitle string `json:"section_title"`
		Content      string `json:"content"`
	}

	err = db.Pool.QueryRow(
		r.Context(),
		`
		SELECT
			t.title,
			t.author,
			ts.title,
			ts.content
		FROM texts t
		JOIN text_sections ts
			ON ts.text_id = t.id
		WHERE
			t.author = $1
			AND t.title = $2
			AND ts.position = $3
		`,
		author,
		title,
		position,
	).Scan(
		&response.TextTitle,
		&response.Author,
		&response.SectionTitle,
		&response.Content,
	)

	if err != nil {
		http.Error(w, "section not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

// =========================================================
// ALL TEXTS
// =========================================================

func TextsHandler(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Pool.Query(
		r.Context(),
		`
		SELECT
			id,
			title,
			author,
			COALESCE(description, ''),
			COALESCE(difficulty, '')
		FROM texts
		ORDER BY title ASC
		`,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	type Text struct {
		ID          int64  `json:"id"`
		Title       string `json:"title"`
		Author      string `json:"author"`
		Description string `json:"description"`
		Difficulty  string `json:"difficulty"`
	}

	var texts []Text

	for rows.Next() {

		var text Text

		err := rows.Scan(
			&text.ID,
			&text.Title,
			&text.Author,
			&text.Description,
			&text.Difficulty,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		texts = append(texts, text)
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(texts)
}

// =========================================================
// READING PROGRESS MODELS
// =========================================================

type ReadingProgress struct {
	TextID          int64     `json:"text_id"`
	SectionID       int64     `json:"section_id"`
	CharacterOffset int       `json:"character_offset"`
	Completed       bool      `json:"completed"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type LatestReadingProgress struct {
	TextID          int64     `json:"text_id"`
	Title           string    `json:"title"`
	Author          string    `json:"author"`
	Difficulty      string    `json:"difficulty"`
	SectionID       int64     `json:"section_id"`
	SectionPosition int       `json:"section_position"`
	CharacterOffset int       `json:"character_offset"`
	Completed       bool      `json:"completed"`
	ProgressPercent int       `json:"progress_percent"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// =========================================================
// GET READING PROGRESS FOR ONE TEXT
//
// GET /api/texts/{textID}/progress
// =========================================================

func ReadingProgressHandler(w http.ResponseWriter, r *http.Request) {

	claims, ok := r.Context().Value(
		middleware.UserContextKey,
	).(*utils.Claims)

	if !ok || claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	path := strings.TrimPrefix(
		r.URL.Path,
		"/api/texts/",
	)

	parts := strings.Split(
		strings.Trim(path, "/"),
		"/",
	)

	if len(parts) != 2 || parts[1] != "progress" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	textID, err := strconv.ParseInt(
		parts[0],
		10,
		64,
	)

	if err != nil {
		http.Error(w, "invalid text id", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {

		getReadingProgress(
			w,
			r,
			claims.UserID,
			textID,
		)

		return
	}

	if r.Method == http.MethodPut {

		saveReadingProgress(
			w,
			r,
			claims.UserID,
			textID,
		)

		return
	}

	w.Header().Set(
		"Allow",
		"GET, PUT",
	)

	http.Error(
		w,
		"method not allowed",
		http.StatusMethodNotAllowed,
	)
}

// =========================================================
// GET ONE TEXT'S PROGRESS
// =========================================================

// =========================================================
// GET ONE TEXT'S PROGRESS
// =========================================================

func getReadingProgress(
	w http.ResponseWriter,
	r *http.Request,
	userID int,
	textID int64,
) {

	rows, err := db.Pool.Query(
		r.Context(),
		`
		SELECT
			text_id,
			section_id,
			character_offset,
			completed,
			updated_at
		FROM reading_progress
		WHERE
			user_id = $1
			AND text_id = $2
		ORDER BY updated_at DESC
		`,
		userID,
		textID,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	defer rows.Close()

	progress := make(
		[]ReadingProgress,
		0,
	)

	for rows.Next() {

		var item ReadingProgress

		err := rows.Scan(
			&item.TextID,
			&item.SectionID,
			&item.CharacterOffset,
			&item.Completed,
			&item.UpdatedAt,
		)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		progress = append(
			progress,
			item,
		)
	}

	if err := rows.Err(); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(progress)
}
// =========================================================
// SAVE READING PROGRESS
// =========================================================

func saveReadingProgress(
	w http.ResponseWriter,
	r *http.Request,
	userID int,
	textID int64,
) {

	var request struct {
		SectionID       int64 `json:"section_id"`
		CharacterOffset int   `json:"character_offset"`
		Completed       bool  `json:"completed"`
	}

	err := json.NewDecoder(
		r.Body,
	).Decode(&request)

	if err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if request.SectionID <= 0 {
		http.Error(
			w,
			"invalid section id",
			http.StatusBadRequest,
		)
		return
	}

	if request.CharacterOffset < 0 {
		http.Error(
			w,
			"invalid character offset",
			http.StatusBadRequest,
		)
		return
	}

	// Make sure the section actually belongs
	// to the requested text.
	var exists bool

	err = db.Pool.QueryRow(
		r.Context(),
		`
		SELECT EXISTS (
			SELECT 1
			FROM text_sections
			WHERE
				id = $1
				AND text_id = $2
		)
		`,
		request.SectionID,
		textID,
	).Scan(&exists)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	if !exists {
		http.Error(
			w,
			"section does not belong to text",
			http.StatusBadRequest,
		)
		return
	}

	_, err = db.Pool.Exec(
		r.Context(),
		`
		INSERT INTO reading_progress (
			user_id,
			text_id,
			section_id,
			character_offset,
			completed,
			updated_at
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			NOW()
		)
		ON CONFLICT (user_id, section_id)
		DO UPDATE SET
			text_id = EXCLUDED.text_id,
			character_offset = EXCLUDED.character_offset,
			completed = EXCLUDED.completed,
			updated_at = NOW()
		`,
		userID,
		textID,
		request.SectionID,
		request.CharacterOffset,
		request.Completed,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// =========================================================
// LATEST READING PROGRESS
//
// GET /api/texts/progress/latest
// =========================================================

func LatestReadingProgressHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	claims, ok := r.Context().Value(
		middleware.UserContextKey,
	).(*utils.Claims)

	if !ok || claims == nil {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set(
			"Allow",
			"GET",
		)

		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	var progress LatestReadingProgress

	var sectionContent string

	err := db.Pool.QueryRow(
		r.Context(),
		`
		SELECT
			rp.text_id,
			t.title,
			t.author,
			COALESCE(t.difficulty, ''),
			rp.section_id,
			ts.position,
			rp.character_offset,
			rp.completed,
			rp.updated_at,
			ts.content
		FROM reading_progress rp
		JOIN texts t
			ON t.id = rp.text_id
		JOIN text_sections ts
			ON ts.id = rp.section_id
		WHERE rp.user_id = $1
		ORDER BY rp.updated_at DESC
		LIMIT 1
		`,
		claims.UserID,
	).Scan(
		&progress.TextID,
		&progress.Title,
		&progress.Author,
		&progress.Difficulty,
		&progress.SectionID,
		&progress.SectionPosition,
		&progress.CharacterOffset,
		&progress.Completed,
		&progress.UpdatedAt,
		&sectionContent,
	)

	if err != nil {

		if err == pgx.ErrNoRows {
			http.Error(
				w,
				"reading progress not found",
				http.StatusNotFound,
			)
			return
		}

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	// Calculate progress across the entire book.
	var totalCharacters int64

	err = db.Pool.QueryRow(
		r.Context(),
		`
		SELECT
			COALESCE(
				SUM(char_length(content)),
				0
			)
		FROM text_sections
		WHERE text_id = $1
		`,
		progress.TextID,
	).Scan(&totalCharacters)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	var previousCharacters int64

	err = db.Pool.QueryRow(
		r.Context(),
		`
		SELECT
			COALESCE(
				SUM(char_length(content)),
				0
			)
		FROM text_sections
		WHERE
			text_id = $1
			AND position < $2
		`,
		progress.TextID,
		progress.SectionPosition,
	).Scan(&previousCharacters)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	if progress.Completed {

		progress.ProgressPercent = 100

	} else if totalCharacters > 0 {

		currentOffset := int64(progress.CharacterOffset)

		sectionLength := int64(
			len([]rune(sectionContent)),
		)

		if currentOffset > sectionLength {
			currentOffset = sectionLength
		}

		readCharacters :=
			previousCharacters +
				currentOffset

		percent :=
			int(
				(readCharacters * 100) /
					totalCharacters,
			)

		if percent < 0 {
			percent = 0
		}

		if percent > 100 {
			percent = 100
		}

		progress.ProgressPercent = percent
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(progress)
}