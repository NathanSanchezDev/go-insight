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

	// Traces endpoints
	router.HandleFunc("/traces", GetTracesHandler).Methods("GET")
	router.HandleFunc("/traces", CreateTraceHandler).Methods("POST")
	router.HandleFunc("/traces/{traceId}/end", EndTraceHandler).Methods("POST")
	router.HandleFunc("/traces/{traceId}/spans", GetSpansHandler).Methods("GET")
	router.HandleFunc("/spans", CreateSpanHandler).Methods("POST")
	router.HandleFunc("/spans/{spanId}/end", EndSpanHandler).Methods("POST")

	return router
}
