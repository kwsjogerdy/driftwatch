package history_test

import (
	"os"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/history"
)

func TestNewStore_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	subdir := dir + "/history"
	_, err := history.NewStore(subdir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(subdir); os.IsNotExist(err) {
		t.Fatal("expected directory to be created")
	}
}

func TestRecord_AndList_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	store, err := history.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	e := history.Entry{
		Timestamp:  time.Now().UTC(),
		Source:     "prod",
		Target:     "staging",
		DriftCount: 2,
		Drifts:     map[string]string{"key1": "value mismatch"},
	}
	if err := store.Record(e); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	got := entries[0]
	if got.Source != e.Source {
		t.Errorf("source: want %q, got %q", e.Source, got.Source)
	}
	if got.DriftCount != e.DriftCount {
		t.Errorf("drift_count: want %d, got %d", e.DriftCount, got.DriftCount)
	}
}

func TestList_EmptyDir_ReturnsEmpty(t *testing.T) {
	store, _ := history.NewStore(t.TempDir())
	entries, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty list, got %d entries", len(entries))
	}
}

func TestRecord_SetsTimestampIfZero(t *testing.T) {
	store, _ := history.NewStore(t.TempDir())
	e := history.Entry{Source: "a", Target: "b"}
	if err := store.Record(e); err != nil {
		t.Fatalf("Record: %v", err)
	}
	entries, _ := store.List()
	if entries[0].Timestamp.IsZero() {
		t.Error("expected timestamp to be set automatically")
	}
}
