package state_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/driftwatch/internal/state"
)

func writeTempState(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp state file: %v", err)
	}
	return path
}

func TestLoadFromFile_Valid(t *testing.T) {
	content := `{"environment":"staging","version":"1.2.0","resources":{"vpc":{"id":"vpc-abc123"}}}`
	path := writeTempState(t, content)

	sf, err := state.LoadFromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sf.Environment != "staging" {
		t.Errorf("expected environment %q, got %q", "staging", sf.Environment)
	}
	if sf.Version != "1.2.0" {
		t.Errorf("expected version %q, got %q", "1.2.0", sf.Version)
	}
	if _, ok := sf.Resources["vpc"]; !ok {
		t.Error("expected 'vpc' key in resources")
	}
}

func TestLoadFromFile_MissingEnvironment(t *testing.T) {
	content := `{"version":"1.0.0","resources":{}}`
	path := writeTempState(t, content)

	_, err := state.LoadFromFile(path)
	if err == nil {
		t.Fatal("expected error for missing environment, got nil")
	}
}

func TestLoadFromFile_InvalidJSON(t *testing.T) {
	path := writeTempState(t, `{not valid json`)

	_, err := state.LoadFromFile(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestLoadFromFile_NotFound(t *testing.T) {
	_, err := state.LoadFromFile("/nonexistent/path/state.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
