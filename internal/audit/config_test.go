package audit

import (
	"testing"
)

func TestDefaultConfig_Values(t *testing.T) {
	c := DefaultConfig()
	if !c.Enabled {
		t.Error("expected Enabled=true by default")
	}
	if c.Format != "text" {
		t.Errorf("expected Format=text, got %s", c.Format)
	}
	if c.Destination != "stdout" {
		t.Errorf("expected Destination=stdout, got %s", c.Destination)
	}
}

func TestValidate_Valid(t *testing.T) {
	cases := []Config{
		{Enabled: true, Format: "text", Destination: "stdout"},
		{Enabled: true, Format: "json", Destination: "stderr"},
		{Enabled: true, Format: "JSON", Destination: "/var/log/drift.log"},
		{Enabled: false, Format: "", Destination: ""},
	}
	for _, c := range cases {
		if err := Validate(c); err != nil {
			t.Errorf("expected no error for %+v, got %v", c, err)
		}
	}
}

func TestValidate_InvalidFormat(t *testing.T) {
	c := Config{Enabled: true, Format: "xml", Destination: "stdout"}
	if err := Validate(c); err == nil {
		t.Error("expected error for invalid format")
	}
}

func TestValidate_EmptyDestination(t *testing.T) {
	c := Config{Enabled: true, Format: "text", Destination: ""}
	if err := Validate(c); err == nil {
		t.Error("expected error for empty destination")
	}
}

func TestValidate_DisabledSkipsChecks(t *testing.T) {
	c := Config{Enabled: false, Format: "bad", Destination: ""}
	if err := Validate(c); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestNormaliseFormat_Defaults(t *testing.T) {
	if got := NormaliseFormat(Config{Format: ""}); got != "text" {
		t.Errorf("expected text, got %s", got)
	}
	if got := NormaliseFormat(Config{Format: "JSON"}); got != "json" {
		t.Errorf("expected json, got %s", got)
	}
}
