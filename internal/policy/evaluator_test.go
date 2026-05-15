package policy

import (
	"testing"
	"time"
)

var fixedNow = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func clockFn() time.Time { return fixedNow }

func makeEvaluator(rules []Rule) *Evaluator {
	return NewEvaluatorWithClock(rules, clockFn)
}

func TestEvaluate_NoRules_ReturnsEmpty(t *testing.T) {
	e := makeEvaluator(nil)
	matches := e.Evaluate([]string{"app.version"})
	if len(matches) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(matches))
	}
}

func TestEvaluate_MatchingPrefix_ReturnsMatch(t *testing.T) {
	rules := []Rule{{ID: "r1", KeyPrefix: "app.", Severity: SeverityWarning}}
	e := makeEvaluator(rules)
	matches := e.Evaluate([]string{"app.version", "db.host"})
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].Key != "app.version" {
		t.Errorf("unexpected key %q", matches[0].Key)
	}
}

func TestEvaluate_ExpiredRule_Skipped(t *testing.T) {
	rules := []Rule{
		{ID: "r1", KeyPrefix: "app.", Severity: SeverityWarning, ExpiresAt: fixedNow.Add(-time.Hour)},
	}
	e := makeEvaluator(rules)
	matches := e.Evaluate([]string{"app.version"})
	if len(matches) != 0 {
		t.Fatalf("expected expired rule to be skipped, got %d matches", len(matches))
	}
}

func TestEvaluate_HigherSeverityWins(t *testing.T) {
	rules := []Rule{
		{ID: "r1", KeyPrefix: "app.", Severity: SeverityWarning},
		{ID: "r2", KeyPrefix: "app.version", Severity: SeverityCritical},
	}
	e := makeEvaluator(rules)
	matches := e.Evaluate([]string{"app.version"})
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].Severity != SeverityCritical {
		t.Errorf("expected critical, got %s", matches[0].Severity)
	}
}

func TestEvaluate_NoMatchingPrefix_ReturnsEmpty(t *testing.T) {
	rules := []Rule{{ID: "r1", KeyPrefix: "infra.", Severity: SeverityWarning}}
	e := makeEvaluator(rules)
	matches := e.Evaluate([]string{"app.version"})
	if len(matches) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(matches))
	}
}
