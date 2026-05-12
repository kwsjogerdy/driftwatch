package suppress

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/diff"
)

func makeDiff(key string) diff.Difference {
	return diff.Difference{Key: key, SourceValue: "a", TargetValue: "b"}
}

func activeRule(prefix, src, tgt string) Rule {
	return Rule{
		ID:        prefix + "-rule",
		KeyPrefix: prefix,
		SourceEnv: src,
		TargetEnv: tgt,
		ExpiresAt: time.Now().Add(time.Hour),
	}
}

func TestApply_NoRules_ReturnAll(t *testing.T) {
	f := NewFilter(nil)
	diffs := []diff.Difference{makeDiff("db.host"), makeDiff("app.port")}
	out := f.Apply(diffs, "staging", "prod")
	if len(out) != 2 {
		t.Fatalf("expected 2 diffs, got %d", len(out))
	}
}

func TestApply_SuppressesMatchingKey(t *testing.T) {
	rules := []Rule{activeRule("db.", "staging", "prod")}
	f := NewFilter(rules)
	diffs := []diff.Difference{makeDiff("db.host"), makeDiff("app.port")}
	out := f.Apply(diffs, "staging", "prod")
	if len(out) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(out))
	}
	if out[0].Key != "app.port" {
		t.Errorf("unexpected key %q", out[0].Key)
	}
}

func TestApply_ExpiredRule_DoesNotSuppress(t *testing.T) {
	r := activeRule("db.", "staging", "prod")
	r.ExpiresAt = time.Now().Add(-time.Minute)
	f := NewFilter([]Rule{r})
	diffs := []diff.Difference{makeDiff("db.host")}
	out := f.Apply(diffs, "staging", "prod")
	if len(out) != 1 {
		t.Fatalf("expected expired rule to not suppress, got %d diffs", len(out))
	}
}

func TestApply_WrongEnv_DoesNotSuppress(t *testing.T) {
	rules := []Rule{activeRule("db.", "staging", "prod")}
	f := NewFilter(rules)
	diffs := []diff.Difference{makeDiff("db.host")}
	out := f.Apply(diffs, "dev", "prod")
	if len(out) != 1 {
		t.Fatal("rule should not match different source env")
	}
}

func TestActiveCount(t *testing.T) {
	expired := activeRule("x.", "", "")
	expired.ExpiresAt = time.Now().Add(-time.Hour)
	rules := []Rule{activeRule("a.", "", ""), expired}
	f := NewFilter(rules)
	if got := f.ActiveCount(); got != 1 {
		t.Fatalf("expected 1 active rule, got %d", got)
	}
}
