package main

import (
	"log"
	"net/http"

	"iuno-api/db"
	"iuno-api/handlers"
	"iuno-api/middleware"
)

func main() {

	// =====================================================
	// INIT DATABASE
	// =====================================================
	db.Init("postgres://lorenz@localhost:5432/iuno?sslmode=disable")

	// =====================================================
	// DICTIONARY
	// =====================================================
	http.HandleFunc(
		"/api/word/",
		handlers.WordHandler,
	)

	// =====================================================
	// SEARCH
	// =====================================================
	http.HandleFunc(
		"/api/search",
		handlers.SearchHandler,
	)

	// =====================================================
	// TRAINER
	// =====================================================
	http.HandleFunc(
		"/api/trainer/random",
		handlers.RandomTrainerHandler,
	)

	// =====================================================
	// MORPHOLOGY / PARSER
	// =====================================================
	http.HandleFunc(
		"/api/parse",
		handlers.ParseHandler,
	)

	// =====================================================
	// ADMIN
	// =====================================================
	http.HandleFunc(
		"/api/admin/lemma",
		handlers.UpsertLemmaHandler,
	)

	// =====================================================
	// AUTH
	// =====================================================
	http.HandleFunc(
		"/api/auth/register",
		handlers.RegisterHandler,
	)

	http.HandleFunc(
		"/api/auth/login",
		handlers.LoginHandler,
	)

	// =====================================================
	// WORD LISTS
	// =====================================================

	// get user word lists
	http.HandleFunc(
		"/api/word-lists",
		middleware.AuthMiddleware(
			handlers.GetWordListsHandler,
		),
	)

	// create word list
	http.HandleFunc(
		"/api/word-lists/create",
		middleware.AuthMiddleware(
			handlers.CreateWordListHandler,
		),
	)

	// add lemma to list
	http.HandleFunc(
		"/api/word-lists/add-lemma",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(
				handlers.AddLemmaToUserListHandler,
			),
		),
	)

	// get lemmas inside a list
	http.HandleFunc(
		"/api/word-lists/lemmas",
		middleware.AuthMiddleware(
			handlers.GetWordListLemmasHandler,
		),
	)

	// =====================================================
	// START SERVER
	// =====================================================
	log.Println("Server running on http://localhost:8080")

	log.Fatal(
		http.ListenAndServe(":8080", nil),
	)
}