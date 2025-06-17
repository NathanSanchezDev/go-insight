package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/NathanSanchezDev/go-insight/internal/db"
	"github.com/NathanSanchezDev/go-insight/internal/models"
)

func GetLogs(serviceName, logLevel, messageContains string, startTime, endTime time.Time, limit, offset int) ([]models.Log, error) {
	query := `SELECT id, service_name, log_level, message, timestamp, trace_id, span_id, metadata 
              FROM logs WHERE 1=1`

	var params []interface{}
	paramCount := 1

	if serviceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", paramCount)
		params = append(params, serviceName)
		paramCount++
	}

	if logLevel != "" {
		query += fmt.Sprintf(" AND log_level = $%d", paramCount)
		params = append(params, logLevel)
		paramCount++
	}

	if messageContains != "" {
		query += fmt.Sprintf(" AND message ILIKE $%d", paramCount)
		params = append(params, "%"+messageContains+"%")
		paramCount++
	}

	if !startTime.IsZero() {
		query += fmt.Sprintf(" AND timestamp >= $%d", paramCount)
		params = append(params, startTime)
		paramCount++
	}

	if !endTime.IsZero() {
		query += fmt.Sprintf(" AND timestamp <= $%d", paramCount)
		params = append(params, endTime)
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

func validateLogEntry(logEntry *models.Log) error {
	if logEntry.ServiceName == "" {
		return errors.New("service name is required")
	}

	if logEntry.Message == "" {
		return errors.New("message is required")
	}

	if logEntry.LogLevel != "" {
		validLevels := map[string]bool{
			"DEBUG": true,
			"INFO":  true,
			"WARN":  true,
			"ERROR": true,
			"FATAL": true,
		}

		if !validLevels[logEntry.LogLevel] {
			return fmt.Errorf("invalid log level: %s", logEntry.LogLevel)
		}
	}

	return nil
}

func sanitizeLogEntry(logEntry *models.Log) {
	logEntry.Message = html.EscapeString(logEntry.Message)
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

	var metadataBytes []byte
	if logEntry.Metadata != nil {
		metadataBytes = *logEntry.Metadata
	} else {
		metadataBytes = []byte("{}")
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
		metadataBytes,
	).Scan(&logEntry.ID)

	if err != nil {
		log.Printf("❌ Error inserting log: %v", err)
		return err
	}

	return nil
}

func GetLogsHandler(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("service")
	logLevel := r.URL.Query().Get("level")
	messageContains := r.URL.Query().Get("message")

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

	logs, err := GetLogs(serviceName, logLevel, messageContains, startTime, endTime, limit, offset)
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
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&logEntry)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validateLogEntry(&logEntry); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sanitizeLogEntry(&logEntry)

	err = PostLog(&logEntry)
	if err != nil {
		http.Error(w, "Failed to save log", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(logEntry)
}

// PostLogsBulk inserts multiple log entries in a single transaction.
func PostLogsBulk(logs []models.Log) error {
	tx, err := db.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(context.Background(), `INSERT INTO logs
                (service_name, log_level, message, timestamp, trace_id, span_id, metadata)
                VALUES ($1, $2, $3, $4, $5, $6, $7)
                RETURNING id`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for i := range logs {
		entry := &logs[i]

		if entry.Timestamp.IsZero() {
			entry.Timestamp = time.Now()
		}

		if entry.TraceID.String == "" {
			entry.TraceID.Valid = false
		} else {
			entry.TraceID.Valid = true
		}

		if entry.SpanID.String == "" {
			entry.SpanID.Valid = false
		} else {
			entry.SpanID.Valid = true
		}

		if entry.Metadata == nil || len(*entry.Metadata) == 0 {
			empty := json.RawMessage("{}")
			entry.Metadata = &empty
		}

		metadataBytes := *entry.Metadata

		if err := stmt.QueryRowContext(context.Background(),
			entry.ServiceName,
			entry.LogLevel,
			entry.Message,
			entry.Timestamp,
			entry.TraceID,
			entry.SpanID,
			metadataBytes,
		).Scan(&entry.ID); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// postLogsBulkFunc allows tests to mock bulk insertion.
var postLogsBulkFunc = PostLogsBulk

// PostLogsBulkHandler handles POST /logs/bulk for inserting multiple logs.
func PostLogsBulkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var entries []models.Log
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&entries); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for i := range entries {
		if err := validateLogEntry(&entries[i]); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sanitizeLogEntry(&entries[i])
	}

	if err := postLogsBulkFunc(entries); err != nil {
		http.Error(w, "Failed to save logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entries)
}
