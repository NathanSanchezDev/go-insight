package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/NathanSanchezDev/go-insight/config"
	"github.com/NathanSanchezDev/go-insight/core/api"
	"github.com/NathanSanchezDev/go-insight/core/db"
	"github.com/NathanSanchezDev/go-insight/core/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db.InitDB(cfg)
	router := api.SetupRoutes()

	// Apply middleware to main router, but auth will check if path needs it
	router.Use(middleware.RateLimitMiddleware(cfg))
	router.Use(conditionalAuthMiddleware)
	router.Use(loggingMiddleware(cfg))

	port := getEnvPort(8080)

	fmt.Printf("🚀 Server started on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

func conditionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !middleware.RequiresAuth(r.URL.Path) {
			log.Printf("📖 Public endpoint accessed: %s %s", r.Method, r.RequestURI)
			next.ServeHTTP(w, r)
			return
		}

		middleware.AuthMiddleware(next).ServeHTTP(w, r)
	})
}

func loggingMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)

			if cfg.Features.DebugLogging {
				log.Printf("DEBUG: %s %s %s", r.Method, r.RequestURI, time.Since(start))
			} else {
				log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
			}
		})
	}
}

func getEnvPort(defaultPort int) int {
	portStr := os.Getenv("GO_PORT")
	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err == nil {
			return port
		}
	}
	return defaultPort
}
