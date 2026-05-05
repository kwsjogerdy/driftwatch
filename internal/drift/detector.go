package drift

import (
	"fmt"

	"github.com/driftwatch/internal/state"
)

// Result holds the outcome of a drift comparison between two environments.
type Result struct {
	Baseline    string
	Target      string
	Drifted     bool
	Mismatches  []Mismatch
}

// Mismatch describes a single resource key that differs between environments.
type Mismatch struct {
	Key           string
	BaselineValue string
	TargetValue   string
}

// Detect compares two environment snapshots and returns a Result describing
// any drift found between them.
func Detect(baseline, target *state.Snapshot) (*Result, error) {
	if baseline == nil {
		return nil, fmt.Errorf("baseline snapshot must not be nil")
	}
	if target == nil {
		return nil, fmt.Errorf("target snapshot must not be nil")
	}

	result := &Result{
		Baseline: baseline.Environment,
		Target:   target.Environment,
	}

	seen := make(map[string]bool)

	for key, bVal := range baseline.Resources {
		seen[key] = true
		tVal, ok := target.Resources[key]
		if !ok {
			result.Mismatches = append(result.Mismatches, Mismatch{
				Key:           key,
				BaselineValue: bVal,
				TargetValue:   "<missing>",
			})
			continue
		}
		if bVal != tVal {
			result.Mismatches = append(result.Mismatches, Mismatch{
				Key:           key,
				BaselineValue: bVal,
				TargetValue:   tVal,
			})
		}
	}

	for key, tVal := range target.Resources {
		if seen[key] {
			continue
		}
		result.Mismatches = append(result.Mismatches, Mismatch{
			Key:           key,
			BaselineValue: "<missing>",
			TargetValue:   tVal,
		})
	}

	result.Drifted = len(result.Mismatches) > 0
	return result, nil
}
