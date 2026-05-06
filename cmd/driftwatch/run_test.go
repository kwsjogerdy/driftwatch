package main

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/compare"
	"github.com/example/driftwatch/internal/config"
	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/state"
)

func makeTestConfig(envs []string, format string) *config.AppConfig {
	return &config.AppConfig{
		StateFile: "state.json",
		Output:    config.OutputConfig{Format: format, Destination: "stdout"},
		Schedule:  config.ScheduleConfig{Environments: envs, Interval: "10s"},
	}
}

func stubLoaderFn(snapshots map[string]*state.Snapshot) func(string, string) (*state.Snapshot, error) {
	return func(_ string, env string) (*state.Snapshot, error) {
		snap, ok := snapshots[env]
		if !ok {
			return nil, nil
		}
		return snap, nil
	}
}

func TestExecute_NoDrift_TextOutput(t *testing.T) {
	snaps := map[string]*state.Snapshot{
		"prod":    {Environment: "prod", Values: map[string]string{"key": "val"}},
		"staging": {Environment: "staging", Values: map[string]string{"key": "val"}},
	}
	cfg := makeTestConfig([]string{"prod", "staging"}, "text")
	engine := compare.NewEngineWithDeps(cfg, stubLoaderFn(snaps), drift.Detect)

	var buf bytes.Buffer
	runner := NewRunner(cfg, engine, &buf)

	if err := runner.Execute(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "no drift") {
		t.Errorf("expected 'no drift' in output, got: %s", buf.String())
	}
}

func TestExecute_WithDrift_JSONOutput(t *testing.T) {
	snaps := map[string]*state.Snapshot{
		"prod":    {Environment: "prod", Values: map[string]string{"key": "a"}},
		"staging": {Environment: "staging", Values: map[string]string{"key": "b"}},
	}
	cfg := makeTestConfig([]string{"prod", "staging"}, "json")
	engine := compare.NewEngineWithDeps(cfg, stubLoaderFn(snaps), drift.Detect)

	var buf bytes.Buffer
	runner := NewRunner(cfg, engine, &buf)

	if err := runner.Execute(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "\"diffs\"") {
		t.Errorf("expected JSON with diffs, got: %s", buf.String())
	}
}
