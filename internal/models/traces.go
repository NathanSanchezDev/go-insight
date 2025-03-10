package models

import "time"

type Trace struct {
	ID          string    `json:"id"`
	ServiceName string    `json:"service_name"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Duration    float64   `json:"duration_ms"`
}
