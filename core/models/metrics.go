package models

import "time"

type MetricSource struct {
	Language  string `json:"language"`
	Framework string `json:"framework"`
	Version   string `json:"version"`
}

type EndpointMetric struct {
	ID          int          `json:"id"`
	ServiceName string       `json:"service_name"`
	Path        string       `json:"path"`
	Method      string       `json:"method"`
	StatusCode  int          `json:"status_code"`
	Duration    float64      `json:"duration_ms"`
	Source      MetricSource `json:"source"`
	Environment string       `json:"environment"`
	Timestamp   time.Time    `json:"timestamp"`
	RequestID   string       `json:"request_id"`
}
