package middleware

import (
	"net/http"
	"net/url"
	"testing"
)

func TestExtractAPIKey(t *testing.T) {
	cases := []struct {
		header http.Header
		query  url.Values
		want   string
	}{
		{header: http.Header{"Authorization": []string{"Bearer token123"}}, want: "token123"},
		{header: http.Header{"Authorization": []string{"ApiKey abcdef"}}, want: "abcdef"},
		{header: http.Header{"X-Api-Key": []string{"headerkey"}}, want: "headerkey"},
		{query: url.Values{"api_key": []string{"querykey"}}, want: "querykey"},
		{header: http.Header{}, want: ""},
	}

	for _, c := range cases {
		r := &http.Request{Header: c.header}
		r.URL = &url.URL{RawQuery: c.query.Encode()}
		got := extractAPIKey(r)
		if got != c.want {
			t.Errorf("extractAPIKey = %q, want %q", got, c.want)
		}
	}
}

func TestRequiresAuth(t *testing.T) {
	if RequiresAuth("/health") {
		t.Errorf("/health should not require auth")
	}
	if RequiresAuth("/api/health") {
		t.Errorf("/api/health should not require auth")
	}
	if !RequiresAuth("/api/logs") {
		t.Errorf("/api/logs should require auth")
	}
	if !RequiresAuth("/api/logs/bulk") {
		t.Errorf("/api/logs/bulk should require auth")
	}
	if RequiresAuth("/") {
		t.Errorf("/ should not require auth (static files)")
	}
	if RequiresAuth("/dashboard") {
		t.Errorf("/dashboard should not require auth (static files)")
	}
}
