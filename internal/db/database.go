package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/NathanSanchezDev/go-insight/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ Database is not reachable:", err)
	}

	log.Println("✅ Connected to PostgreSQL!")
}

func StoreTrace(trace models.Trace) error {
	query := `INSERT INTO traces (id, service_name, start_time) VALUES ($1, $2, $3)`
	_, err := DB.Exec(query, trace.ID, trace.ServiceName, trace.StartTime)
	if err != nil {
		log.Println("Failed to store trace:", err)
	}
	return err
}

func StoreSpan(span models.Span) error {
	query := `INSERT INTO spans (id, trace_id, parent_id, service, operation, start_time) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := DB.Exec(query, span.ID, span.TraceID, span.ParentID, span.Service, span.Operation, span.StartTime)
	if err != nil {
		log.Println("Failed to store span:", err)
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

func UpdateSpan(span *models.Span) error {
	query := `UPDATE spans SET end_time = $1, duration_ms = $2 WHERE id = $3`
	_, err := DB.Exec(query, span.EndTime, span.Duration, span.ID)
	if err != nil {
		log.Println("Failed to update span:", err)
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
