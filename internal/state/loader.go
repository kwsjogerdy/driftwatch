package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// StateFile represents a parsed infrastructure state file.
type StateFile struct {
	Environment string                 `json:"environment"`
	Version     string                 `json:"version"`
	Resources   map[string]interface{} `json:"resources"`
	Checksum    string                 `json:"checksum,omitempty"`
}

// LoadFromFile reads and parses a JSON state file from the given path.
func LoadFromFile(path string) (*StateFile, error) {
	path = filepath.Clean(path)

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("state: failed to open file %q: %w", path, err)
	}
	defer f.Close()

	var sf StateFile
	decoder := json.NewDecoder(f)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&sf); err != nil {
		return nil, fmt.Errorf("state: failed to decode file %q: %w", path, err)
	}

	if err := sf.validate(); err != nil {
		return nil, fmt.Errorf("state: invalid state file %q: %w", path, err)
	}

	return &sf, nil
}

// validate checks that required fields are present.
func (s *StateFile) validate() error {
	if s.Environment == "" {
		return fmt.Errorf("missing required field: environment")
	}
	if s.Version == "" {
		return fmt.Errorf("missing required field: version")
	}
	if s.Resources == nil {
		return fmt.Errorf("missing required field: resources")
	}
	return nil
}
