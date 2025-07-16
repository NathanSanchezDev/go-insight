package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/NathanSanchezDev/go-insight/core/db"
	"github.com/NathanSanchezDev/go-insight/core/models"
	"github.com/NathanSanchezDev/go-insight/core/observability"
	"github.com/gorilla/mux"
)

func GetTracesHandler(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("service")

	var startTime, endTime time.Time
	if startTimeStr := r.URL.Query().Get("start_time"); startTimeStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err == nil {
			startTime = parsedTime
		}
	}

	if endTimeStr := r.URL.Query().Get("end_time"); endTimeStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err == nil {
			endTime = parsedTime
		}
	}

	limit := 100
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	traces, err := GetTraces(serviceName, startTime, endTime, limit, offset)
	if err != nil {
		http.Error(w, "Error fetching traces", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(traces)
}

func GetTraces(serviceName string, startTime, endTime time.Time, limit, offset int) ([]models.Trace, error) {
	query := `SELECT id, service_name, start_time, end_time, duration_ms 
              FROM traces WHERE 1=1`

	var params []interface{}
	paramCount := 1

	if serviceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", paramCount)
		params = append(params, serviceName)
		paramCount++
	}

	if !startTime.IsZero() {
		query += fmt.Sprintf(" AND start_time >= $%d", paramCount)
		params = append(params, startTime)
		paramCount++
	}

	if !endTime.IsZero() {
		query += fmt.Sprintf(" AND start_time <= $%d", paramCount)
		params = append(params, endTime)
		paramCount++
	}

	query += " ORDER BY start_time DESC"

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
		log.Println("❌ Error fetching traces:", err)
		return nil, err
	}
	defer rows.Close()

	var traces []models.Trace

	for rows.Next() {
		var trace models.Trace
		err := rows.Scan(
			&trace.ID,
			&trace.ServiceName,
			&trace.StartTime,
			&trace.EndTime,
			&trace.Duration,
		)
		if err != nil {
			log.Println("❌ Error scanning trace row:", err)
			continue
		}

		traces = append(traces, trace)
	}

	return traces, nil
}

func GetTraceByID(traceID string) (*models.Trace, error) {
	query := `SELECT id, service_name, start_time, end_time, duration_ms 
	          FROM traces WHERE id = $1`

	var trace models.Trace
	err := db.DB.QueryRowContext(context.Background(), query, traceID).Scan(
		&trace.ID,
		&trace.ServiceName,
		&trace.StartTime,
		&trace.EndTime,
		&trace.Duration,
	)

	if err != nil {
		log.Printf("❌ Error fetching trace by ID %s: %v", traceID, err)
		return nil, err
	}

	return &trace, nil
}

func GetSpansHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	traceID := vars["traceId"]

	if traceID == "" {
		http.Error(w, "Trace ID is required", http.StatusBadRequest)
		return
	}

	spans, err := db.FetchSpans(traceID)
	if err != nil {
		log.Printf("❌ Error fetching spans: %v", err)
		http.Error(w, "Error fetching spans", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spans)
}

func CreateTraceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var trace models.Trace
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&trace)
	if err != nil {
		log.Printf("❌ Error decoding trace JSON: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validateTrace(&trace); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if trace.ID == "" {
		trace.ID = observability.GenerateUUID()
	}

	if trace.StartTime.IsZero() {
		trace.StartTime = time.Now()
	}

	err = db.StoreTrace(trace)
	if err != nil {
		log.Printf("❌ Error storing trace: %v", err)
		http.Error(w, "Failed to store trace", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(trace)
}

func CreateSpanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var span models.Span
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&span)
	if err != nil {
		log.Printf("❌ Error decoding span JSON: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validateSpan(&span); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if span.ID == "" {
		span.ID = observability.GenerateUUID()
	}

	if span.StartTime.IsZero() {
		span.StartTime = time.Now()
	}

	err = db.StoreSpan(span)
	if err != nil {
		log.Printf("❌ Error storing span: %v", err)
		http.Error(w, "Failed to store span", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(span)
}

func EndTraceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	traceID := vars["traceId"]

	if traceID == "" {
		http.Error(w, "Trace ID is required", http.StatusBadRequest)
		return
	}

	trace, err := GetTraceByID(traceID)
	if err != nil {
		http.Error(w, "Trace not found", http.StatusNotFound)
		return
	}

	now := time.Now()
	trace.EndTime = sql.NullTime{Time: now, Valid: true}
	trace.Duration = sql.NullFloat64{
		Float64: now.Sub(trace.StartTime).Seconds() * 1000,
		Valid:   true,
	}

	err = db.UpdateTrace(trace)
	if err != nil {
		log.Printf("❌ Error updating trace: %v", err)
		http.Error(w, "Failed to update trace", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trace)
}

func EndSpanHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	spanID := vars["spanId"]

	if spanID == "" {
		http.Error(w, "Span ID is required", http.StatusBadRequest)
		return
	}

	span, err := db.FetchSpanByID(spanID)
	if err != nil {
		http.Error(w, "Span not found", http.StatusNotFound)
		return
	}

	span.EndTime = time.Now()
	span.Duration = span.EndTime.Sub(span.StartTime).Seconds() * 1000

	err = db.UpdateSpan(&span)
	if err != nil {
		log.Printf("❌ Error updating span: %v", err)
		http.Error(w, "Failed to update span", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(span)
}

func validateTrace(trace *models.Trace) error {
	if trace.ServiceName == "" {
		return errors.New("service name is required")
	}

	if trace.EndTime.Valid && !trace.StartTime.IsZero() {
		if trace.EndTime.Time.Before(trace.StartTime) {
			return errors.New("end time cannot be before start time")
		}
	}

	return nil
}

func validateSpan(span *models.Span) error {
	if span.Service == "" {
		return errors.New("service name is required")
	}

	if span.TraceID == "" {
		return errors.New("trace ID is required")
	}

	if span.Operation == "" {
		return errors.New("operation is required")
	}

	if !span.EndTime.IsZero() && !span.StartTime.IsZero() {
		if span.EndTime.Before(span.StartTime) {
			return errors.New("end time cannot be before start time")
		}
	}

	return nil
}
