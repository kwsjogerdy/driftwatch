package policy

import (
	"strings"
	"time"
)

// Match pairs a diff key with the rule that applies to it.
type Match struct {
	Key      string
	Rule     Rule
	Severity Severity
}

// Evaluator applies a set of policy rules to a collection of diff keys.
type Evaluator struct {
	rules []Rule
	now   func() time.Time
}

// NewEvaluator constructs an Evaluator with the given rules.
func NewEvaluator(rules []Rule) *Evaluator {
	return &Evaluator{
		rules: rules,
		now:   time.Now,
	}
}

// NewEvaluatorWithClock constructs an Evaluator with a custom clock (useful for testing).
func NewEvaluatorWithClock(rules []Rule, now func() time.Time) *Evaluator {
	return &Evaluator{rules: rules, now: now}
}

// Evaluate returns all rule matches for the provided keys.
// Expired rules are skipped. Each key is tested against all active rules;
// the highest-severity match wins when multiple rules apply.
func (e *Evaluator) Evaluate(keys []string) []Match {
	now := e.now()
	var matches []Match
	for _, key := range keys {
		var best *Match
		for _, rule := range e.rules {
			if rule.IsExpired(now) {
				continue
			}
			if !strings.HasPrefix(key, rule.KeyPrefix) {
				continue
			}
			m := Match{Key: key, Rule: rule, Severity: rule.Severity}
			if best == nil || severityRank(m.Severity) > severityRank(best.Severity) {
				copy := m
				best = &copy
			}
		}
		if best != nil {
			matches = append(matches, *best)
		}
	}
	return matches
}

func severityRank(s Severity) int {
	switch s {
	case SeverityCritical:
		return 2
	case SeverityWarning:
		return 1
	}
	return 0
}
