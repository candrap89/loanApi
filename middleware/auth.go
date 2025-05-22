package middleware

import (
	"net/http"
)

// APIKeyMiddleware verifies the presence and validity of the API key
func APIKeyMiddleware(validKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get API key from header
			apiKey := r.Header.Get("X-API-Key")

			// Check if API key is valid
			if apiKey != validKey {
				http.Error(w, "Unauthorized: Invalid or missing API key", http.StatusUnauthorized)
				return
			}

			// Proceed to next handler if valid
			next.ServeHTTP(w, r)
		})
	}
}
