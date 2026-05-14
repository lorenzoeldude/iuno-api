package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"iuno-api/services"
	"iuno-api/utils"
)

func WordHandler(w http.ResponseWriter, r *http.Request) {

	utils.EnableCORS(w)

	slug := strings.TrimPrefix(r.URL.Path, "/api/word/")
	slug = strings.Trim(slug, "/")
	slug = strings.TrimSpace(slug)
	slug = strings.ToLower(slug)

	response, err := services.GetWord(slug)

	if err != nil {
		http.Error(w, "Word not found", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}