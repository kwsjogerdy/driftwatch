package history

import (
	"testing"
	"time"
)

func makeQueryEntry(envPair string, status DriftStatus, hoursAgo float64) Entry {
	return Entry{
		Timestamp: time.Now().Add(-time.Duration(hoursAgo * float64(time.Hour))),
		EnvPair:   envPair,
		Status:    status,
		DiffCount: 1,
	}
}

func TestQuery_NoFilter_ReturnsAll(t *testing.T) {
	entries := []Entry{
		makeQueryEntry("a->b", StatusClean, 1),
		makeQueryEntry("a->b", StatusDrifted, 2),
		makeQueryEntry("c->d", StatusClean, 3),
	}
	got := Query(entries, QueryOptions{})
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
}

func TestQuery_FilterByEnvPair(t *testing.T) {
	entries := []Entry{
		makeQueryEntry("a->b", StatusClean, 1),
		makeQueryEntry("c->d", StatusDrifted, 2),
	}
	got := Query(entries, QueryOptions{EnvPair: "a->b"})
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if got[0].EnvPair != "a->b" {
		t.Errorf("unexpected env pair: %s", got[0].EnvPair)
	}
}

func TestQuery_FilterBySince(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		{Timestamp: now.Add(-30 * time.Minute), EnvPair: "a->b", Status: StatusClean},
		{Timestamp: now.Add(-3 * time.Hour), EnvPair: "a->b", Status: StatusDrifted},
	}
	got := Query(entries, QueryOptions{Since: now.Add(-1 * time.Hour)})
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if got[0].Status != StatusClean {
		t.Errorf("expected clean entry")
	}
}

func TestQuery_MaxItems(t *testing.T) {
	entries := []Entry{
		makeQueryEntry("a->b", StatusClean, 1),
		makeQueryEntry("a->b", StatusClean, 2),
		makeQueryEntry("a->b", StatusDrifted, 3),
	}
	got := Query(entries, QueryOptions{MaxItems: 2})
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestQuery_SortedDescending(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		{Timestamp: now.Add(-3 * time.Hour), Status: StatusDrifted},
		{Timestamp: now.Add(-1 * time.Hour), Status: StatusClean},
		{Timestamp: now.Add(-2 * time.Hour), Status: StatusDrifted},
	}
	got := Query(entries, QueryOptions{})
	if !got[0].Timestamp.After(got[1].Timestamp) {
		t.Errorf("expected descending order")
	}
}
