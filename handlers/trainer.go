package handlers

import (
	"encoding/json"
	"net/http"

	"iuno-api/services"
)

func TrainerHandler(w http.ResponseWriter, r *http.Request) {

	words, err := services.GetTrainerWords()

	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	json.NewEncoder(w).Encode(words)
}