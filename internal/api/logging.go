package api

import (
	"context"
	"log"

	"github.com/NathanSanchezDev/go-insight/internal/db"
	"github.com/NathanSanchezDev/go-insight/internal/models"
)

func InsertLog(service, level, message string) {
	query := `INSERT INTO logs (service_name, log_level, message) VALUES ($1, $2, $3)`

	_, err := db.DB.ExecContext(context.Background(), query, service, level, message)
	if err != nil {
		log.Println("❌ Error inserting log:", err)
		return
	}

	log.Println("✅ Log inserted successfully!")
}

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
