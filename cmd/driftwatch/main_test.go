package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "driftwatch.json")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return path
}

func TestMain_MissingConfigFlag(t *testing.T) {
	if os.Getenv("RUN_MAIN_TEST") != "1" {
		cmd := exec.Command(os.Args[0], "-test.run=TestMain_MissingConfigFlag")
		cmd.Env = append(os.Environ(), "RUN_MAIN_TEST=1")
		err := cmd.Run()
		if err == nil {
			t.Fatal("expected non-zero exit for missing config")
		}
		return
	}
	os.Args = []string{"driftwatch", "--config", "/nonexistent/path.json"}
	main()
}

func TestMain_InvalidConfig(t *testing.T) {
	if os.Getenv("RUN_INVALID_TEST") != "1" {
		cmd := exec.Command(os.Args[0], "-test.run=TestMain_InvalidConfig")
		cmd.Env = append(os.Environ(), "RUN_INVALID_TEST=1")
		err := cmd.Run()
		if err == nil {
			t.Fatal("expected non-zero exit for invalid config")
		}
		return
	}
	path := writeTempConfig(t, `{"state_file": ""}`)
	os.Args = []string{"driftwatch", "--config", path}
	main()
}
