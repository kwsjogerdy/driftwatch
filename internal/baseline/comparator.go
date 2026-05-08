package baseline

import (
	"fmt"

	"driftwatch/internal/drift"
	"driftwatch/internal/state"
)

// Comparator compares a live environment snapshot against a saved baseline.
type Comparator struct {
	store *Store
}

// NewComparator creates a Comparator backed by the given Store.
func NewComparator(store *Store) *Comparator {
	return &Comparator{store: store}
}

// CompareResult holds the outcome of a baseline comparison.
type CompareResult struct {
	SourceEnv  string
	TargetEnv  string
	Diffs      []drift.Diff
	BaselineAge string
}

// Compare loads the saved baseline for the env pair and diffs it against
// the provided live snapshot. Returns CompareResult with any detected diffs.
func (c *Comparator) Compare(sourceEnv, targetEnv string, live state.Snapshot) (CompareResult, error) {
	rec, err := c.store.Load(sourceEnv, targetEnv)
	if err != nil {
		return CompareResult{}, fmt.Errorf("comparator: load baseline: %w", err)
	}

	baselineSnap := state.Snapshot{
		Environment: rec.TargetEnv,
		Values:      rec.Values,
	}

	diffs := drift.Detect(baselineSnap, live)

	age := "unknown"
	if !rec.CapturedAt.IsZero() {
		age = rec.CapturedAt.Format("2006-01-02T15:04:05Z")
	}

	return CompareResult{
		SourceEnv:   sourceEnv,
		TargetEnv:   targetEnv,
		Diffs:       diffs,
		BaselineAge: age,
	}, nil
}

// CaptureBaseline saves the current live snapshot as the new baseline.
func (c *Comparator) CaptureBaseline(sourceEnv, targetEnv string, live state.Snapshot) error {
	rec := Record{
		SourceEnv: sourceEnv,
		TargetEnv: targetEnv,
		Values:    live.Values,
	}
	if err := c.store.Save(rec); err != nil {
		return fmt.Errorf("comparator: capture baseline: %w", err)
	}
	return nil
}
