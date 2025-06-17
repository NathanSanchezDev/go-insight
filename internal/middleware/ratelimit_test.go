package middleware

import (
	"net/http"
	"testing"
)

func TestGetClientIPSingleXFF(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("X-Forwarded-For", "1.2.3.4")
	ip := getClientIP(req)
	if ip != "1.2.3.4" {
		t.Errorf("expected 1.2.3.4, got %s", ip)
	}
}

func TestGetClientIPMultipleXFF(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("X-Forwarded-For", "1.2.3.4, 5.6.7.8")
	ip := getClientIP(req)
	if ip != "1.2.3.4" {
		t.Errorf("expected 1.2.3.4, got %s", ip)
	}
}
