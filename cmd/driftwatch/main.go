package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/example/driftwatch/internal/config"
	"github.com/example/driftwatch/internal/compare"
	"github.com/example/driftwatch/internal/output"
	"github.com/example/driftwatch/internal/schedule"
)

func main() {
	configPath := flag.String("config", "driftwatch.json", "path to config file")
	flag.Parse()

	cfg, err := config.LoadAppConfig(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	config.ApplyDefaults(cfg)

	if err := config.Validate(cfg); err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	writer, err := output.NewWriter(cfg.Output)
	if err != nil {
		log.Fatalf("failed to create output writer: %v", err)
	}
	defer writer.Close()

	engine := compare.NewEngine(cfg)

	watcher, err := schedule.NewWatcher(cfg.Schedule, engine, writer)
	if err != nil {
		log.Fatalf("failed to create watcher: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	fmt.Fprintf(os.Stderr, "driftwatch started (interval: %s)\n", cfg.Schedule.Interval)

	if err := watcher.Run(ctx); err != nil {
		log.Fatalf("watcher error: %v", err)
	}

	fmt.Fprintln(os.Stderr, "driftwatch stopped")
}
