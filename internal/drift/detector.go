package drift

import (
	"github.com/driftwatch/internal/state"
)

// DifferenceKind describes the nature of a detected drift.
type DifferenceKind string

const (
	// KindValueMismatch means the key exists in both snapshots but values differ.
	KindValueMismatch DifferenceKind = "VALUE_MISMATCH"
	// KindMissing means the key is present in source but absent in target.
	KindMissing DifferenceKind = "MISSING_IN_TARGET"
	// KindExtra means the key is present in target but absent in source.
	KindExtra DifferenceKind = "EXTRA_IN_TARGET"
)

// Difference represents a single drift finding between two snapshots.
type Difference struct {
	Key         string
	SourceValue interface{}
	TargetValue interface{}
	Kind        DifferenceKind
}

// Detect compares source and target snapshots and returns all differences.
// An empty slice means the environments are in sync.
func Detect(source, target state.Snapshot) []Difference {
	var diffs []Difference

	for key, srcVal := range source.Values {
		tgtVal, exists := target.Values[key]
		if !exists {
			diffs = append(diffs, Difference{
				Key:         key,
				SourceValue: srcVal,
				TargetValue: nil,
				Kind:        KindMissing,
			})
			continue
		}
		if !valuesEqual(srcVal, tgtVal) {
			diffs = append(diffs, Difference{
				Key:         key,
				SourceValue: srcVal,
				TargetValue: tgtVal,
				Kind:        KindValueMismatch,
			})
		}
	}

	for key, tgtVal := range target.Values {
		if _, exists := source.Values[key]; !exists {
			diffs = append(diffs, Difference{
				Key:         key,
				SourceValue: nil,
				TargetValue: tgtVal,
				Kind:        KindExtra,
			})
		}
	}

	return diffs
}

// valuesEqual performs a simple equality check supporting primitive types.
func valuesEqual(a, b interface{}) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}
