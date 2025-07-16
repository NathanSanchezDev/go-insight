package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/NathanSanchezDev/go-insight/core/db"
	"github.com/NathanSanchezDev/go-insight/core/models"
)

func GetMetrics(serviceName, path, method string, minStatus, maxStatus int, limit, offset int) ([]models.EndpointMetric, error) {
	query := `SELECT id, service_name, path, method, status_code, duration, 
                     language, framework, version, environment, timestamp, request_id 
              FROM metrics WHERE 1=1`

	var params []interface{}
	paramCount := 1

	if serviceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", paramCount)
		params = append(params, serviceName)
		paramCount++
	}

	if path != "" {
		query += fmt.Sprintf(" AND path LIKE $%d", paramCount)
		params = append(params, "%"+path+"%")
		paramCount++
	}

	if method != "" {
		query += fmt.Sprintf(" AND method = $%d", paramCount)
		params = append(params, method)
		paramCount++
	}

	if minStatus > 0 {
		query += fmt.Sprintf(" AND status_code >= $%d", paramCount)
		params = append(params, minStatus)
		paramCount++
	}

	if maxStatus > 0 {
		query += fmt.Sprintf(" AND status_code <= $%d", paramCount)
		params = append(params, maxStatus)
		paramCount++
	}

	query += " ORDER BY timestamp DESC"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", paramCount)
		params = append(params, limit)
		paramCount++

		if offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", paramCount)
			params = append(params, offset)
		}
	}

	rows, err := db.DB.QueryContext(context.Background(), query, params...)
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

func validateMetric(metric *models.EndpointMetric) error {
	if metric.ServiceName == "" {
		return errors.New("service name is required")
	}

	if metric.Path == "" {
		return errors.New("path is required")
	}

	if metric.Method == "" {
		return errors.New("method is required")
	}

	validMethods := map[string]bool{
		"GET":     true,
		"POST":    true,
		"PUT":     true,
		"DELETE":  true,
		"PATCH":   true,
		"OPTIONS": true,
		"HEAD":    true,
	}

	if !validMethods[metric.Method] {
		return errors.New("invalid HTTP method")
	}

	if metric.StatusCode < 100 || metric.StatusCode > 599 {
		return errors.New("status code must be between 100 and 599")
	}

	if metric.Duration < 0 {
		return errors.New("duration cannot be negative")
	}

	if metric.Source.Language == "" {
		return errors.New("source language is required")
	}

	return nil
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
	serviceName := r.URL.Query().Get("service")
	path := r.URL.Query().Get("path")
	method := r.URL.Query().Get("method")

	var minStatus, maxStatus int
	if minStatusStr := r.URL.Query().Get("min_status"); minStatusStr != "" {
		if parsed, err := strconv.Atoi(minStatusStr); err == nil {
			minStatus = parsed
		}
	}

	if maxStatusStr := r.URL.Query().Get("max_status"); maxStatusStr != "" {
		if parsed, err := strconv.Atoi(maxStatusStr); err == nil {
			maxStatus = parsed
		}
	}

	limit := 100 // Default limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0 // Default offset
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	metrics, err := GetMetrics(serviceName, path, method, minStatus, maxStatus, limit, offset)
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
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&metric)
	if err != nil {
		log.Printf("❌ Error decoding metric JSON: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validateMetric(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
