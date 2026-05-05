package compare

import (
	"fmt"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/state"
)

// CompareResult holds the outcome of comparing two environments.
type CompareResult struct {
	SourceEnv string
	TargetEnv string
	Diffs     []drift.Diff
	HasDrift  bool
}

// String returns a human-readable summary of the compare result.
func (r CompareResult) String() string {
	if !r.HasDrift {
		return fmt.Sprintf("No drift detected between %s and %s", r.SourceEnv, r.TargetEnv)
	}
	return fmt.Sprintf("%d drift(s) detected between %s and %s", len(r.Diffs), r.SourceEnv, r.TargetEnv)
}

// Engine performs environment comparisons using loaded state snapshots.
type Engine struct {
	loader  func(path, env string) (*state.Snapshot, error)
	detector func(source, target *state.Snapshot) []drift.Diff
}

// NewEngine creates a new Engine with default loader and detector.
func NewEngine() *Engine {
	return &Engine{
		loader:  state.LoadFromFile,
		detector: drift.Detect,
	}
}

// NewEngineWithDeps creates an Engine with injected dependencies (useful for testing).
func NewEngineWithDeps(
	loader func(path, env string) (*state.Snapshot, error),
	detector func(source, target *state.Snapshot) []drift.Diff,
) *Engine {
	return &Engine{loader: loader, detector: detector}
}

// Compare loads snapshots for two environments from the given file and detects drift.
func (e *Engine) Compare(filePath, sourceEnv, targetEnv string) (*CompareResult, error) {
	source, err := e.loader(filePath, sourceEnv)
	if err != nil {
		return nil, fmt.Errorf("loading source environment %q: %w", sourceEnv, err)
	}

	target, err := e.loader(filePath, targetEnv)
	if err != nil {
		return nil, fmt.Errorf("loading target environment %q: %w", targetEnv, err)
	}

	diffs := e.detector(source, target)

	return &CompareResult{
		SourceEnv: sourceEnv,
		TargetEnv: targetEnv,
		Diffs:     diffs,
		HasDrift:  len(diffs) > 0,
	}, nil
}
