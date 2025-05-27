package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/NathanSanchezDev/go-insight/internal/api"
	"github.com/NathanSanchezDev/go-insight/internal/db"
	"github.com/NathanSanchezDev/go-insight/internal/middleware"
)

func main() {
	db.InitDB()
	router := api.SetupRoutes()

	router.Use(middleware.RateLimitMiddleware)
	router.Use(conditionalAuthMiddleware)
	router.Use(loggingMiddleware)

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

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}

func getEnvPort(defaultPort int) int {
	portStr := os.Getenv("PORT")
	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err == nil {
			return port
		}
	}
	return defaultPort
}
