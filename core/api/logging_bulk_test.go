package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NathanSanchezDev/go-insight/core/models"
)

// TestPostLogsBulkHandlerSuccess verifies that the handler accepts valid input
// and returns a 201 status code.
func TestPostLogsBulkHandlerSuccess(t *testing.T) {
	// stub bulk insertion to avoid database dependency
	called := false
	postLogsBulkFunc = func(entries []models.Log) error {
		called = true
		return nil
	}
	defer func() { postLogsBulkFunc = PostLogsBulk }()

	body := `[{"service_name":"svc","message":"m1"},{"service_name":"svc","message":"m2"}]`
	req := httptest.NewRequest(http.MethodPost, "/logs/bulk", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	PostLogsBulkHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rr.Code)
	}
	if !called {
		t.Fatalf("postLogsBulkFunc was not called")
	}
}

// TestPostLogsBulkHandlerBadRequest verifies that invalid JSON returns 400.
func TestPostLogsBulkHandlerBadRequest(t *testing.T) {
	postLogsBulkFunc = func(entries []models.Log) error { return nil }
	defer func() { postLogsBulkFunc = PostLogsBulk }()

	req := httptest.NewRequest(http.MethodPost, "/logs/bulk", bytes.NewBufferString("{"))
	rr := httptest.NewRecorder()

	PostLogsBulkHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

// TestPostLogsBulkHandlerValidation verifies that invalid log entries are rejected.
func TestPostLogsBulkHandlerValidation(t *testing.T) {
	postLogsBulkFunc = func(entries []models.Log) error { return nil }
	defer func() { postLogsBulkFunc = PostLogsBulk }()

	body := `[{"message":"m1"}]` // missing service_name
	req := httptest.NewRequest(http.MethodPost, "/logs/bulk", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	PostLogsBulkHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}
