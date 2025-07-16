package api

import (
	"testing"

	"github.com/NathanSanchezDev/go-insight/core/models"
)

func TestValidateMetric(t *testing.T) {
	metric := models.EndpointMetric{
		ServiceName: "svc",
		Path:        "/",
		Method:      "GET",
		StatusCode:  200,
		Duration:    1.0,
		Source:      models.MetricSource{Language: "go"},
	}
	if err := validateMetric(&metric); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	badMethod := metric
	badMethod.Method = "INVALID"
	if err := validateMetric(&badMethod); err == nil {
		t.Errorf("expected error for invalid method")
	}

	badStatus := metric
	badStatus.StatusCode = 99
	if err := validateMetric(&badStatus); err == nil {
		t.Errorf("expected error for status code")
	}

	negativeDuration := metric
	negativeDuration.Duration = -1
	if err := validateMetric(&negativeDuration); err == nil {
		t.Errorf("expected error for negative duration")
	}

	missingSource := metric
	missingSource.Source.Language = ""
	if err := validateMetric(&missingSource); err == nil {
		t.Errorf("expected error for missing source language")
	}
}
