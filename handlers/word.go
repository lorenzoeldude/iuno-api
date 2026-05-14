package handlers

import (
	"encoding/json"
	// "log"
	"net/http"
	"strings"

	"iuno-api/services"
	"iuno-api/utils"
	"iuno-api/models"
)

func WordHandler(w http.ResponseWriter, r *http.Request) {

    utils.EnableCORS(w)

    slug := strings.TrimPrefix(r.URL.Path, "/api/word/")
    slug = strings.ToLower(slug)

    word, err := services.GetWord(slug)
    if err != nil {
        http.Error(w, "Word not found", 404)
        return
    }

    forms, err := services.GetForms(word.ID)
    if err != nil {
        http.Error(w, "Forms not found", 500)
        return
    }

    json.NewEncoder(w).Encode(struct {
        Word  models.Word   `json:"word"`
        Forms []models.Form `json:"forms"`
    }{
        Word:  word,
        Forms: forms,
    })
}