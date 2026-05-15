package policy

import (
	"testing"
	"time"
)

func baseRule() Rule {
	return Rule{
		ID:        "r1",
		KeyPrefix: "app.",
		Severity:  SeverityWarning,
	}
}

func TestValidate_Valid(t *testing.T) {
	if err := baseRule().Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidate_MissingID(t *testing.T) {
	r := baseRule()
	r.ID = ""
	if err := r.Validate(); err == nil {
		t.Fatal("expected error for missing id")
	}
}

func TestValidate_MissingKeyPrefix(t *testing.T) {
	r := baseRule()
	r.KeyPrefix = ""
	if err := r.Validate(); err == nil {
		t.Fatal("expected error for missing key_prefix")
	}
}

func TestValidate_InvalidSeverity(t *testing.T) {
	r := baseRule()
	r.Severity = "unknown"
	if err := r.Validate(); err == nil {
		t.Fatal("expected error for invalid severity")
	}
}

func TestIsExpired_ZeroTime_NotExpired(t *testing.T) {
	r := baseRule()
	if r.IsExpired(time.Now()) {
		t.Fatal("zero expiry should never be expired")
	}
}

func TestIsExpired_PastTime_IsExpired(t *testing.T) {
	r := baseRule()
	r.ExpiresAt = time.Now().Add(-time.Hour)
	if !r.IsExpired(time.Now()) {
		t.Fatal("past expiry should be expired")
	}
}

func TestIsExpired_FutureTime_NotExpired(t *testing.T) {
	r := baseRule()
	r.ExpiresAt = time.Now().Add(time.Hour)
	if r.IsExpired(time.Now()) {
		t.Fatal("future expiry should not be expired")
	}
}
