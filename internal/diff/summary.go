package diff

import "fmt"

// Severity represents the importance of a detected difference.
type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

// Difference describes a single detected drift between two environments.
type Difference struct {
	Key      string
	Source   interface{}
	Target   interface{}
	Kind     string // "mismatch", "missing", "extra"
	Severity Severity
}

// Summary aggregates the results of a drift comparison.
type Summary struct {
	SourceEnv   string
	TargetEnv   string
	Differences []Difference
	TotalDrift  int
	HasCritical bool
	HasWarning  bool
}

// NewSummary builds a Summary from a list of differences.
func NewSummary(source, target string, diffs []Difference) Summary {
	s := Summary{
		SourceEnv:   source,
		TargetEnv:   target,
		Differences: diffs,
		TotalDrift:  len(diffs),
	}
	for _, d := range diffs {
		switch d.Severity {
		case SeverityCritical:
			s.HasCritical = true
		case SeverityWarning:
			s.HasWarning = true
		}
	}
	return s
}

// Label returns a short human-readable status string.
func (s Summary) Label() string {
	if s.TotalDrift == 0 {
		return "clean"
	}
	if s.HasCritical {
		return "critical"
	}
	if s.HasWarning {
		return "warning"
	}
	return "info"
}

// String implements the Stringer interface.
func (s Summary) String() string {
	return fmt.Sprintf("[%s] %s → %s: %d difference(s)",
		s.Label(), s.SourceEnv, s.TargetEnv, s.TotalDrift)
}
