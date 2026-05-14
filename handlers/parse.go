package handlers

import (
	"encoding/json"
	"net/http"

	"iuno-api/services"
	"iuno-api/utils"
)

func ParseHandler(w http.ResponseWriter, r *http.Request) {

	utils.EnableCORS(w)

	word := r.URL.Query().Get("word")
	if word == "" {
		http.Error(w, "missing word", 400)
		return
	}

	result := services.ParseWord(word)

	json.NewEncoder(w).Encode(result)
}