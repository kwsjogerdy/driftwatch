package audit

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func makeEvent(level Level, action, envPair string, driftCount int) Event {
	return Event{
		Timestamp:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Level:      level,
		Action:     action,
		EnvPair:    envPair,
		Message:    "test message",
		DriftCount: driftCount,
	}
}

func TestWrite_TextFormat_ContainsFields(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf, "text")
	e := makeEvent(LevelInfo, "compare", "prod->staging", 3)
	if err := l.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"info", "compare", "prod->staging", "drift_count=3"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output: %s", want, out)
		}
	}
}

func TestWrite_JSONFormat_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf, "json")
	e := makeEvent(LevelWarn, "drift_detected", "dev->prod", 5)
	if err := l.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded Event
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, buf.String())
	}
	if decoded.Action != "drift_detected" {
		t.Errorf("expected action drift_detected, got %s", decoded.Action)
	}
	if decoded.DriftCount != 5 {
		t.Errorf("expected drift_count 5, got %d", decoded.DriftCount)
	}
}

func TestWrite_SetsTimestampIfZero(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf, "json")
	e := Event{Level: LevelInfo, Action: "check"}
	if err := l.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded Event
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if decoded.Timestamp.IsZero() {
		t.Error("expected timestamp to be set automatically")
	}
}

func TestWrite_TextFormat_NoEnvPair_OmitsField(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf, "text")
	e := Event{Level: LevelError, Action: "load_failed", Message: "file not found"}
	if err := l.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "env_pair") {
		t.Errorf("did not expect env_pair in output: %s", out)
	}
}

func TestNewLogger_DefaultsToTextAndStdout(t *testing.T) {
	l := NewLogger(nil, "")
	if l.format != "text" {
		t.Errorf("expected default format text, got %s", l.format)
	}
	if l.out == nil {
		t.Error("expected non-nil writer")
	}
}
