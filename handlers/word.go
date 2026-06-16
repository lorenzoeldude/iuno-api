package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	// "context"
	"log"

	"iuno-api/models"
	"iuno-api/services"
	"iuno-api/services/morphology"
	"iuno-api/utils"
	// "iuno-api/db"
)

func WordHandler(w http.ResponseWriter, r *http.Request) {

	utils.EnableCORS(w)

	slug := strings.TrimPrefix(r.URL.Path, "/api/word/")
	slug = strings.Trim(slug, "/")
	slug = strings.TrimSpace(slug)
	slug = strings.ToLower(slug)

	log.Println("word: ", slug)

	response, err := services.GetWord(slug)

	if err != nil {
		http.Error(w, "Word not found", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

func WriteWordHandler(w http.ResponseWriter, r *http.Request) {

	utils.EnableCORS(w)

	// CORS PRE-FLIGHT
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// METHOD CHECK
	// =====================================================
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// PARSE REQUEST
	var body models.WriteRequest

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println("JSON ERROR:", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	log.Println("pronounType handler: ", body.Lemma.PronounType)

	if body.Lemma.Lemma == "" {
		http.Error(w, "lemma is required", http.StatusBadRequest)
		return
	}

	body.Lemma.LemmaNormalized = morphology.NormalizeLatin(body.Lemma.Lemma)

	err = services.WriteWord(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Word created",
	})
}