package history_test

import (
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/history"
)

func makeEntry(driftCount int, ts time.Time) history.Entry {
	return history.Entry{
		Timestamp:  ts,
		Source:     "prod",
		Target:     "staging",
		DriftCount: driftCount,
	}
}

func TestSummarise_Empty(t *testing.T) {
	s := history.Summarise(nil)
	if s.TotalRuns != 0 || s.DriftingRuns != 0 || s.CleanRuns != 0 {
		t.Errorf("expected all zeros for empty input, got %+v", s)
	}
}

func TestSummarise_AllClean(t *testing.T) {
	now := time.Now()
	entries := []history.Entry{
		makeEntry(0, now.Add(-2*time.Minute)),
		makeEntry(0, now.Add(-time.Minute)),
	}
	s := history.Summarise(entries)
	if s.TotalRuns != 2 {
		t.Errorf("TotalRuns: want 2, got %d", s.TotalRuns)
	}
	if s.CleanRuns != 2 {
		t.Errorf("CleanRuns: want 2, got %d", s.CleanRuns)
	}
	if s.DriftingRuns != 0 {
		t.Errorf("DriftingRuns: want 0, got %d", s.DriftingRuns)
	}
	if !s.LastDrift.IsZero() {
		t.Error("LastDrift should be zero when no drift occurred")
	}
}

func TestSummarise_MixedRuns(t *testing.T) {
	now := time.Now()
	t1 := now.Add(-3 * time.Minute)
	t2 := now.Add(-2 * time.Minute)
	t3 := now.Add(-time.Minute)
	entries := []history.Entry{
		makeEntry(0, t1),
		makeEntry(3, t2),
		makeEntry(0, t3),
	}
	s := history.Summarise(entries)
	if s.TotalRuns != 3 {
		t.Errorf("TotalRuns: want 3, got %d", s.TotalRuns)
	}
	if s.DriftingRuns != 1 {
		t.Errorf("DriftingRuns: want 1, got %d", s.DriftingRuns)
	}
	if s.CleanRuns != 2 {
		t.Errorf("CleanRuns: want 2, got %d", s.CleanRuns)
	}
	if !s.LastDrift.Equal(t2) {
		t.Errorf("LastDrift: want %v, got %v", t2, s.LastDrift)
	}
	if !s.LastChecked.Equal(t3) {
		t.Errorf("LastChecked: want %v, got %v", t3, s.LastChecked)
	}
}
