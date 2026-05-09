package metrics

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func makeExporter(t *testing.T) (*Exporter, *Collector) {
	t.Helper()
	c := NewCollector()
	return NewExporter(c), c
}

func TestExport_TextFormat_NoRecords(t *testing.T) {
	ex, _ := makeExporter(t)
	var buf bytes.Buffer
	if err := ex.Export(&buf, ExportOptions{Format: ExportText}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no metrics recorded") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestExport_TextFormat_WithRecords(t *testing.T) {
	ex, c := makeExporter(t)
	c.Record(Record{EnvPair: "prod:staging", Drifted: true, DiffCount: 3})
	var buf bytes.Buffer
	if err := ex.Export(&buf, ExportOptions{Format: ExportText}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "prod:staging") {
		t.Errorf("expected env pair in output, got: %s", out)
	}
	if !strings.Contains(out, "diffs=3") {
		t.Errorf("expected diff count in output, got: %s", out)
	}
}

func TestExport_JSONFormat_ValidStructure(t *testing.T) {
	ex, c := makeExporter(t)
	c.Record(Record{EnvPair: "dev:prod", Drifted: false, DiffCount: 0})
	var buf bytes.Buffer
	if err := ex.Export(&buf, ExportOptions{Format: ExportJSON}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var records []Record
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 1 || records[0].EnvPair != "dev:prod" {
		t.Errorf("unexpected records: %+v", records)
	}
}

func TestExport_FilterBySince(t *testing.T) {
	ex, c := makeExporter(t)
	old := Record{EnvPair: "a:b", Drifted: true, DiffCount: 1, Timestamp: time.Now().Add(-2 * time.Hour)}
	recent := Record{EnvPair: "a:b", Drifted: false, DiffCount: 0, Timestamp: time.Now()}
	c.Record(old)
	c.Record(recent)
	var buf bytes.Buffer
	opts := ExportOptions{Format: ExportText, Since: time.Now().Add(-30 * time.Minute)}
	if err := ex.Export(&buf, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 record after since filter, got %d: %s", len(lines), buf.String())
	}
}

func TestExport_FilterByEnvPair(t *testing.T) {
	ex, c := makeExporter(t)
	c.Record(Record{EnvPair: "prod:staging", Drifted: true, DiffCount: 2})
	c.Record(Record{EnvPair: "dev:staging", Drifted: false, DiffCount: 0})
	var buf bytes.Buffer
	opts := ExportOptions{Format: ExportText, EnvFilter: "prod:staging"}
	if err := ex.Export(&buf, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(buf.String(), "dev:staging") {
		t.Errorf("expected dev:staging to be filtered out, got: %s", buf.String())
	}
}

func TestExport_UnknownFormat_ReturnsError(t *testing.T) {
	ex, _ := makeExporter(t)
	var buf bytes.Buffer
	err := ex.Export(&buf, ExportOptions{Format: "xml"})
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}
