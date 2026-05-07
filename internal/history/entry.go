package history

import "time"

// DriftStatus represents whether drift was detected in a run.
type DriftStatus string

const (
	StatusClean  DriftStatus = "clean"
	StatusDrifted DriftStatus = "drifted"
)

// Entry records the outcome of a single drift-check run.
type Entry struct {
	Timestamp  time.Time   `json:"timestamp"`
	EnvPair    string      `json:"env_pair"`
	Status     DriftStatus `json:"status"`
	DiffCount  int         `json:"diff_count"`
	StateFile  string      `json:"state_file"`
}

// IsClean returns true when no drift was found.
func (e Entry) IsClean() bool {
	return e.Status == StatusClean
}
