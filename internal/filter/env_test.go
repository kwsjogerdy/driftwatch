package filter

import (
	"testing"
)

func TestNewEnvFilter_Valid(t *testing.T) {
	f, err := NewEnvFilter([]string{"staging", "production"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewEnvFilter_Empty(t *testing.T) {
	_, err := NewEnvFilter([]string{})
	if err == nil {
		t.Fatal("expected error for empty slice, got nil")
	}
}

func TestNewEnvFilter_BlankName(t *testing.T) {
	_, err := NewEnvFilter([]string{"staging", "  "})
	if err == nil {
		t.Fatal("expected error for blank name, got nil")
	}
}

func TestAllow_CaseInsensitive(t *testing.T) {
	f, _ := NewEnvFilter([]string{"Staging", "PRODUCTION"})

	if !f.Allow("staging") {
		t.Error("expected 'staging' to be allowed")
	}
	if !f.Allow("STAGING") {
		t.Error("expected 'STAGING' to be allowed")
	}
	if !f.Allow("production") {
		t.Error("expected 'production' to be allowed")
	}
	if f.Allow("dev") {
		t.Error("expected 'dev' to be rejected")
	}
}

func TestFilterKeys_ReturnsOnlyAllowed(t *testing.T) {
	f, _ := NewEnvFilter([]string{"staging", "production"})

	input := map[string]interface{}{
		"staging":    "val1",
		"production": "val2",
		"dev":        "val3",
	}

	result := f.FilterKeys(input)

	if len(result) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["dev"]; ok {
		t.Error("expected 'dev' to be excluded from result")
	}
	if result["staging"] != "val1" {
		t.Errorf("expected staging=val1, got %v", result["staging"])
	}
}

func TestFilterKeys_EmptyInput(t *testing.T) {
	f, _ := NewEnvFilter([]string{"staging"})
	result := f.FilterKeys(map[string]interface{}{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d keys", len(result))
	}
}

func TestNames_ReturnsAllowed(t *testing.T) {
	f, _ := NewEnvFilter([]string{"staging", "production"})
	names := f.Names()
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}
