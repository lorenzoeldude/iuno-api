package handlers

import (
	"encoding/json"
	"net/http"

	"iuno-api/services"
)

func ParseHandler(w http.ResponseWriter, r *http.Request) {

	word := r.URL.Query().Get("word")
	if word == "" {
		http.Error(w, "missing word", 400)
		return
	}

	result := services.ParseWord(word)

	json.NewEncoder(w).Encode(result)
}