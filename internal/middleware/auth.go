package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedAPIKey := os.Getenv("API_KEY")

		if expectedAPIKey == "" {
			log.Println("⚠️  WARNING: No API_KEY configured, authentication disabled")
			next.ServeHTTP(w, r)
			return
		}

		apiKey := extractAPIKey(r)

		if apiKey == "" {
			log.Printf("🔒 Authentication failed: No API key provided from %s %s", r.Method, r.RequestURI)
			http.Error(w, `{"error": "API key required", "hint": "Provide API key in Authorization header, X-API-Key header, or api_key query parameter"}`, http.StatusUnauthorized)
			return
		}

		if apiKey != expectedAPIKey {
			log.Printf("🔒 Authentication failed: Invalid API key from %s %s", r.Method, r.RequestURI)
			http.Error(w, `{"error": "Invalid API key"}`, http.StatusUnauthorized)
			return
		}

		log.Printf("✅ Authenticated request: %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func extractAPIKey(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		if strings.HasPrefix(authHeader, "Bearer ") {
			return strings.TrimPrefix(authHeader, "Bearer ")
		}
		if strings.HasPrefix(authHeader, "ApiKey ") {
			return strings.TrimPrefix(authHeader, "ApiKey ")
		}
	}

	if apiKey := r.Header.Get("X-API-Key"); apiKey != "" {
		return apiKey
	}

	if apiKey := r.URL.Query().Get("api_key"); apiKey != "" {
		return apiKey
	}

	return ""
}

var PublicEndpoints = map[string]bool{
	"/health": true,
}

func RequiresAuth(path string) bool {
	return !PublicEndpoints[path]
}
