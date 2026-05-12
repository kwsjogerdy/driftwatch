package rollup

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/diff"
)

var (
	t0 = time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	t1 = t0.Add(1 * time.Hour)
	t2 = t0.Add(2 * time.Hour)
	period = Period{Start: t0, End: t0.Add(24 * time.Hour)}
)

func makeEntry(src, tgt string, ts time.Time, keys ...string) Entry {
	diffs := make([]diff.Difference, 0, len(keys))
	for _, k := range keys {
		diffs = append(diffs, diff.Difference{Key: k, SourceValue: "a", TargetValue: "b"})
	}
	return Entry{Timestamp: ts, Source: src, Target: tgt, Diffs: diffs}
}

func TestAggregate_Empty(t *testing.T) {
	results := Aggregate(nil, period)
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestAggregate_AllClean(t *testing.T) {
	entries := []Entry{
		makeEntry("dev", "prod", t1),
		makeEntry("dev", "prod", t2),
	}
	results := Aggregate(entries, period)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if r.DriftedRuns != 0 {
		t.Errorf("expected 0 drifted runs, got %d", r.DriftedRuns)
	}
	if r.DriftRate != 0.0 {
		t.Errorf("expected drift rate 0, got %f", r.DriftRate)
	}
	if r.CleanRuns != 2 {
		t.Errorf("expected 2 clean runs, got %d", r.CleanRuns)
	}
}

func TestAggregate_MixedDrift(t *testing.T) {
	entries := []Entry{
		makeEntry("dev", "prod", t1, "db_host", "db_port"),
		makeEntry("dev", "prod", t2),
	}
	results := Aggregate(entries, period)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if r.DriftedRuns != 1 {
		t.Errorf("expected 1 drifted run, got %d", r.DriftedRuns)
	}
	if r.DriftRate != 0.5 {
		t.Errorf("expected drift rate 0.5, got %f", r.DriftRate)
	}
	if len(r.TopDriftKeys) == 0 {
		t.Error("expected top drift keys to be populated")
	}
}

func TestAggregate_OutsidePeriod_Excluded(t *testing.T) {
	outside := t0.Add(-1 * time.Hour)
	entries := []Entry{
		makeEntry("dev", "prod", outside, "key"),
		makeEntry("dev", "prod", t1),
	}
	results := Aggregate(entries, period)
	if results[0].TotalRuns != 1 {
		t.Errorf("expected 1 run within period, got %d", results[0].TotalRuns)
	}
}

func TestAggregate_MultipleEnvPairs(t *testing.T) {
	entries := []Entry{
		makeEntry("dev", "prod", t1),
		makeEntry("staging", "prod", t1, "timeout"),
	}
	results := Aggregate(entries, period)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestTopKeys_LimitedToN(t *testing.T) {
	counts := map[string]int{
		"a": 5, "b": 4, "c": 3, "d": 2, "e": 1, "f": 1,
	}
	keys := topKeys(counts, 3)
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "a" {
		t.Errorf("expected first key to be 'a', got %q", keys[0])
	}
}
