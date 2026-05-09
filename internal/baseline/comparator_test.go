package baseline_test

import (
	"testing"

	"driftwatch/internal/baseline"
	"driftwatch/internal/state"
)

func makeStore(t *testing.T) *baseline.Store {
	t.Helper()
	store, err := baseline.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return store
}

func makeComparator(t *testing.T) *baseline.Comparator {
	t.Helper()
	return baseline.NewComparator(makeStore(t))
}

func TestCompare_NoDrift_WhenMatchesBaseline(t *testing.T) {
	cmp := makeComparator(t)

	live := state.Snapshot{
		Environment: "production",
		Values:      map[string]string{"replicas": "3", "region": "us-east-1"},
	}
	if err := cmp.CaptureBaseline("staging", "production", live); err != nil {
		t.Fatalf("CaptureBaseline: %v", err)
	}

	result, err := cmp.Compare("staging", "production", live)
	if err != nil {
		t.Fatalf("Compare: %v", err)
	}
	if len(result.Diffs) != 0 {
		t.Errorf("expected no diffs, got %d", len(result.Diffs))
	}
}

func TestCompare_DetectsDrift_WhenValueChanged(t *testing.T) {
	cmp := makeComparator(t)

	original := state.Snapshot{
		Environment: "production",
		Values:      map[string]string{"replicas": "3"},
	}
	_ = cmp.CaptureBaseline("staging", "production", original)

	modified := state.Snapshot{
		Environment: "production",
		Values:      map[string]string{"replicas": "5"},
	}

	result, err := cmp.Compare("staging", "production", modified)
	if err != nil {
		t.Fatalf("Compare: %v", err)
	}
	if len(result.Diffs) == 0 {
		t.Error("expected diffs, got none")
	}
}

func TestCompare_MissingBaseline_ReturnsError(t *testing.T) {
	cmp := makeComparator(t)

	live := state.Snapshot{Environment: "production", Values: map[string]string{}}
	_, err := cmp.Compare("staging", "production", live)
	if err == nil {
		t.Error("expected error when baseline not found")
	}
}

func TestCompare_ReturnsBaselineAge(t *testing.T) {
	cmp := makeComparator(t)

	live := state.Snapshot{
		Environment: "production",
		Values:      map[string]string{"key": "val"},
	}
	_ = cmp.CaptureBaseline("staging", "production", live)

	result, err := cmp.Compare("staging", "production", live)
	if err != nil {
		t.Fatalf("Compare: %v", err)
	}
	if result.BaselineAge == "unknown" || result.BaselineAge == "" {
		t.Errorf("expected baseline age to be set, got %q", result.BaselineAge)
	}
}
