package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/NathanSanchezDev/go-insight/internal/db"
	"github.com/NathanSanchezDev/go-insight/internal/models"
)

func GetMetrics() ([]models.EndpointMetric, error) {
	query := `SELECT id, service_name, path, method, status_code, duration, 
                     language, framework, version, environment, timestamp, request_id 
              FROM metrics`

	rows, err := db.DB.QueryContext(context.Background(), query)
	if err != nil {
		log.Println("❌ Error fetching metrics:", err)
		return nil, err
	}
	defer rows.Close()

	var metrics []models.EndpointMetric

	for rows.Next() {
		var metric models.EndpointMetric
		var source models.MetricSource

		err := rows.Scan(
			&metric.ID, &metric.ServiceName, &metric.Path, &metric.Method, &metric.StatusCode,
			&metric.Duration, &source.Language, &source.Framework, &source.Version,
			&metric.Environment, &metric.Timestamp, &metric.RequestID,
		)
		if err != nil {
			log.Println("❌ Error scanning metric row:", err)
			continue
		}

		metric.Source = source
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func GetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics, err := GetMetrics()
	if err != nil {
		http.Error(w, "Failed to fetch metrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
