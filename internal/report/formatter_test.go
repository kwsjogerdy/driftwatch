package report_test

import (
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/report"
)

func makeDiffs() []drift.Diff {
	return []drift.Diff{
		{Key: "db_host", Type: "mismatch", BaselineValue: "prod-db", TargetValue: "staging-db"},
		{Key: "replica_count", Type: "missing", BaselineValue: "3", TargetValue: ""},
	}
}

func TestNewReport_StoresFields(t *testing.T) {
	diffs := makeDiffs()
	r := report.NewReport("staging", diffs, report.FormatText)
	if r.Environment != "staging" {
		t.Errorf("expected environment staging, got %s", r.Environment)
	}
	if len(r.Diffs) != 2 {
		t.Errorf("expected 2 diffs, got %d", len(r.Diffs))
	}
}

func TestRender_TextFormat_NoDrift(t *testing.T) {
	r := report.NewReport("prod", []drift.Diff{}, report.FormatText)
	out := r.Render()
	if !strings.Contains(out, "No drift detected") {
		t.Errorf("expected 'No drift detected' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Environment: prod") {
		t.Errorf("expected environment name in output")
	}
}

func TestRender_TextFormat_WithDrift(t *testing.T) {
	r := report.NewReport("staging", makeDiffs(), report.FormatText)
	out := r.Render()
	if !strings.Contains(out, "Total diffs: 2") {
		t.Errorf("expected diff count in output, got:\n%s", out)
	}
	if !strings.Contains(out, "db_host") {
		t.Errorf("expected key db_host in output")
	}
	if !strings.Contains(out, "MISMATCH") {
		t.Errorf("expected MISMATCH type label in output")
	}
}

func TestRender_JSONFormat_WithDrift(t *testing.T) {
	r := report.NewReport("staging", makeDiffs(), report.FormatJSON)
	out := r.Render()
	if !strings.Contains(out, `"environment":"staging"`) {
		t.Errorf("expected environment field in JSON, got:\n%s", out)
	}
	if !strings.Contains(out, `"total_diffs":2`) {
		t.Errorf("expected total_diffs in JSON")
	}
	if !strings.Contains(out, `"key":"db_host"`) {
		t.Errorf("expected key field in JSON diffs")
	}
}

func TestRender_JSONFormat_NoDrift(t *testing.T) {
	r := report.NewReport("prod", []drift.Diff{}, report.FormatJSON)
	out := r.Render()
	if !strings.Contains(out, `"total_diffs":0`) {
		t.Errorf("expected zero diffs in JSON, got:\n%s", out)
	}
	if !strings.Contains(out, `"diffs":[]`) {
		t.Errorf("expected empty diffs array in JSON")
	}
}
