package middleware

import "net/http"

func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

		// Allow browsers to cache preflight responses
		w.Header().Set(
			"Access-Control-Max-Age",
			"86400",
		)

		// =====================================================
		// PRE-FLIGHT REQUEST
		// =====================================================
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// =====================================================
		// NEXT HANDLER
		// =====================================================
		next.ServeHTTP(w, r)
	}
}