package alert_test

import (
	"strings"
	"testing"

	"github.com/driftwatch/internal/alert"
	"github.com/driftwatch/internal/drift"
)

func makeDiff(key string, src, tgt interface{}, kind drift.DifferenceKind) drift.Difference {
	return drift.Difference{Key: key, SourceValue: src, TargetValue: tgt, Kind: kind}
}

func TestNotify_NoDifferences(t *testing.T) {
	var buf strings.Builder
	n := &alert.Notifier{Writer: &buf}

	result := n.Notify("production", nil)

	if result != nil {
		t.Errorf("expected nil alert when no differences, got %+v", result)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output, got %q", buf.String())
	}
}

func TestNotify_ValueMismatch_WarningLevel(t *testing.T) {
	var buf strings.Builder
	n := &alert.Notifier{Writer: &buf}

	diffs := []drift.Difference{
		makeDiff("db_pool_size", 10, 20, drift.KindValueMismatch),
	}

	a := n.Notify("staging", diffs)

	if a == nil {
		t.Fatal("expected non-nil alert")
	}
	if a.Severity != alert.SeverityWarning {
		t.Errorf("expected WARNING severity, got %s", a.Severity)
	}
	if a.Environment != "staging" {
		t.Errorf("expected environment 'staging', got %s", a.Environment)
	}
	output := buf.String()
	if !strings.Contains(output, "db_pool_size") {
		t.Errorf("expected key in output, got: %s", output)
	}
}

func TestNotify_MissingKey_CriticalLevel(t *testing.T) {
	var buf strings.Builder
	n := &alert.Notifier{Writer: &buf}

	diffs := []drift.Difference{
		makeDiff("api_secret", "set", nil, drift.KindMissing),
	}

	a := n.Notify("production", diffs)

	if a.Severity != alert.SeverityCritical {
		t.Errorf("expected CRITICAL severity, got %s", a.Severity)
	}
	output := buf.String()
	if !strings.Contains(output, "CRITICAL") {
		t.Errorf("expected CRITICAL in output, got: %s", output)
	}
}

func TestNotify_OutputContainsDiffCount(t *testing.T) {
	var buf strings.Builder
	n := &alert.Notifier{Writer: &buf}

	diffs := []drift.Difference{
		makeDiff("key1", "a", "b", drift.KindValueMismatch),
		makeDiff("key2", "x", "y", drift.KindValueMismatch),
	}

	n.Notify("dev", diffs)

	if !strings.Contains(buf.String(), "2 difference(s)") {
		t.Errorf("expected diff count in output, got: %s", buf.String())
	}
}
