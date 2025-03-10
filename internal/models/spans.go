package models

import "time"

type Span struct {
	ID        string    `json:"id"`
	TraceID   string    `json:"trace_id"`
	ParentID  string    `json:"parent_id"`
	Service   string    `json:"service"`
	Operation string    `json:"operation"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  float64   `json:"duration_ms"`
}
