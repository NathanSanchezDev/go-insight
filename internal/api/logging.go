package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

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
		err := rows.Scan(
			&logEntry.ID,
			&logEntry.ServiceName,
			&logEntry.LogLevel,
			&logEntry.Message,
			&logEntry.Timestamp,
			&logEntry.TraceID,
			&logEntry.SpanID,
			&logEntry.Metadata,
		)
		if err != nil {
			log.Println("❌ Error scanning log row:", err)
			continue
		}
		logs = append(logs, logEntry)
	}

	return logs, nil
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
