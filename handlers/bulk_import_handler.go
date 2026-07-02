package handlers

import (
	"encoding/json"
	"net/http"

	"iuno-api/models"
	"iuno-api/services"
)

type BulkImportResponse struct {
	Imported int      `json:"imported"`
	Failed   int      `json:"failed"`
	Errors   []string `json:"errors"`
}

func BulkImportHandler(w http.ResponseWriter, r *http.Request) {

	// log.Println("bulk import hit")

	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}

	var requests []models.WriteRequest

	err := json.NewDecoder(r.Body).Decode(&requests)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	response := BulkImportResponse{
		Errors: []string{},
	}


	for _, req := range requests {

		err := services.WriteWord(req)

		if err != nil {
			response.Failed++

			response.Errors = append(
				response.Errors,
				err.Error(),
			)

			continue
		}

		response.Imported++
	}


	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(response)
}