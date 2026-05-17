package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"math/rand"

	"iuno-api/db"
	"iuno-api/middleware"
	"iuno-api/utils"
)

type TrainerResponse struct {
	Lemma   string   `json:"lemma"`
	Correct string   `json:"correct"`
	Answers []string `json:"answers"`
}

func ListTrainerHandler(w http.ResponseWriter, r *http.Request) {

	utils.EnableCORS(w)

	log.Println("METHOD:", r.Method)

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claimsRaw := r.Context().Value(middleware.UserContextKey)
	claims, ok := claimsRaw.(*utils.Claims)
	if !ok || claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID

	// =====================================================
	// 1. pick random lemma from user's list
	// =====================================================
	var lemmaID int
	var lemma string

	err := db.Pool.QueryRow(r.Context(), `
		SELECT l.id, l.lemma
		FROM word_list_lemmas wll
		JOIN word_lists wl ON wl.id = wll.list_id
		JOIN lemmas l ON l.id = wll.lemma_id
		WHERE wl.user_id = $1
		ORDER BY RANDOM()
		LIMIT 1
	`, userID).Scan(&lemmaID, &lemma)

	if err != nil {
		log.Println("NO WORDS ERROR:", err)
		http.Error(w, "no words in list", http.StatusBadRequest)
		return
	}

	// =====================================================
	// 2. get correct meaning (1 row only)
	// =====================================================
	var correct string

	err = db.Pool.QueryRow(r.Context(), `
		SELECT m.meaning
		FROM meanings m
		WHERE m.lemma_id = $1
		ORDER BY m.sort_order ASC
		LIMIT 1
	`, lemmaID).Scan(&correct)

	if err != nil {
		log.Println("MEANING FETCH ERROR:", err)
		http.Error(w, "no meaning found", http.StatusBadRequest)
		return
	}

	// =====================================================
	// 3. get WRONG meanings (2 random ones)
	// =====================================================
	rows, err := db.Pool.Query(r.Context(), `
		SELECT m.meaning
		FROM meanings m
		WHERE m.meaning IS NOT NULL
		AND m.meaning != $1
		ORDER BY RANDOM()
		LIMIT 10
	`, correct)

	if err != nil {
		log.Println("WRONG MEANINGS ERROR:", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var answers []string
	answers = append(answers, correct)

	for rows.Next() {
		var m string
		rows.Scan(&m)
		answers = append(answers, m)

		if len(answers) == 3 {
			break
		}
	}

	// fallback safety (rare but prevents crashes)
	if len(answers) < 3 {
		for len(answers) < 3 {
			answers = append(answers, "—")
		}
	}

	// shuffle
	rand.Shuffle(len(answers), func(i, j int) {
		answers[i], answers[j] = answers[j], answers[i]
	})

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(TrainerResponse{
		Lemma:   lemma,
		Correct: correct,
		Answers: answers,
	})
}