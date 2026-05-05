package drift

import (
	"fmt"

	"github.com/driftwatch/internal/state"
)

// Diff represents a single detected difference between two snapshots.
type Diff struct {
	Key           string
	Type          string // "mismatch", "missing", "extra"
	BaselineValue string
	TargetValue   string
}

// String returns a human-readable representation of the diff.
func (d Diff) String() string {
	return fmt.Sprintf("[%s] key=%q baseline=%q target=%q",
		d.Type, d.Key, d.BaselineValue, d.TargetValue)
}

// Detect compares a baseline snapshot against a target snapshot and
// returns all detected diffs.
func Detect(baseline, target state.Snapshot) []Diff {
	var diffs []Diff

	for k, bv := range baseline.Values {
		tv, ok := target.Values[k]
		if !ok {
			diffs = append(diffs, Diff{
				Key:           k,
				Type:          "missing",
				BaselineValue: bv,
				TargetValue:   "",
			})
			continue
		}
		if !valuesEqual(bv, tv) {
			diffs = append(diffs, Diff{
				Key:           k,
				Type:          "mismatch",
				BaselineValue: bv,
				TargetValue:   tv,
			})
		}
	}

	for k, tv := range target.Values {
		if _, ok := baseline.Values[k]; !ok {
			diffs = append(diffs, Diff{
				Key:           k,
				Type:          "extra",
				BaselineValue: "",
				TargetValue:   tv,
			})
		}
	}

	return diffs
}

// valuesEqual performs a case-sensitive equality check on two state values.
func valuesEqual(a, b string) bool {
	return a == b
}
