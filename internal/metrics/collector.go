package metrics

import (
	"sync"
	"time"
)

// RunMetrics holds statistics for a single drift-check run.
type RunMetrics struct {
	Environment string
	DriftCount  int
	Duration    time.Duration
	Timestamp   time.Time
	HadError    bool
}

// Collector accumulates run metrics in memory.
type Collector struct {
	mu      sync.Mutex
	entries []RunMetrics
}

// NewCollector returns an initialised Collector.
func NewCollector() *Collector {
	return &Collector{}
}

// Record appends a RunMetrics entry. If Timestamp is zero it is set to now.
func (c *Collector) Record(m RunMetrics) {
	if m.Timestamp.IsZero() {
		m.Timestamp = time.Now().UTC()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = append(c.entries, m)
}

// All returns a copy of every recorded entry.
func (c *Collector) All() []RunMetrics {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]RunMetrics, len(c.entries))
	copy(out, c.entries)
	return out
}

// Summary returns aggregate statistics over all recorded entries.
func (c *Collector) Summary() Summary {
	c.mu.Lock()
	defer c.mu.Unlock()
	return summarise(c.entries)
}
