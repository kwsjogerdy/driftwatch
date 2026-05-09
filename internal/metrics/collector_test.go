package metrics

import (
	"testing"
	"time"
)

func TestRecord_AndAll_RoundTrip(t *testing.T) {
	c := NewCollector()
	c.Record(RunMetrics{Environment: "prod", DriftCount: 2, Duration: time.Second})
	c.Record(RunMetrics{Environment: "staging", DriftCount: 0, Duration: 500 * time.Millisecond})

	all := c.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	if all[0].Environment != "prod" {
		t.Errorf("expected prod, got %s", all[0].Environment)
	}
}

func TestRecord_SetsTimestampIfZero(t *testing.T) {
	c := NewCollector()
	before := time.Now().UTC()
	c.Record(RunMetrics{Environment: "dev"})
	after := time.Now().UTC()

	all := c.All()
	ts := all[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", ts, before, after)
	}
}

func TestSummary_Empty(t *testing.T) {
	c := NewCollector()
	s := c.Summary()
	if s.TotalRuns != 0 {
		t.Errorf("expected 0 total runs, got %d", s.TotalRuns)
	}
}

func TestSummary_AggregatesCorrectly(t *testing.T) {
	c := NewCollector()
	c.Record(RunMetrics{DriftCount: 3, Duration: 2 * time.Second, HadError: false})
	c.Record(RunMetrics{DriftCount: 0, Duration: 1 * time.Second, HadError: true})
	c.Record(RunMetrics{DriftCount: 1, Duration: 3 * time.Second, HadError: false})

	s := c.Summary()
	if s.TotalRuns != 3 {
		t.Errorf("TotalRuns: want 3, got %d", s.TotalRuns)
	}
	if s.ErrorRuns != 1 {
		t.Errorf("ErrorRuns: want 1, got %d", s.ErrorRuns)
	}
	if s.DriftRuns != 2 {
		t.Errorf("DriftRuns: want 2, got %d", s.DriftRuns)
	}
	if s.TotalDriftKeys != 4 {
		t.Errorf("TotalDriftKeys: want 4, got %d", s.TotalDriftKeys)
	}
	wantAvg := 2 * time.Second
	if s.AvgDuration != wantAvg {
		t.Errorf("AvgDuration: want %v, got %v", wantAvg, s.AvgDuration)
	}
}
