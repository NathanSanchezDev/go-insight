package models

import (
	"database/sql"
	"time"
)

type Log struct {
	ID          int            `json:"id"`
	ServiceName string         `json:"service_name"`
	LogLevel    string         `json:"log_level"`
	Message     string         `json:"message"`
	Timestamp   time.Time      `json:"timestamp"`
	TraceID     sql.NullString `json:"trace_id,omitempty"`
	SpanID      sql.NullString `json:"span_id,omitempty"`
	Metadata    sql.NullString `json:"metadata,omitempty"`
}
