package main

import (
	"context"
	"fmt"
	"io"

	"github.com/example/driftwatch/internal/compare"
	"github.com/example/driftwatch/internal/config"
	"github.com/example/driftwatch/internal/report"
)

// Runner executes a single drift comparison cycle and writes the result.
type Runner struct {
	cfg    *config.AppConfig
	engine *compare.Engine
	out    io.Writer
}

// NewRunner creates a Runner bound to the given config, engine, and output.
func NewRunner(cfg *config.AppConfig, engine *compare.Engine, out io.Writer) *Runner {
	return &Runner{cfg: cfg, engine: engine, out: out}
}

// Execute runs one comparison cycle, formats the report, and writes output.
func (r *Runner) Execute(ctx context.Context) error {
	diffs, err := r.engine.Compare(ctx)
	if err != nil {
		return fmt.Errorf("compare failed: %w", err)
	}

	rep := report.NewReport(
		r.cfg.Schedule.Environments[0],
		r.cfg.Schedule.Environments[1],
		diffs,
	)

	output, err := rep.Render(r.cfg.Output.Format)
	if err != nil {
		return fmt.Errorf("render failed: %w", err)
	}

	_, err = fmt.Fprintln(r.out, output)
	return err
}
