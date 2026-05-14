package handlers

import (
	"encoding/json"
	"net/http"

	"iuno-api/services"
	"iuno-api/utils"
)

func TrainerHandler(w http.ResponseWriter, r *http.Request) {

	utils.EnableCORS(w)

	words, err := services.GetTrainerWords()

	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	json.NewEncoder(w).Encode(words)
}