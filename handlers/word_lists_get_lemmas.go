package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"iuno-api/db"
	"iuno-api/middleware"
	"iuno-api/models"
	"iuno-api/utils"
)

func GetWordListLemmasHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("HANDLER HIT")

	utils.EnableCORS(w)

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// =====================================================
	// AUTH USER
	// =====================================================

	claimsRaw := r.Context().Value(middleware.UserContextKey)
	if claimsRaw == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	claims, ok := claimsRaw.(*utils.Claims)
	if !ok || claims == nil {
		http.Error(w, "invalid auth context", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID

	// =====================================================
	// GET USER LIST
	// =====================================================

	var listID int

	err := db.Pool.QueryRow(r.Context(), `
		SELECT id
		FROM word_lists
		WHERE user_id = $1
		LIMIT 1
	`, userID).Scan(&listID)

	if err != nil {
		log.Println("LIST NOT FOUND:", err)
		http.Error(w, "word list not found", http.StatusBadRequest)
		return
	}

	// =====================================================
	// QUERY LEMMAS
	// =====================================================

	rows, err := db.Pool.Query(r.Context(), `
		SELECT
			l.id,
			l.lemma_normalized,
			l.lemma,
			l.part_of_speech,
			l.gender,
			l.declension,
			l.conjugation,
			l.perfect,
			l.supine,
			l.irregular,
			l.genitive
		FROM word_list_lemmas wll
		JOIN lemmas l ON l.id = wll.lemma_id
		WHERE wll.list_id = $1
		ORDER BY l.lemma
	`, listID)

	if err != nil {
		log.Println("GET LIST LEMMAS ERROR:", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	lemmas := []models.Lemma{}

	for rows.Next() {

		var lemma models.Lemma

		err := rows.Scan(
			&lemma.ID,
			&lemma.LemmaNormalized,
			&lemma.Lemma,
			&lemma.PartOfSpeech,
			&lemma.Gender,
			&lemma.Declension,
			&lemma.Conjugation,
			&lemma.Perfect,
			&lemma.Supine,
			&lemma.Irregular,
			&lemma.Genitive,
		)

		if err != nil {
			continue
		}

		lemmas = append(lemmas, lemma)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lemmas)
}