package api

import (
	"net/http"
)

// HealthCheck is a simple endpoint to check if the server is running
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
