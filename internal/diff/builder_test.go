package diff

import (
	"testing"
)

func TestBuild_NoDiff(t *testing.T) {
	b := NewBuilder(nil)
	src := map[string]interface{}{"key": "val"}
	tgt := map[string]interface{}{"key": "val"}
	diffs := b.Build(src, tgt)
	if len(diffs) != 0 {
		t.Fatalf("expected no diffs, got %d", len(diffs))
	}
}

func TestBuild_Mismatch_IsWarningByDefault(t *testing.T) {
	b := NewBuilder(nil)
	src := map[string]interface{}{"region": "us-east-1"}
	tgt := map[string]interface{}{"region": "eu-west-1"}
	diffs := b.Build(src, tgt)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Kind != "mismatch" {
		t.Errorf("expected kind 'mismatch', got %q", diffs[0].Kind)
	}
	if diffs[0].Severity != SeverityWarning {
		t.Errorf("expected warning severity, got %q", diffs[0].Severity)
	}
}

func TestBuild_CriticalKey_ElevatesSeverity(t *testing.T) {
	b := NewBuilder([]string{"db_password"})
	src := map[string]interface{}{"db_password": "secret1"}
	tgt := map[string]interface{}{"db_password": "secret2"}
	diffs := b.Build(src, tgt)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Severity != SeverityCritical {
		t.Errorf("expected critical severity, got %q", diffs[0].Severity)
	}
}

func TestBuild_MissingKeyInTarget(t *testing.T) {
	b := NewBuilder(nil)
	src := map[string]interface{}{"timeout": 30}
	tgt := map[string]interface{}{}
	diffs := b.Build(src, tgt)
	if len(diffs) != 1 || diffs[0].Kind != "missing" {
		t.Fatalf("expected one 'missing' diff, got %+v", diffs)
	}
}

func TestBuild_ExtraKeyInTarget(t *testing.T) {
	b := NewBuilder(nil)
	src := map[string]interface{}{}
	tgt := map[string]interface{}{"debug": true}
	diffs := b.Build(src, tgt)
	if len(diffs) != 1 || diffs[0].Kind != "extra" {
		t.Fatalf("expected one 'extra' diff, got %+v", diffs)
	}
	if diffs[0].Severity != SeverityInfo {
		t.Errorf("extra keys should be info severity, got %q", diffs[0].Severity)
	}
}

func TestApplyIgnore_RemovesKeys(t *testing.T) {
	m := map[string]interface{}{"a": 1, "b": 2, "c": 3}
	out := ApplyIgnore(m, []string{"b"})
	if _, ok := out["b"]; ok {
		t.Error("expected key 'b' to be removed")
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}
