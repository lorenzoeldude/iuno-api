package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

// --------------------
// STRICT API SCHEMA (DTO)
// --------------------

type WordResponse struct {
	Latin       string       `json:"latin"`
	Translation string       `json:"translation"`
	Definition  string       `json:"definition"`
	Examples    []string     `json:"examples"`
	Declension  []Declension `json:"declension"`
}

type Declension struct {
	Case     string `json:"case"`
	Singular string `json:"singular"`
	Plural   string `json:"plural"`
}

type SearchResult struct {
	Slug        string `json:"slug"`
	Latin       string `json:"latin"`
	Translation string `json:"translation"`
}

// --------------------
// DB
// --------------------

var db *pgx.Conn

func initDB() {
	connStr := "postgres://localhost:5432/iuno?sslmode=disable"

	var err error
	db, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	log.Println("Connected to PostgreSQL")
}

// --------------------
// CORS
// --------------------

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
}

// --------------------
// SAFE BUILDER
// --------------------

func newWordResponse() WordResponse {
	return WordResponse{
		Examples:   []string{},
		Declension: []Declension{},
	}
}

// --------------------
// WORD ENDPOINT
// --------------------

func wordHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	slug := strings.TrimPrefix(r.URL.Path, "/api/word/")
	slug = strings.ToLower(slug)

	word := newWordResponse()
	var wordID int

	// --------------------
	// WORD QUERY
	// --------------------
	err := db.QueryRow(context.Background(),
		`SELECT id, latin, translation, definition, COALESCE(examples, ARRAY[]::text[])
		 FROM words
		 WHERE slug=$1`,
		slug,
	).Scan(
		&wordID,
		&word.Latin,
		&word.Translation,
		&word.Definition,
		&word.Examples,
	)

	if err != nil {
		http.Error(w, "Word not found", http.StatusNotFound)
		return
	}

	// --------------------
	// DECLENSIONS QUERY
	// --------------------
	rows, err := db.Query(context.Background(),
		`SELECT case_name, singular, plural
		 FROM declensions
		 WHERE word_id=$1`,
		wordID,
	)

	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var d Declension

		if err := rows.Scan(&d.Case, &d.Singular, &d.Plural); err != nil {
			log.Println("scan error:", err)
			continue
		}

		word.Declension = append(word.Declension, d)
	}

	// --------------------
	// RESPONSE
	// --------------------
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(word)
}

// --------------------
// SEARCH ENDPOINT
// --------------------

func searchHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	query := strings.ToLower(r.URL.Query().Get("q"))
	if query == "" {
		json.NewEncoder(w).Encode([]SearchResult{})
		return
	}

	rows, err := db.Query(context.Background(),
		`SELECT slug, latin, translation
		 FROM words
		 WHERE slug LIKE $1 OR LOWER(latin) LIKE $1
		 ORDER BY slug
		 LIMIT 10`,
		query+"%",
	)

	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	defer rows.Close()

	results := []SearchResult{}

	for rows.Next() {
		var r SearchResult

		if err := rows.Scan(&r.Slug, &r.Latin, &r.Translation); err != nil {
			log.Println("scan error:", err)
			continue
		}

		results = append(results, r)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// --------------------
// MAIN
// --------------------

func main() {
	initDB()

	http.HandleFunc("/api/word/", wordHandler)
	http.HandleFunc("/api/search", searchHandler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}