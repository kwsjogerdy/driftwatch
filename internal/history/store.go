package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single recorded drift check result.
type Entry struct {
	Timestamp   time.Time         `json:"timestamp"`
	Source      string            `json:"source"`
	Target      string            `json:"target"`
	DriftCount  int               `json:"drift_count"`
	Drifts      map[string]string `json:"drifts,omitempty"`
}

// Store persists drift history entries to a directory.
type Store struct {
	dir string
}

// NewStore creates a Store that writes entries under dir.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("history: create dir %q: %w", dir, err)
	}
	return &Store{dir: dir}, nil
}

// Record writes an Entry as a JSON file named by its timestamp.
func (s *Store) Record(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	name := fmt.Sprintf("%d.json", e.Timestamp.UnixNano())
	path := filepath.Join(s.dir, name)

	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Errorf("history: marshal entry: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// List returns all stored entries sorted by filename (chronological).
func (s *Store) List() ([]Entry, error) {
	glob := filepath.Join(s.dir, "*.json")
	matches, err := filepath.Glob(glob)
	if err != nil {
		return nil, fmt.Errorf("history: glob: %w", err)
	}

	entries := make([]Entry, 0, len(matches))
	for _, m := range matches {
		data, err := os.ReadFile(m)
		if err != nil {
			return nil, fmt.Errorf("history: read %q: %w", m, err)
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("history: unmarshal %q: %w", m, err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
