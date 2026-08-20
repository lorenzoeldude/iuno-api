package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"iuno-api/db"
	"iuno-api/email"
	"iuno-api/handlers"
	"iuno-api/middleware"
	"iuno-api/stripe"
	stripehandlers "iuno-api/handlers/stripehandlers"
)

func main() {

	if err := godotenv.Load(); err != nil {
        log.Println(".env file not found")
    }

	// =====================================================
	// INIT DATABASE
	// =====================================================

	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		log.Fatal("DATABASE_URL is missing")
	}

	email.Init()

	key := os.Getenv("RESEND_API_KEY")
	log.Printf(
		"RESEND_API_KEY loaded: %v (length=%d)",
		key != "",
		len(key),
	)

	db.Init(dbURL)
	stripe.Init()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("GLOBAL HIT:", r.URL.Path)
	})

	// =====================================================
	// STRIPE
	// =====================================================

	http.HandleFunc(
		"/api/stripe/create-checkout-session",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(
				stripehandlers.CreateCheckoutSessionHandler,
			),
		),
	)

	http.HandleFunc(
		"/api/stripe/webhook",
		middleware.CORSMiddleware(
			stripehandlers.StripeWebhookHandler,
		),
	)

	// =====================================================
	// STRIPE CUSTOMER PORTAL
	// =====================================================

	http.HandleFunc(
		"/api/stripe/create-portal-session",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(
				stripehandlers.CreateStripePortalSessionHandler,
			),
		),
	)

	// =====================================================
	// BILLING STATUS
	// =====================================================

	http.HandleFunc(
		"/api/billing/status",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(
				stripehandlers.GetBillingStatusHandler,
			),
		),
	)

	// =====================================================
	// APPLE STOREKIT
	// =====================================================

	http.HandleFunc(
		"/api/apple/storekit/transaction",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(
				handlers.AppleStoreKitTransactionHandler,
			),
		),
	)

	// =====================================================
	// DICTIONARY
	// =====================================================

	http.HandleFunc(
		"/api/word/",
		middleware.CORSMiddleware(
			handlers.WordHandler,
		),
	)

	// =====================================================
	// WORD OF THE DAY
	// =====================================================

	http.HandleFunc(
		"/api/word-of-the-day",
		middleware.CORSMiddleware(
			handlers.WordOfTheDayHandler,
		),
	)

	// =====================================================
	// TEXTS
	// =====================================================

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
			middleware.AnonymousTrainerMiddleware(
				handlers.RandomTrainerHandler,
			),
		),
	)

	http.HandleFunc(
		"/api/trainer/list/random",
		middleware.CORSMiddleware(
			middleware.AnonymousTrainerMiddleware(
				handlers.ListTrainerHandler,
			),
		),
	)

	http.HandleFunc(
		"/api/trainer/book/random",
		middleware.CORSMiddleware(
			middleware.AnonymousTrainerMiddleware(
				handlers.BookTrainerHandler,
			),
		),
	)

	// =====================================================
	// TRAINING STATS
	// =====================================================

	http.HandleFunc(
		"/api/training/stats",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(
				handlers.GetTrainingStatsHandler,
			),
		),
	)

	http.HandleFunc(
		"/api/training/attempt",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(
				handlers.RecordTrainingAttemptHandler,
			),
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
	// LESSONS (ADMIN)
	// =====================================================

	http.HandleFunc(
		"/api/admin/lessons/",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(
				middleware.AdminOnly(
					func(w http.ResponseWriter, r *http.Request) {

						// =====================================================
						// LESSON VOCABULARY
						// =====================================================

						if strings.HasSuffix(
							r.URL.Path,
							"/vocabulary",
						) {

							switch r.Method {

							case http.MethodGet:
								handlers.GetLessonVocabularyHandler(
									w,
									r,
								)

							case http.MethodPut:
								handlers.UpdateLessonVocabularyHandler(
									w,
									r,
								)

							default:
								http.Error(
									w,
									"method not allowed",
									http.StatusMethodNotAllowed,
								)
							}

							return
						}

						// =====================================================
						// ADMIN LESSON
						// =====================================================

						switch r.Method {

						case http.MethodGet:
							handlers.GetLessonHandler(w, r)

						case http.MethodPost:
							handlers.CreateLessonHandler(w, r)

						case http.MethodPut:
							handlers.UpdateLessonHandler(w, r)

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
		),
	)

	// =====================================================
	// LESSONS - PUBLIC LIST
	// GET /api/lessons
	// =====================================================

	http.HandleFunc(
		"/api/lessons",
		middleware.CORSMiddleware(
			handlers.GetLessonsHandler,
		),
	)

	// =====================================================
	// LESSONS - PUBLIC / USER
	// =====================================================

	http.HandleFunc(
		"/api/lessons/",
		middleware.CORSMiddleware(
			func(w http.ResponseWriter, r *http.Request) {

				// =====================================================
				// LESSON PROGRESS
				// GET /api/lessons/{id}/progress
				// PUT /api/lessons/{id}/progress
				// =====================================================

				if strings.HasSuffix(
					r.URL.Path,
					"/progress",
				) {

					middleware.AuthMiddleware(
						func(w http.ResponseWriter, r *http.Request) {

							switch r.Method {

							case http.MethodGet:
								handlers.GetLessonProgressHandler(
									w,
									r,
								)

							case http.MethodPut:
								handlers.UpdateLessonProgressHandler(
									w,
									r,
								)

							default:
								http.Error(
									w,
									"method not allowed",
									http.StatusMethodNotAllowed,
								)
							}
						},
					)(w, r)

					return
				}

				// =====================================================
				// PREVIOUS LESSONS TRAINER
				// GET /api/lessons/{id}/trainer/previous/random
				// =====================================================

				if strings.HasSuffix(
					r.URL.Path,
					"/trainer/previous/random",
				) {

					middleware.AnonymousTrainerMiddleware(
						func(w http.ResponseWriter, r *http.Request) {

							switch r.Method {

							case http.MethodGet:
								handlers.LessonPreviousTrainerHandler(w, r)

							default:
								http.Error(
									w,
									"method not allowed",
									http.StatusMethodNotAllowed,
								)
							}

						},
					)(w, r)

					return
				}

				// =====================================================
				// LESSON TRAINER
				// GET /api/lessons/{id}/trainer/random
				// =====================================================

				if strings.HasSuffix(
					r.URL.Path,
					"/trainer/random",
				) {

					middleware.AnonymousTrainerMiddleware(
						func(w http.ResponseWriter, r *http.Request) {

							switch r.Method {

							case http.MethodGet:
								handlers.LessonTrainerHandler(w, r)

							default:
								http.Error(
									w,
									"method not allowed",
									http.StatusMethodNotAllowed,
								)
							}

						},
					)(w, r)

					return
				}

				// =====================================================
				// LESSON VOCABULARY
				// GET /api/lessons/{id}/vocabulary
				// =====================================================

				if strings.HasSuffix(
					r.URL.Path,
					"/vocabulary",
				) {

					switch r.Method {

					case http.MethodGet:
						handlers.GetLessonVocabularyTrainerHandler(
							w,
							r,
						)

					default:
						http.Error(
							w,
							"method not allowed",
							http.StatusMethodNotAllowed,
						)
					}

					return
				}

				// =====================================================
				// LESSON
				// GET /api/lessons/{id}
				// =====================================================

				switch r.Method {

				case http.MethodGet:
					handlers.GetLessonHandler(w, r)

				default:
					http.Error(
						w,
						"method not allowed",
						http.StatusMethodNotAllowed,
					)
				}
			},
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

	http.HandleFunc(
		"/api/word-lists",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(
				handlers.GetWordListsHandler,
			),
		),
	)

	http.HandleFunc(
		"/api/word-lists/create",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(
				handlers.CreateWordListHandler,
			),
		),
	)

	http.HandleFunc(
		"/api/word-lists/add-lemma",
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(
				handlers.AddLemmaToUserListHandler,
			),
		),
	)

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
						handlers.DeleteLemmaFromUserListHandler(
							w,
							r,
						)

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

	// =====================================================
	// EMAIL VERIFICATION
	// =====================================================

	http.HandleFunc(
		"/verify-email",
		middleware.CORSMiddleware(
			handlers.VerifyEmailHandler,
		),
	)

	// =====================================================
	// HEALTH
	// =====================================================

	http.HandleFunc(
		"/health",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		},
	)

	// =====================================================
	// START SERVER
	// =====================================================

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Printf(
		"Server running on port %s",
		port,
	)

	log.Fatal(
		http.ListenAndServe(
			":"+port,
			nil,
		),
	)
}
