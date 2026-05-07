package history

import (
	"fmt"
	"log"
	"time"
)

// CleanerConfig holds configuration for the history cleaner.
type CleanerConfig struct {
	Retention RetentionPolicy
	DryRun    bool
	Logger    *log.Logger
}

// Cleaner applies retention policies to a history store.
type Cleaner struct {
	store  *Store
	cfg    CleanerConfig
	policy RetentionPolicy
}

// NewCleaner creates a Cleaner bound to the given Store.
func NewCleaner(store *Store, cfg CleanerConfig) *Cleaner {
	policy := cfg.Retention
	if policy.MaxAgeDays == 0 && policy.MaxEntries == 0 {
		policy = DefaultRetentionPolicy()
	}
	return &Cleaner{
		store:  store,
		cfg:    cfg,
		policy: policy,
	}
}

// CleanResult summarises what was removed during a clean run.
type CleanResult struct {
	Removed   int
	Retained  int
	RanAt     time.Time
	DryRun    bool
}

// Run executes the retention policy against the store directory.
// When DryRun is true, no files are deleted.
func (c *Cleaner) Run() (CleanResult, error) {
	result := CleanResult{
		RanAt:  time.Now().UTC(),
		DryRun: c.cfg.DryRun,
	}

	if c.cfg.DryRun {
		entries, err := c.store.List()
		if err != nil {
			return result, fmt.Errorf("cleaner: listing entries: %w", err)
		}
		result.Retained = len(entries)
		if c.cfg.Logger != nil {
			c.cfg.Logger.Printf("dry-run: would evaluate %d entries against retention policy", len(entries))
		}
		return result, nil
	}

	removed, err := Apply(c.store.dir, c.policy)
	if err != nil {
		return result, fmt.Errorf("cleaner: applying retention policy: %w", err)
	}

	result.Removed = removed

	entries, err := c.store.List()
	if err != nil {
		return result, fmt.Errorf("cleaner: listing remaining entries: %w", err)
	}
	result.Retained = len(entries)

	if c.cfg.Logger != nil {
		c.cfg.Logger.Printf("cleaner: removed %d entries, %d retained", removed, result.Retained)
	}

	return result, nil
}
