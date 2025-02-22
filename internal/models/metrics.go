package models

import "time"

type MetricSource struct {
	Langauge  string `json:"language"`
	Framework string `json:"framework"`
	Version   string `json:"version"`
}

type EndpointMetric struct {
	Path       string  `json:"path"`
	Method     string  `json:"method"`
	StatusCode int     `json:"status_code"`
	Duration   float64 `json:"duration_ms"`

	ServiceName string       `json:"service_name"`
	Source      MetricSource `json:"source"`
	Environment string       `json:"environment"`

	Timestamp time.Time `json:"timestamp"`
	RequestID string    `json:"request_id"`
}
