package models

import (
	"database/sql"
	"time"
)

type Trace struct {
	ID          string          `json:"id"`
	ServiceName string          `json:"service_name"`
	StartTime   time.Time       `json:"start_time"`
	EndTime     sql.NullTime    `json:"end_time"`
	Duration    sql.NullFloat64 `json:"duration_ms"`
}
