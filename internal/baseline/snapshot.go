package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Record represents a saved baseline snapshot for an environment pair.
type Record struct {
	SourceEnv string            `json:"source_env"`
	TargetEnv string            `json:"target_env"`
	CapturedAt time.Time        `json:"captured_at"`
	Values    map[string]string `json:"values"`
}

// Store manages persisted baseline records on disk.
type Store struct {
	dir string
}

// NewStore creates a Store rooted at dir, creating the directory if needed.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("baseline: create dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Save writes a baseline record to disk, keyed by source+target env names.
func (s *Store) Save(r Record) error {
	if r.CapturedAt.IsZero() {
		r.CapturedAt = time.Now().UTC()
	}
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	path := s.filePath(r.SourceEnv, r.TargetEnv)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write: %w", err)
	}
	return nil
}

// Load retrieves the baseline record for the given environment pair.
// Returns os.ErrNotExist if no baseline has been saved yet.
func (s *Store) Load(sourceEnv, targetEnv string) (Record, error) {
	path := s.filePath(sourceEnv, targetEnv)
	data, err := os.ReadFile(path)
	if err != nil {
		return Record{}, fmt.Errorf("baseline: read: %w", err)
	}
	var r Record
	if err := json.Unmarshal(data, &r); err != nil {
		return Record{}, fmt.Errorf("baseline: unmarshal: %w", err)
	}
	return r, nil
}

// Delete removes the stored baseline for the given environment pair.
func (s *Store) Delete(sourceEnv, targetEnv string) error {
	path := s.filePath(sourceEnv, targetEnv)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("baseline: delete: %w", err)
	}
	return nil
}

func (s *Store) filePath(sourceEnv, targetEnv string) string {
	name := fmt.Sprintf("%s__%s.json", sourceEnv, targetEnv)
	return filepath.Join(s.dir, name)
}
