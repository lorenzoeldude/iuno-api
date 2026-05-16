package main

import (
	"log"
	"net/http"

	"iuno-api/db"
	"iuno-api/handlers"
)

func main() {

	// --------------------
	// INIT DATABASE
	// --------------------
	db.Init("postgres://lorenz@localhost:5432/iuno?sslmode=disable")

	// --------------------
	// ROUTES
	// --------------------

	// dictionary
	http.HandleFunc("/api/word/", handlers.WordHandler)

	// search
	http.HandleFunc("/api/search", handlers.SearchHandler)

	// trainer
	http.HandleFunc("/api/trainer/random", handlers.RandomTrainerHandler)

	// morphology / popup parser
	http.HandleFunc("/api/parse", handlers.ParseHandler)

	//admin adds lemma and meaning (upsert = update or insert)
	http.HandleFunc("/api/admin/lemma", handlers.UpsertLemmaHandler)

	// --------------------
	// START SERVER
	// --------------------
	log.Println("Server running on http://localhost:8080")

	log.Fatal(
		http.ListenAndServe(":8080", nil),
	)
}