package output

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewWriter_Stdout(t *testing.T) {
	w, err := NewWriter(DestStdout, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.Destination() != DestStdout {
		t.Errorf("expected DestStdout, got %q", w.Destination())
	}
}

func TestNewWriter_Stderr(t *testing.T) {
	w, err := NewWriter(DestStderr, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.Destination() != DestStderr {
		t.Errorf("expected DestStderr, got %q", w.Destination())
	}
}

func TestNewWriter_File_WritesContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")

	w, err := NewWriter(DestFile, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer w.Close()

	_, err = w.Write([]byte("hello drift"))
	if err != nil {
		t.Fatalf("write error: %v", err)
	}
	w.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}
	if string(data) != "hello drift" {
		t.Errorf("expected %q, got %q", "hello drift", string(data))
	}
}

func TestNewWriter_File_EmptyPath(t *testing.T) {
	_, err := NewWriter(DestFile, "")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestNewWriter_UnknownDestination(t *testing.T) {
	_, err := NewWriter(Destination("syslog"), "")
	if err == nil {
		t.Fatal("expected error for unknown destination, got nil")
	}
}

func TestNewWriter_File_InvalidPath(t *testing.T) {
	_, err := NewWriter(DestFile, "/nonexistent/dir/out.txt")
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}
