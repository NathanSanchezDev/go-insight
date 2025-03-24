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

func GetLogs() ([]models.Log, error) {
	query := `SELECT id, service_name, log_level, message, timestamp, trace_id, span_id, metadata FROM logs`
	rows, err := db.DB.QueryContext(context.Background(), query)
	if err != nil {
		log.Println("❌ Error fetching logs:", err)
		return nil, err
	}
	defer rows.Close()

	var logs []models.Log

	for rows.Next() {
		var logEntry models.Log
		var metadataBytes []byte

		err := rows.Scan(
			&logEntry.ID,
			&logEntry.ServiceName,
			&logEntry.LogLevel,
			&logEntry.Message,
			&logEntry.Timestamp,
			&logEntry.TraceID,
			&logEntry.SpanID,
			&metadataBytes,
		)
		if err != nil {
			log.Println("❌ Error scanning log row:", err)
			continue
		}

		if len(metadataBytes) > 0 {
			rawJSON := json.RawMessage(metadataBytes)
			logEntry.Metadata = &rawJSON
		} else {
			emptyJSON := json.RawMessage("{}")
			logEntry.Metadata = &emptyJSON
		}

		logs = append(logs, logEntry)
	}

	return logs, nil
}

func PostLog(logEntry *models.Log) error {
	if logEntry.Timestamp.IsZero() {
		logEntry.Timestamp = time.Now()
	}

	if logEntry.TraceID.String == "" {
		logEntry.TraceID.Valid = false
	} else {
		logEntry.TraceID.Valid = true
	}

	if logEntry.SpanID.String == "" {
		logEntry.SpanID.Valid = false
	} else {
		logEntry.SpanID.Valid = true
	}

	if logEntry.Metadata == nil {
		emptyJSON := json.RawMessage("{}")
		logEntry.Metadata = &emptyJSON
	} else if len(*logEntry.Metadata) == 0 {
		emptyJSON := json.RawMessage("{}")
		logEntry.Metadata = &emptyJSON
	}

	query := `INSERT INTO logs 
		(service_name, log_level, message, timestamp, trace_id, span_id, metadata) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id`

	err := db.DB.QueryRowContext(
		context.Background(),
		query,
		logEntry.ServiceName,
		logEntry.LogLevel,
		logEntry.Message,
		logEntry.Timestamp,
		logEntry.TraceID,
		logEntry.SpanID,
		logEntry.Metadata,
	).Scan(&logEntry.ID)

	if err != nil {
		log.Printf("❌ Error inserting log: %v", err)
		return err
	}

	return nil
}

func GetLogsHandler(w http.ResponseWriter, r *http.Request) {
	logs, err := GetLogs()
	if err != nil {
		http.Error(w, "Failed to fetch logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func PostLogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var logEntry models.Log
	err := json.NewDecoder(r.Body).Decode(&logEntry)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if logEntry.ServiceName == "" || logEntry.Message == "" {
		http.Error(w, "Service name and message are required", http.StatusBadRequest)
		return
	}

	err = PostLog(&logEntry)
	if err != nil {
		http.Error(w, "Failed to save log", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(logEntry)
}
