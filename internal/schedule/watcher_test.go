package schedule_test

import (
	"context"
	"testing"
	"time"

	"github.com/driftwatch/internal/compare"
	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/schedule"
	"github.com/driftwatch/internal/state"
)

type stubLoader struct {
	snaps map[string]state.Snapshot
}

func (s *stubLoader) LoadFromFile(path, env string) (state.Snapshot, error) {
	snap, ok := s.snaps[env]
	if !ok {
		return state.Snapshot{}, nil
	}
	return snap, nil
}

func makeEngine(snaps map[string]state.Snapshot) *compare.Engine {
	loader := &stubLoader{snaps: snaps}
	detector := drift.NewDetector()
	return compare.NewEngineWithDeps(loader, detector)
}

func TestWatcher_RunCancelsCleanly(t *testing.T) {
	engine := makeEngine(map[string]state.Snapshot{
		"prod":    {Environment: "prod", Values: map[string]string{"key": "val"}},
		"staging": {Environment: "staging", Values: map[string]string{"key": "val"}},
	})

	cfg := schedule.WatchConfig{
		Interval:  50 * time.Millisecond,
		SourceEnv: "prod",
		TargetEnv: "staging",
		StateFile: "state.json",
	}

	w := schedule.NewWatcher(cfg, engine)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	err := w.Run(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}

func TestWatcher_RunTicksMultipleTimes(t *testing.T) {
	engine := makeEngine(map[string]state.Snapshot{
		"prod":    {Environment: "prod", Values: map[string]string{"a": "1"}},
		"staging": {Environment: "staging", Values: map[string]string{"a": "2"}},
	})

	cfg := schedule.WatchConfig{
		Interval:  30 * time.Millisecond,
		SourceEnv: "prod",
		TargetEnv: "staging",
		StateFile: "state.json",
	}

	w := schedule.NewWatcher(cfg, engine)
	ctx, cancel := context.WithTimeout(context.Background(), 110*time.Millisecond)
	defer cancel()

	// Should tick ~3 times without panicking
	err := w.Run(ctx)
	if err == nil {
		t.Error("expected context error, got nil")
	}
}
