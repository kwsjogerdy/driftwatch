package metrics

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func makeSummary() Summary {
	return Summary{
		TotalRuns:      5,
		DriftRuns:      2,
		ErrorRuns:      1,
		TotalDriftKeys: 7,
		AvgDuration:    250 * time.Millisecond,
	}
}

func TestRender_TextFormat_ContainsFields(t *testing.T) {
	var buf bytes.Buffer
	if err := Render(&buf, makeSummary(), FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"total_runs", "drift_runs", "error_runs", "avg_duration_ms"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing field %q", want)
		}
	}
}

func TestRender_JSONFormat_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := Render(&buf, makeSummary(), FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if v, ok := m["total_runs"]; !ok || v == nil {
		t.Error("JSON missing total_runs")
	}
}

func TestRender_UnknownFormat_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	if err := Render(&buf, makeSummary(), Format("xml")); err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}

func TestRender_TextFormat_CorrectValues(t *testing.T) {
	s := Summary{TotalRuns: 3, DriftRuns: 1, AvgDuration: 100 * time.Millisecond}
	var buf bytes.Buffer
	if err := Render(&buf, s, FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "total_runs:       3") {
		t.Errorf("unexpected output:\n%s", out)
	}
}
