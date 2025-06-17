package observability

import "testing"

func TestGenerateUUIDUnique(t *testing.T) {
	a := GenerateUUID()
	b := GenerateUUID()
	if a == "" || b == "" {
		t.Fatalf("uuid should not be empty")
	}
	if a == b {
		t.Errorf("expected different uuids, got same")
	}
}
