package api

import (
	"encoding/json"
	"net/http"

	"github.com/NathanSanchezDev/go-insight/internal/db"
)

func GetTraces(w http.ResponseWriter, r *http.Request) {
	service := r.URL.Query().Get("service")

	traces, err := db.FetchTraces(service)
	if err != nil {
		http.Error(w, "Error fetching traces", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(traces)
}
