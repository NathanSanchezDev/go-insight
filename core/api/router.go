package api

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// API routes FIRST (with /api prefix)
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Health check
	apiRouter.HandleFunc("/health", HealthCheck).Methods("GET")

	// Metrics endpoints
	apiRouter.HandleFunc("/metrics", GetMetricsHandler).Methods("GET")
	apiRouter.HandleFunc("/metrics", PostMetricHandler).Methods("POST")

	// Logs endpoints
	apiRouter.HandleFunc("/logs", GetLogsHandler).Methods("GET")
	apiRouter.HandleFunc("/logs", PostLogHandler).Methods("POST")
	apiRouter.HandleFunc("/logs/bulk", PostLogsBulkHandler).Methods("POST")

	// Traces endpoints
	apiRouter.HandleFunc("/traces", GetTracesHandler).Methods("GET")
	apiRouter.HandleFunc("/traces", CreateTraceHandler).Methods("POST")
	apiRouter.HandleFunc("/traces/{traceId}/end", EndTraceHandler).Methods("POST")
	apiRouter.HandleFunc("/traces/{traceId}/spans", GetSpansHandler).Methods("GET")
	apiRouter.HandleFunc("/spans", CreateSpanHandler).Methods("POST")
	apiRouter.HandleFunc("/spans/{spanId}/end", EndSpanHandler).Methods("POST")

	// Serve Next.js static files LAST (catches everything else)
	webDir := "./web"
	if _, err := os.Stat(webDir); err == nil {
		log.Printf("✅ Serving static files from %s", webDir)
		router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(webDir))))
	} else {
		log.Printf("❌ Web directory not found: %s", webDir)
	}

	return router
}
