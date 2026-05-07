package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeFile(t *testing.T, dir, name string, modTime time.Time) {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte("{}\n"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := os.Chtimes(p, modTime, modTime); err != nil {
		t.Fatalf("chtimes: %v", err)
	}
}

func TestApply_EmptyDir_ReturnsZero(t *testing.T) {
	dir := t.TempDir()
	n, err := Apply(dir, DefaultRetentionPolicy())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 removed, got %d", n)
	}
}

func TestApply_NonExistentDir_ReturnsZero(t *testing.T) {
	n, err := Apply("/tmp/driftwatch-does-not-exist-xyz", DefaultRetentionPolicy())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 removed, got %d", n)
	}
}

func TestApply_RemovesOldFiles(t *testing.T) {
	dir := t.TempDir()
	now := time.Now()

	writeFile(t, dir, "recent.json", now.Add(-1*time.Hour))
	writeFile(t, dir, "old.json", now.Add(-40*24*time.Hour))

	policy := RetentionPolicy{MaxAge: 30 * 24 * time.Hour}
	n, err := Apply(dir, policy)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 removed, got %d", n)
	}

	if _, err := os.Stat(filepath.Join(dir, "recent.json")); err != nil {
		t.Error("recent.json should still exist")
	}
	if _, err := os.Stat(filepath.Join(dir, "old.json")); !os.IsNotExist(err) {
		t.Error("old.json should have been removed")
	}
}

func TestApply_EnforcesMaxEntries(t *testing.T) {
	dir := t.TempDir()
	now := time.Now()

	for i := 0; i < 5; i++ {
		name := filepath.Join(dir, time.Now().Format("20060102150405.000000000")+".json")
		mod := now.Add(-time.Duration(i) * time.Minute)
		if err := os.WriteFile(name, []byte("{}"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Chtimes(name, mod, mod); err != nil {
			t.Fatal(err)
		}
	}

	policy := RetentionPolicy{MaxEntries: 3}
	n, err := Apply(dir, policy)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 removed, got %d", n)
	}
}

func TestDefaultRetentionPolicy_Values(t *testing.T) {
	p := DefaultRetentionPolicy()
	if p.MaxAge != 30*24*time.Hour {
		t.Errorf("unexpected MaxAge: %v", p.MaxAge)
	}
	if p.MaxEntries != 100 {
		t.Errorf("unexpected MaxEntries: %d", p.MaxEntries)
	}
}
