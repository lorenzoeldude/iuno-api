package utils

import "net/http"

func EnableCORS(w http.ResponseWriter, r *http.Request) {

	// =====================================================
	// ALLOWED ORIGINS
	// =====================================================
	origin := r.Header.Get("Origin")

	allowedOrigins := map[string]bool{
		"http://localhost:3000":  true,
		"https://iunoni.com":     true,
		"https://www.iunoni.com": true,
	}

	if allowedOrigins[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	// =====================================================
	// CORS HEADERS
	// =====================================================
	w.Header().Set(
		"Access-Control-Allow-Methods",
		"GET, POST, PUT, DELETE, OPTIONS",
	)

	w.Header().Set(
		"Access-Control-Allow-Headers",
		"Content-Type, Authorization",
	)

	w.Header().Set(
		"Access-Control-Max-Age",
		"86400",
	)
}