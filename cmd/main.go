package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NathanSanchezDev/go-insight/internal/api"
	"github.com/NathanSanchezDev/go-insight/internal/db"
	"github.com/gorilla/mux"
)

func main() {
	db.InitDB()

	router := mux.NewRouter()
	router.HandleFunc("/metrics", api.GetMetricsHandler).Methods("GET")
	router.HandleFunc("/logs", api.GetLogsHandler).Methods("GET")

	port := 8080
	fmt.Printf("🚀 Server started on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
