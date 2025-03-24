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
)

func main() {
	db.InitDB()
	router := api.SetupRoutes()
	router.Use(loggingMiddleware)
	port := getEnvPort(8080)

	// Start server
	fmt.Printf("🚀 Server started on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

// Middleware for logging requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}

// Get port from environment variable or use default
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
