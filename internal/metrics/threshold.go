package metrics

import "fmt"

// BreachLevel indicates the severity of a threshold breach.
type BreachLevel string

const (
	BreachWarn     BreachLevel = "warn"
	BreachCritical BreachLevel = "critical"
)

// Breach describes a single threshold violation.
type Breach struct {
	Field  string
	Level  BreachLevel
	Value  float64
	Limit  float64
}

func (b Breach) String() string {
	return fmt.Sprintf("%s [%s] value=%.2f limit=%.2f", b.Field, b.Level, b.Value, b.Limit)
}

// Evaluator checks a MetricSummary against a set of thresholds.
type Evaluator struct {
	thresholds map[string]ThresholdConfig
}

// NewEvaluator creates an Evaluator from a pre-built threshold map.
func NewEvaluator(thresholds map[string]ThresholdConfig) *Evaluator {
	return &Evaluator{thresholds: thresholds}
}

// Evaluate returns all threshold breaches found in the given summary.
func (e *Evaluator) Evaluate(s Summary) []Breach {
	fields := map[string]float64{
		"drift_count":    float64(s.DriftCount),
		"clean_count":    float64(s.CleanCount),
		"total_runs":     float64(s.TotalRuns),
		"drift_rate":     s.DriftRate,
		"avg_diff_count": s.AvgDiffCount,
	}

	var breaches []Breach
	for key, val := range fields {
		cfg, ok := e.thresholds[key]
		if !ok {
			continue
		}
		switch {
		case val >= cfg.Critical:
			breaches = append(breaches, Breach{Field: key, Level: BreachCritical, Value: val, Limit: cfg.Critical})
		case val >= cfg.Warn:
			breaches = append(breaches, Breach{Field: key, Level: BreachWarn, Value: val, Limit: cfg.Warn})
		}
	}
	return breaches
}
