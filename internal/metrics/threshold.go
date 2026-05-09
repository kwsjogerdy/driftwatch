package metrics

import "fmt"

// ThresholdLevel represents the severity of a threshold breach.
type ThresholdLevel string

const (
	ThresholdWarn     ThresholdLevel = "warn"
	ThresholdCritical ThresholdLevel = "critical"
)

// Threshold defines a limit for a named metric field.
type Threshold struct {
	Field    string
	WarnAt   float64
	CritAt   float64
}

// Breach describes a single threshold violation.
type Breach struct {
	Field   string
	Value   float64
	Limit   float64
	Level   ThresholdLevel
}

func (b Breach) String() string {
	return fmt.Sprintf("%s: value %.2f exceeds %s limit %.2f", b.Field, b.Value, b.Level, b.Limit)
}

// Evaluator checks a Summary against a set of Thresholds.
type Evaluator struct {
	thresholds []Threshold
}

// NewEvaluator creates an Evaluator with the given thresholds.
func NewEvaluator(thresholds []Threshold) *Evaluator {
	return &Evaluator{thresholds: thresholds}
}

// Evaluate returns any threshold breaches found in the given Summary.
func (e *Evaluator) Evaluate(s Summary) []Breach {
	fields := map[string]float64{
		"total_runs":    float64(s.TotalRuns),
		"drift_count":   float64(s.DriftCount),
		"error_count":   float64(s.ErrorCount),
		"drift_percent": s.DriftPercent,
	}

	var breaches []Breach
	for _, t := range e.thresholds {
		v, ok := fields[t.Field]
		if !ok {
			continue
		}
		switch {
		case t.CritAt > 0 && v >= t.CritAt:
			breaches = append(breaches, Breach{Field: t.Field, Value: v, Limit: t.CritAt, Level: ThresholdCritical})
		case t.WarnAt > 0 && v >= t.WarnAt:
			breaches = append(breaches, Breach{Field: t.Field, Value: v, Limit: t.WarnAt, Level: ThresholdWarn})
		}
	}
	return breaches
}
