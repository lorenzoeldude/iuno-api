package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"iuno-api/db"
	"iuno-api/utils"
)


func GetUserCountHandler(w http.ResponseWriter, r *http.Request) {

	utils.EnableCORS(w)

	// CORS preflight
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}


	if r.Method != http.MethodGet {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}


	var count int

	err := db.Pool.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM users",
	).Scan(&count)


	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}


	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(map[string]int{
		"count": count,
	})
}



func GetLemmaCountHandler(w http.ResponseWriter, r *http.Request) {

	utils.EnableCORS(w)

	// CORS preflight
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}


	if r.Method != http.MethodGet {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}


	var count int

	err := db.Pool.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM lemmas",
	).Scan(&count)


	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}


	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(map[string]int{
		"count": count,
	})
}