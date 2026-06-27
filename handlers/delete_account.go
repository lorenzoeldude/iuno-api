package handlers

import (
	"net/http"

	"iuno-api/middleware"
	"iuno-api/services"
	"iuno-api/utils"
)

func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {

	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err := services.DeleteUser(claims.UserID)
	if err != nil {
		http.Error(w, "failed to delete account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}