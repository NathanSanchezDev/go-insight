package db

import (
	"log"

	"github.com/NathanSanchezDev/go-insight/core/models"
)

func StoreSpan(span models.Span) error {
	query := `INSERT INTO spans (id, trace_id, parent_id, service, operation, start_time) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := DB.Exec(query, span.ID, span.TraceID, span.ParentID, span.Service, span.Operation, span.StartTime)
	if err != nil {
		log.Println("Failed to store span:", err)
	}
	return err
}

func UpdateSpan(span *models.Span) error {
	query := `UPDATE spans SET end_time = $1, duration_ms = $2 WHERE id = $3`
	_, err := DB.Exec(query, span.EndTime, span.Duration, span.ID)
	if err != nil {
		log.Println("Failed to update span:", err)
	}
	return err
}

func FetchSpans(traceID string) ([]models.Span, error) {
	query := `SELECT id, trace_id, parent_id, service, operation, start_time, end_time, duration_ms FROM spans WHERE trace_id = $1 ORDER BY start_time`
	rows, err := DB.Query(query, traceID)
	if err != nil {
		log.Println("Failed to fetch spans:", err)
		return nil, err
	}
	defer rows.Close()

	var spans []models.Span
	for rows.Next() {
		var span models.Span
		err := rows.Scan(&span.ID, &span.TraceID, &span.ParentID, &span.Service, &span.Operation, &span.StartTime, &span.EndTime, &span.Duration)
		if err != nil {
			log.Println("Error scanning span row:", err)
			continue
		}
		spans = append(spans, span)
	}
	return spans, nil
}

func FetchSpanByID(spanID string) (models.Span, error) {
	query := `SELECT id, trace_id, parent_id, service, operation, start_time, end_time, duration_ms FROM spans WHERE id = $1`

	var span models.Span
	err := DB.QueryRow(query, spanID).Scan(
		&span.ID,
		&span.TraceID,
		&span.ParentID,
		&span.Service,
		&span.Operation,
		&span.StartTime,
		&span.EndTime,
		&span.Duration,
	)

	if err != nil {
		log.Println("Failed to fetch span:", err)
		return models.Span{}, err
	}

	return span, nil
}
