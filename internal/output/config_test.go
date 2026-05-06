package output

import (
	"testing"
)

func TestParseDestination_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  Destination
	}{
		{"stdout", DestStdout},
		{"STDOUT", DestStdout},
		{"stderr", DestStderr},
		{"file", DestFile},
		{" file ", DestFile},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ParseDestination(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestParseDestination_Invalid(t *testing.T) {
	_, err := ParseDestination("syslog")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestValidateConfig_Valid(t *testing.T) {
	cases := []Config{
		{Destination: DestStdout, Format: "text"},
		{Destination: DestStderr, Format: "json"},
		{Destination: DestFile, FilePath: "/tmp/out.txt", Format: "json"},
	}
	for _, c := range cases {
		if err := ValidateConfig(c); err != nil {
			t.Errorf("unexpected error for %+v: %v", c, err)
		}
	}
}

func TestValidateConfig_FileMissingPath(t *testing.T) {
	c := Config{Destination: DestFile, Format: "text"}
	if err := ValidateConfig(c); err == nil {
		t.Fatal("expected error for file destination without path")
	}
}

func TestValidateConfig_EmptyFormat(t *testing.T) {
	c := Config{Destination: DestStdout, Format: ""}
	if err := ValidateConfig(c); err == nil {
		t.Fatal("expected error for empty format")
	}
}

func TestValidateConfig_InvalidFormat(t *testing.T) {
	c := Config{Destination: DestStdout, Format: "xml"}
	if err := ValidateConfig(c); err == nil {
		t.Fatal("expected error for unsupported format")
	}
}
