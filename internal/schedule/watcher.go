package schedule

import (
	"context"
	"log"
	"time"

	"github.com/driftwatch/internal/compare"
)

// WatchConfig holds configuration for the periodic watcher.
type WatchConfig struct {
	Interval    time.Duration
	SourceEnv   string
	TargetEnv   string
	StateFile   string
}

// Watcher periodically runs drift comparisons on a schedule.
type Watcher struct {
	cfg    WatchConfig
	engine *compare.Engine
}

// NewWatcher creates a Watcher with the given config and compare engine.
func NewWatcher(cfg WatchConfig, engine *compare.Engine) *Watcher {
	return &Watcher{
		cfg:    cfg,
		engine: engine,
	}
}

// Run starts the watch loop, blocking until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	log.Printf("[watcher] starting drift watch: %s -> %s every %s",
		w.cfg.SourceEnv, w.cfg.TargetEnv, w.cfg.Interval)

	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("[watcher] stopping")
			return ctx.Err()
		case <-ticker.C:
			w.runOnce()
		}
	}
}

func (w *Watcher) runOnce() {
	diffs, err := w.engine.Compare(w.cfg.StateFile, w.cfg.SourceEnv, w.cfg.TargetEnv)
	if err != nil {
		log.Printf("[watcher] compare error: %v", err)
		return
	}
	if len(diffs) == 0 {
		log.Printf("[watcher] no drift detected between %s and %s",
			w.cfg.SourceEnv, w.cfg.TargetEnv)
		return
	}
	log.Printf("[watcher] drift detected: %d difference(s) between %s and %s",
		len(diffs), w.cfg.SourceEnv, w.cfg.TargetEnv)
}
