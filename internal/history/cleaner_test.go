package history

import (
	"log"
	"os"
	"testing"
	"time"
)

func makeCleanerStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	store, err := NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return store
}

func recordEntries(t *testing.T, store *Store, n int, age time.Duration) {
	t.Helper()
	for i := 0; i < n; i++ {
		e := Entry{
			Timestamp: time.Now().UTC().Add(-age),
			Source:    "prod",
			Target:    "staging",
			DriftCount: i,
		}
		if err := store.Record(e); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}
}

func TestCleaner_DryRun_DoesNotRemove(t *testing.T) {
	store := makeCleanerStore(t)
	recordEntries(t, store, 5, 0)

	cleaner := NewCleaner(store, CleanerConfig{
		Retention: RetentionPolicy{MaxEntries: 2},
		DryRun:    true,
		Logger:    log.New(os.Stderr, "", 0),
	})

	result, err := cleaner.Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.Removed != 0 {
		t.Errorf("expected 0 removed in dry-run, got %d", result.Removed)
	}
	if result.Retained != 5 {
		t.Errorf("expected 5 retained in dry-run, got %d", result.Retained)
	}
	if !result.DryRun {
		t.Error("expected DryRun flag to be true")
	}
}

func TestCleaner_Run_RemovesExcessEntries(t *testing.T) {
	store := makeCleanerStore(t)
	recordEntries(t, store, 6, 0)

	cleaner := NewCleaner(store, CleanerConfig{
		Retention: RetentionPolicy{MaxEntries: 3},
	})

	result, err := cleaner.Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.Removed != 3 {
		t.Errorf("expected 3 removed, got %d", result.Removed)
	}
	if result.Retained != 3 {
		t.Errorf("expected 3 retained, got %d", result.Retained)
	}
}

func TestCleaner_DefaultPolicy_WhenZero(t *testing.T) {
	store := makeCleanerStore(t)
	cleaner := NewCleaner(store, CleanerConfig{})
	if cleaner.policy.MaxAgeDays == 0 && cleaner.policy.MaxEntries == 0 {
		t.Error("expected default policy to be applied when both fields are zero")
	}
}

func TestCleaner_Run_SetsTimestamp(t *testing.T) {
	store := makeCleanerStore(t)
	before := time.Now().UTC().Add(-time.Second)
	cleaner := NewCleaner(store, CleanerConfig{
		Retention: RetentionPolicy{MaxEntries: 10},
	})
	result, err := cleaner.Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.RanAt.Before(before) {
		t.Errorf("expected RanAt to be recent, got %v", result.RanAt)
	}
}
