package middleware

import (
	"net/http"
	"testing"
)

func TestGetClientIP(t *testing.T) {
	r := &http.Request{Header: http.Header{"X-Forwarded-For": []string{"1.2.3.4"}}}
	if ip := getClientIP(r); ip != "1.2.3.4" {
		t.Errorf("expected X-Forwarded-For to be used, got %s", ip)
	}

	r = &http.Request{Header: http.Header{"X-Real-IP": []string{"2.3.4.5"}}}
	if ip := getClientIP(r); ip != "2.3.4.5" {
		t.Errorf("expected X-Real-IP to be used, got %s", ip)
	}

	r = &http.Request{RemoteAddr: "3.4.5.6:789"}
	if ip := getClientIP(r); ip != "3.4.5.6" {
		t.Errorf("expected remote address host, got %s", ip)
	}

	r = &http.Request{RemoteAddr: "nonsense"}
	if ip := getClientIP(r); ip != "nonsense" {
		t.Errorf("expected entire remote addr when invalid, got %s", ip)
	}
}
