package baseline_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"driftwatch/internal/baseline"
)

func TestSave_AndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	store, err := baseline.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	rec := baseline.Record{
		SourceEnv:  "staging",
		TargetEnv:  "production",
		CapturedAt: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
		Values:     map[string]string{"db_host": "db.prod", "replicas": "3"},
	}

	if err := store.Save(rec); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := store.Load("staging", "production")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if got.SourceEnv != rec.SourceEnv || got.TargetEnv != rec.TargetEnv {
		t.Errorf("env mismatch: got %s/%s", got.SourceEnv, got.TargetEnv)
	}
	if got.Values["db_host"] != "db.prod" {
		t.Errorf("values not preserved")
	}
}

func TestLoad_NotFound_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	store, _ := baseline.NewStore(dir)

	_, err := store.Load("dev", "prod")
	if err == nil {
		t.Fatal("expected error for missing baseline")
	}
}

func TestSave_SetsTimestampIfZero(t *testing.T) {
	dir := t.TempDir()
	store, _ := baseline.NewStore(dir)

	rec := baseline.Record{
		SourceEnv: "dev",
		TargetEnv: "staging",
		Values:    map[string]string{"key": "val"},
	}

	before := time.Now().UTC()
	if err := store.Save(rec); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, _ := store.Load("dev", "staging")
	if got.CapturedAt.Before(before) {
		t.Errorf("timestamp not set: %v", got.CapturedAt)
	}
}

func TestDelete_RemovesRecord(t *testing.T) {
	dir := t.TempDir()
	store, _ := baseline.NewStore(dir)

	rec := baseline.Record{SourceEnv: "a", TargetEnv: "b", Values: map[string]string{}}
	_ = store.Save(rec)

	if err := store.Delete("a", "b"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err := store.Load("a", "b")
	if !errors.Is(err, os.ErrNotExist) && err == nil {
		t.Error("expected record to be gone after delete")
	}
}

func TestDelete_NonExistent_NoError(t *testing.T) {
	dir := t.TempDir()
	store, _ := baseline.NewStore(dir)

	if err := store.Delete("x", "y"); err != nil {
		t.Errorf("expected no error deleting non-existent baseline, got: %v", err)
	}
}
