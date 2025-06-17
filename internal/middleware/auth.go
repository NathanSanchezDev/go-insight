package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedAPIKey := os.Getenv("API_KEY")
		jwtSecret := os.Getenv("JWT_SECRET")

		if expectedAPIKey == "" && jwtSecret == "" {
			log.Println("⚠️  WARNING: No API_KEY or JWT_SECRET configured, authentication disabled")
			next.ServeHTTP(w, r)
			return
		}

		apiKey := extractAPIKey(r)
		token := extractJWT(r)

		var role string
		var authenticated bool

		if expectedAPIKey != "" && apiKey == expectedAPIKey {
			role = "admin"
			authenticated = true
		} else if jwtSecret != "" && token != "" {
			claims, err := parseJWT(token, jwtSecret)
			if err != nil {
				log.Printf("🔒 JWT validation error: %v", err)
			} else {
				role = claims.Role
				authenticated = true
			}
		}

		if !authenticated {
			log.Printf("🔒 Authentication failed from %s %s", r.Method, r.RequestURI)
			http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
			return
		}

		requiredRole := EndpointRoles[r.URL.Path]
		if requiredRole != "" && !hasRole(role, requiredRole) {
			log.Printf("🔒 Access denied: role %s required for %s", requiredRole, r.URL.Path)
			http.Error(w, `{"error": "Forbidden"}`, http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "role", role)
		log.Printf("✅ Authenticated request: %s %s as %s", r.Method, r.RequestURI, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractAPIKey(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		if strings.HasPrefix(authHeader, "ApiKey ") {
			return strings.TrimPrefix(authHeader, "ApiKey ")
		}
		if strings.HasPrefix(authHeader, "Bearer ") {
			return strings.TrimPrefix(authHeader, "Bearer ")
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

func extractJWT(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return ""
}

type jwtClaims struct {
	Role string `json:"role"`
	Exp  int64  `json:"exp"`
}

func parseJWT(token, secret string) (*jwtClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token")
	}

	base64Raw := base64.RawURLEncoding
	sig, err := base64Raw.DecodeString(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid signature")
	}

	unsigned := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(unsigned))
	expected := mac.Sum(nil)
	if !hmac.Equal(sig, expected) {
		return nil, fmt.Errorf("signature mismatch")
	}

	payloadBytes, err := base64Raw.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid payload")
	}

	var claims jwtClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, err
	}

	if claims.Exp > 0 && time.Now().Unix() > claims.Exp {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}

var EndpointRoles = map[string]string{
	"/metrics": "user",
	"/logs":    "user",
	"/traces":  "user",
	"/spans":   "user",
}

func hasRole(userRole, required string) bool {
	if required == "" {
		return true
	}
	if userRole == "admin" {
		return true
	}
	return userRole == required
}
