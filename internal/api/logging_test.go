package api

import (
	"testing"

	"github.com/NathanSanchezDev/go-insight/internal/models"
)

func TestValidateLogEntry(t *testing.T) {
	valid := models.Log{ServiceName: "svc", Message: "msg", LogLevel: "INFO"}
	if err := validateLogEntry(&valid); err != nil {
		t.Fatalf("valid log returned error: %v", err)
	}

	missingService := models.Log{Message: "msg"}
	if err := validateLogEntry(&missingService); err == nil {
		t.Errorf("expected error for missing service name")
	}

	missingMessage := models.Log{ServiceName: "svc"}
	if err := validateLogEntry(&missingMessage); err == nil {
		t.Errorf("expected error for missing message")
	}

	invalidLevel := models.Log{ServiceName: "svc", Message: "msg", LogLevel: "BAD"}
	if err := validateLogEntry(&invalidLevel); err == nil {
		t.Errorf("expected error for bad log level")
	}
}
