package api

import (
	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/health", HealthCheck).Methods("GET")

	// Metrics endpoints
	router.HandleFunc("/metrics", GetMetricsHandler).Methods("GET")
	router.HandleFunc("/metrics", PostMetricHandler).Methods("POST")

	// Logs endpoints
	router.HandleFunc("/logs", GetLogsHandler).Methods("GET")
	router.HandleFunc("/logs", PostLogHandler).Methods("POST")

	return router
}
