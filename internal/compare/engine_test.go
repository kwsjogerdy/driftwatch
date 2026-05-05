package compare

import (
	"errors"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/state"
)

func makeSnap(values map[string]string) *state.Snapshot {
	return &state.Snapshot{Values: values}
}

func stubLoader(snapshots map[string]*state.Snapshot) func(string, string) (*state.Snapshot, error) {
	return func(_ string, env string) (*state.Snapshot, error) {
		snap, ok := snapshots[env]
		if !ok {
			return nil, errors.New("environment not found: " + env)
		}
		return snap, nil
	}
}

func TestCompare_NoDrift(t *testing.T) {
	snaps := map[string]*state.Snapshot{
		"staging": makeSnap(map[string]string{"key": "value"}),
		"prod":    makeSnap(map[string]string{"key": "value"}),
	}
	engine := NewEngineWithDeps(stubLoader(snaps), drift.Detect)

	result, err := engine.Compare("state.json", "staging", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasDrift {
		t.Errorf("expected no drift, got %d diff(s)", len(result.Diffs))
	}
}

func TestCompare_WithDrift(t *testing.T) {
	snaps := map[string]*state.Snapshot{
		"staging": makeSnap(map[string]string{"key": "staging-value"}),
		"prod":    makeSnap(map[string]string{"key": "prod-value"}),
	}
	engine := NewEngineWithDeps(stubLoader(snaps), drift.Detect)

	result, err := engine.Compare("state.json", "staging", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasDrift {
		t.Error("expected drift but got none")
	}
	if len(result.Diffs) != 1 {
		t.Errorf("expected 1 diff, got %d", len(result.Diffs))
	}
}

func TestCompare_SourceEnvMissing(t *testing.T) {
	snaps := map[string]*state.Snapshot{
		"prod": makeSnap(map[string]string{"key": "value"}),
	}
	engine := NewEngineWithDeps(stubLoader(snaps), drift.Detect)

	_, err := engine.Compare("state.json", "staging", "prod")
	if err == nil {
		t.Error("expected error for missing source environment, got nil")
	}
}

func TestCompare_TargetEnvMissing(t *testing.T) {
	snaps := map[string]*state.Snapshot{
		"staging": makeSnap(map[string]string{"key": "value"}),
	}
	engine := NewEngineWithDeps(stubLoader(snaps), drift.Detect)

	_, err := engine.Compare("state.json", "staging", "prod")
	if err == nil {
		t.Error("expected error for missing target environment, got nil")
	}
}

func TestCompareResult_String_NoDrift(t *testing.T) {
	r := &CompareResult{SourceEnv: "staging", TargetEnv: "prod", HasDrift: false}
	got := r.String()
	expected := "No drift detected between staging and prod"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
