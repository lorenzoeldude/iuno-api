package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"iuno-api/db"
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
	ID       int64  `json:"id"`
	Position int    `json:"position"`
	Title    string `json:"title"`
}

func TextSectionHandler(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/api/text/")
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