package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"iuno-api/services/wordoftheday"
)

func WordOfTheDayHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodGet {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	word, err := wordoftheday.GetToday(
		r.Context(),
	)

	if err != nil {

		log.Println(
			"Word of the day error:",
			err,
		)

		http.Error(
			w,
			"word of the day not found",
			http.StatusNotFound,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(word)
}