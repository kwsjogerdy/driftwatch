package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of an audit event.
type Level string

const (
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Event captures a single audit log entry.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Level     Level     `json:"level"`
	Action    string    `json:"action"`
	EnvPair   string    `json:"env_pair,omitempty"`
	Message   string    `json:"message"`
	DriftCount int      `json:"drift_count,omitempty"`
}

// Logger writes structured audit events to a destination.
type Logger struct {
	out    io.Writer
	format string
}

// NewLogger creates a Logger writing to out in the given format ("text" or "json").
func NewLogger(out io.Writer, format string) *Logger {
	if out == nil {
		out = os.Stdout
	}
	if format == "" {
		format = "text"
	}
	return &Logger{out: out, format: format}
}

// Write emits an audit event.
func (l *Logger) Write(e Event) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	switch l.format {
	case "json":
		return l.writeJSON(e)
	default:
		return l.writeText(e)
	}
}

func (l *Logger) writeText(e Event) error {
	line := fmt.Sprintf("%s [%s] action=%s", e.Timestamp.Format(time.RFC3339), e.Level, e.Action)
	if e.EnvPair != "" {
		line += fmt.Sprintf(" env_pair=%s", e.EnvPair)
	}
	if e.DriftCount > 0 {
		line += fmt.Sprintf(" drift_count=%d", e.DriftCount)
	}
	if e.Message != "" {
		line += fmt.Sprintf(" msg=%q", e.Message)
	}
	_, err := fmt.Fprintln(l.out, line)
	return err
}

func (l *Logger) writeJSON(e Event) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(l.out, string(b))
	return err
}
