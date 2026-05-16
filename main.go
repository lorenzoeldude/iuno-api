package main

import (
	"log"
	"net/http"

	"iuno-api/db"
	"iuno-api/handlers"
)

func main() {

	// =====================================================
	// INIT DATABASE
	// =====================================================
	db.Init("postgres://lorenz@localhost:5432/iuno?sslmode=disable")

	// =====================================================
	// DICTIONARY
	// =====================================================
	http.HandleFunc("/api/word/", handlers.WordHandler)

	// =====================================================
	// SEARCH
	// =====================================================
	http.HandleFunc("/api/search", handlers.SearchHandler)

	// =====================================================
	// TRAINER
	// =====================================================
	http.HandleFunc("/api/trainer/random", handlers.RandomTrainerHandler)

	// =====================================================
	// MORPHOLOGY / PARSER
	// =====================================================
	http.HandleFunc("/api/parse", handlers.ParseHandler)

	// =====================================================
	// ADMIN
	// =====================================================
	http.HandleFunc("/api/admin/lemma", handlers.UpsertLemmaHandler)

	// =====================================================
	// AUTH
	// =====================================================
	http.HandleFunc("/api/auth/register", handlers.RegisterHandler)
	http.HandleFunc("/api/auth/login", handlers.LoginHandler)

	// =====================================================
	// WORD LISTS
	// =====================================================
	http.HandleFunc("/api/word-lists", handlers.GetWordListsHandler)

	http.HandleFunc(
		"/api/word-lists/create",
		handlers.CreateWordListHandler,
	)

	http.HandleFunc(
		"/api/word-lists/add-lemma",
		handlers.AddLemmaToListHandler,
	)

	http.HandleFunc(
		"/api/word-lists/lemmas",
		handlers.GetWordListLemmasHandler,
	)

	// =====================================================
	// START SERVER
	// =====================================================
	log.Println("Server running on http://localhost:8080")

	log.Fatal(
		http.ListenAndServe(":8080", nil),
	)
}