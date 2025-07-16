package db

import (
	"log"

	"github.com/NathanSanchezDev/go-insight/core/models"
)

func StoreTrace(trace models.Trace) error {
	query := `INSERT INTO traces (id, service_name, start_time) VALUES ($1, $2, $3)`
	_, err := DB.Exec(query, trace.ID, trace.ServiceName, trace.StartTime)
	if err != nil {
		log.Println("Failed to store trace:", err)
	}
	return err
}

func UpdateTrace(trace *models.Trace) error {
	query := `UPDATE traces SET end_time = $1, duration_ms = $2 WHERE id = $3`
	_, err := DB.Exec(query, trace.EndTime, trace.Duration, trace.ID)
	if err != nil {
		log.Println("Failed to update trace:", err)
	}
	return err
}

func FetchTraces(service string) ([]models.Trace, error) {
	query := `SELECT id, service_name, start_time, end_time, duration_ms FROM traces WHERE service_name = $1 ORDER BY start_time DESC`
	rows, err := DB.Query(query, service)
	if err != nil {
		log.Println("Failed to fetch traces:", err)
		return nil, err
	}
	defer rows.Close()

	var traces []models.Trace
	for rows.Next() {
		var trace models.Trace
		err := rows.Scan(&trace.ID, &trace.ServiceName, &trace.StartTime, &trace.EndTime, &trace.Duration)
		if err != nil {
			log.Println("Error scanning trace row:", err)
			continue
		}
		traces = append(traces, trace)
	}
	return traces, nil
}
