package main

import (
	"log"
	"net/http"
	"os"

	"iuno-api/db"
	"iuno-api/handlers"
	"iuno-api/middleware"
	"iuno-api/email"
)

func main() {

	// =====================================================
	// INIT DATABASE
	// =====================================================
	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		log.Fatal("DATABASE_URL is missing")
	}

	email.Init()

	key := os.Getenv("RESEND_API_KEY")
	log.Printf("RESEND_API_KEY loaded: %v (length=%d)", key != "", len(key))

	db.Init(dbURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("GLOBAL HIT:", r.URL.Path)
	})

	// =====================================================
	// DICTIONARY
	// =====================================================
	http.HandleFunc(
		"/api/word/",
		middleware.CORSMiddleware(
			handlers.WordHandler,
		),
	)

	// TEXTS
	http.HandleFunc(
		"/api/texts",
		middleware.CORSMiddleware(
			handlers.TextsHandler,
		),
	)

	http.HandleFunc(
		"/api/text/",
		middleware.CORSMiddleware(
			handlers.TextHandler,
		),
	)

	http.HandleFunc(
		"/api/text-section/",
		middleware.CORSMiddleware(
			handlers.TextSectionHandler,
		),
	)


	// =====================================================
	// SEARCH
	// =====================================================
	http.HandleFunc(
		"/api/search",
		middleware.CORSMiddleware(
			handlers.SearchFormHandler,
		),
	)

	// =====================================================
	// TRAINER
	// =====================================================
	http.HandleFunc(
		"/api/trainer/random",
		middleware.CORSMiddleware(
			handlers.RandomTrainerHandler,
		),
	)

	http.HandleFunc(
		"/api/trainer/list/random",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(handlers.ListTrainerHandler),
		),
	)

	http.HandleFunc(
		"/api/trainer/book/random",
		middleware.CORSMiddleware(
			handlers.BookTrainerHandler,
		),
	)

	// =====================================================
	// MORPHOLOGY / PARSER
	// =====================================================
	http.HandleFunc(
		"/api/parse",
		middleware.CORSMiddleware(
			handlers.ParseHandler,
		),
	)

	// =====================================================
	// LOOKUP
	// =====================================================
	http.HandleFunc(
		"/api/lookup",
		middleware.CORSMiddleware(
			handlers.ParseFormHandler,
		),
	)

	// =====================================================
	// ADMIN
	// =====================================================

	http.HandleFunc(
		"/admin/users/count",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(
				middleware.AdminOnly(
					handlers.GetUserCountHandler,
				),
			),
		),
	)

	http.HandleFunc(
		"/admin/lemmas/count",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(
				middleware.AdminOnly(
					handlers.GetLemmaCountHandler,
				),
			),
		),
	)

	http.HandleFunc(
		"/api/admin/lemma/",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(
				middleware.AdminOnly(
					handlers.GetLemmaByIDHandler,
				),
			),
		),
	)

	http.HandleFunc(
		"/api/admin/write-word/",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(
				middleware.AdminOnly(
					handlers.WriteWordHandler,
				),
			),
		),
	)

	http.HandleFunc(
		"/api/admin/bulk-import",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(
				middleware.AdminOnly(
					handlers.BulkImportHandler,
				),
			),
		),
	)

	// =====================================================
	// AUTH
	// =====================================================
	http.HandleFunc(
		"/api/auth/register",
		middleware.CORSMiddleware(
			handlers.RegisterHandler,
		),
	)

	http.HandleFunc(
		"/api/auth/login",
		middleware.CORSMiddleware(
			handlers.LoginHandler,
		),
	)

	// =====================================================
	// USER SETTINGS
	// =====================================================
	http.HandleFunc(
		"/api/settings",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(
				handlers.UpdateSettingsHandler,
			),
		),
	)

	// =====================================================
	// WORD LISTS
	// =====================================================

	// get user word lists
	http.HandleFunc(
		"/api/word-lists",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(
				handlers.GetWordListsHandler,
			),
		),
	)

	// create word list
	http.HandleFunc(
		"/api/word-lists/create",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(
				handlers.CreateWordListHandler,
			),
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
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(
				handlers.GetWordListLemmasHandler,
			),
		),
	)

	// =====================================================
	// LEMMA CHECK + DELETE
	// /api/word-lists/lemma/:id
	// =====================================================
	http.HandleFunc(
		"/api/word-lists/lemma/",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(
				func(w http.ResponseWriter, r *http.Request) {

					switch r.Method {

					case http.MethodGet:
						handlers.CheckLemmaSavedHandler(w, r)

					case http.MethodDelete:
						handlers.DeleteLemmaFromUserListHandler(w, r)

					default:
						http.Error(
							w,
							"method not allowed",
							http.StatusMethodNotAllowed,
						)
					}
				},
			),
		),
	)

	// =====================================================
	// ACCOUNT
	// =====================================================
	http.HandleFunc(
		"/api/account",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(
				func(w http.ResponseWriter, r *http.Request) {

					switch r.Method {

					case http.MethodDelete:
						handlers.DeleteAccountHandler(w, r)

					default:
						http.Error(
							w,
							"method not allowed",
							http.StatusMethodNotAllowed,
						)
					}
				},
			),
		),
	)

	// EMAIL VERIFICATION
	http.HandleFunc(
		"/verify-email",
		middleware.CORSMiddleware(
			handlers.VerifyEmailHandler,
		),
	)

	// HEALTH
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// =====================================================
	// START SERVER
	// =====================================================
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)

	log.Fatal(
		http.ListenAndServe(":"+port, nil),
	)
}