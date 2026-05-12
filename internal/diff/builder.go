package diff

import "fmt"

// Builder constructs a list of Differences from two flat key-value maps.
type Builder struct {
	criticalKeys map[string]bool
}

// NewBuilder creates a Builder. Keys listed in criticalKeys will be assigned
// SeverityCritical when they differ or are absent.
func NewBuilder(criticalKeys []string) *Builder {
	cm := make(map[string]bool, len(criticalKeys))
	for _, k := range criticalKeys {
		cm[k] = true
	}
	return &Builder{criticalKeys: cm}
}

// Build compares source and target maps and returns all detected differences.
func (b *Builder) Build(source, target map[string]interface{}) []Difference {
	var diffs []Difference

	for k, sv := range source {
		tv, ok := target[k]
		if !ok {
			diffs = append(diffs, Difference{
				Key:      k,
				Source:   sv,
				Target:   nil,
				Kind:     "missing",
				Severity: b.severity(k),
			})
			continue
		}
		if !valuesEqual(sv, tv) {
			diffs = append(diffs, Difference{
				Key:      k,
				Source:   sv,
				Target:   tv,
				Kind:     "mismatch",
				Severity: b.severity(k),
			})
		}
	}

	for k, tv := range target {
		if _, ok := source[k]; !ok {
			diffs = append(diffs, Difference{
				Key:      k,
				Source:   nil,
				Target:   tv,
				Kind:     "extra",
				Severity: SeverityInfo,
			})
		}
	}

	return diffs
}

func (b *Builder) severity(key string) Severity {
	if b.criticalKeys[key] {
		return SeverityCritical
	}
	return SeverityWarning
}

func valuesEqual(a, b interface{}) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}
