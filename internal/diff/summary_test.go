package diff

import (
	"strings"
	"testing"
)

func makeDiff(key, kind string, sev Severity) Difference {
	return Difference{Key: key, Source: "a", Target: "b", Kind: kind, Severity: sev}
}

func TestNewSummary_NoDifferences(t *testing.T) {
	s := NewSummary("prod", "staging", nil)
	if s.TotalDrift != 0 {
		t.Fatalf("expected 0 drift, got %d", s.TotalDrift)
	}
	if s.Label() != "clean" {
		t.Errorf("expected label 'clean', got %q", s.Label())
	}
}

func TestNewSummary_CriticalFlagSet(t *testing.T) {
	diffs := []Difference{makeDiff("db_password", "mismatch", SeverityCritical)}
	s := NewSummary("prod", "staging", diffs)
	if !s.HasCritical {
		t.Error("expected HasCritical to be true")
	}
	if s.Label() != "critical" {
		t.Errorf("expected label 'critical', got %q", s.Label())
	}
}

func TestNewSummary_WarningOnly(t *testing.T) {
	diffs := []Difference{makeDiff("log_level", "mismatch", SeverityWarning)}
	s := NewSummary("prod", "staging", diffs)
	if s.HasCritical {
		t.Error("expected HasCritical to be false")
	}
	if !s.HasWarning {
		t.Error("expected HasWarning to be true")
	}
	if s.Label() != "warning" {
		t.Errorf("expected label 'warning', got %q", s.Label())
	}
}

func TestSummary_String_ContainsEnvNames(t *testing.T) {
	s := NewSummary("prod", "staging", nil)
	str := s.String()
	if !strings.Contains(str, "prod") || !strings.Contains(str, "staging") {
		t.Errorf("String() missing env names: %q", str)
	}
}
