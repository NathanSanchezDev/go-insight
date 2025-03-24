package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

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

func PostMetric(metric *models.EndpointMetric) error {
	if metric.Timestamp.IsZero() {
		metric.Timestamp = time.Now()
	}

	query := `INSERT INTO metrics 
		(service_name, path, method, status_code, duration, language, framework, version, environment, timestamp, request_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
		RETURNING id`

	err := db.DB.QueryRowContext(
		context.Background(),
		query,
		metric.ServiceName,
		metric.Path,
		metric.Method,
		metric.StatusCode,
		metric.Duration,
		metric.Source.Language,
		metric.Source.Framework,
		metric.Source.Version,
		metric.Environment,
		metric.Timestamp,
		metric.RequestID,
	).Scan(&metric.ID)

	if err != nil {
		log.Printf("❌ Error inserting metric: %v", err)
		return err
	}

	return nil
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

func PostMetricHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var metric models.EndpointMetric
	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		log.Printf("❌ Error decoding metric JSON: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = PostMetric(&metric)
	if err != nil {
		http.Error(w, "Failed to save metric", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(metric)
}
