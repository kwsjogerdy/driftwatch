package drift_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/state"
)

func makeSnapshot(env string, resources map[string]string) *state.Snapshot {
	return &state.Snapshot{
		Environment: env,
		Resources:   resources,
	}
}

func TestDetect_NoDrift(t *testing.T) {
	baseline := makeSnapshot("staging", map[string]string{"aws.s3.bucket": "my-bucket", "aws.ec2.ami": "ami-123"})
	target := makeSnapshot("production", map[string]string{"aws.s3.bucket": "my-bucket", "aws.ec2.ami": "ami-123"})

	result, err := drift.Detect(baseline, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Drifted {
		t.Errorf("expected no drift, got %d mismatches", len(result.Mismatches))
	}
}

func TestDetect_ValueMismatch(t *testing.T) {
	baseline := makeSnapshot("staging", map[string]string{"aws.ec2.ami": "ami-old"})
	target := makeSnapshot("production", map[string]string{"aws.ec2.ami": "ami-new"})

	result, err := drift.Detect(baseline, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Drifted {
		t.Fatal("expected drift to be detected")
	}
	if len(result.Mismatches) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(result.Mismatches))
	}
	if result.Mismatches[0].Key != "aws.ec2.ami" {
		t.Errorf("unexpected mismatch key: %s", result.Mismatches[0].Key)
	}
}

func TestDetect_MissingKeyInTarget(t *testing.T) {
	baseline := makeSnapshot("staging", map[string]string{"aws.rds.instance": "db-prod"})
	target := makeSnapshot("production", map[string]string{})

	result, err := drift.Detect(baseline, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Drifted {
		t.Fatal("expected drift due to missing key in target")
	}
	if result.Mismatches[0].TargetValue != "<missing>" {
		t.Errorf("expected <missing>, got %s", result.Mismatches[0].TargetValue)
	}
}

func TestDetect_ExtraKeyInTarget(t *testing.T) {
	baseline := makeSnapshot("staging", map[string]string{})
	target := makeSnapshot("production", map[string]string{"aws.lambda.arn": "arn:aws:lambda:us-east-1:fn"})

	result, err := drift.Detect(baseline, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Drifted {
		t.Fatal("expected drift due to extra key in target")
	}
	if result.Mismatches[0].BaselineValue != "<missing>" {
		t.Errorf("expected <missing>, got %s", result.Mismatches[0].BaselineValue)
	}
}

func TestDetect_NilBaseline(t *testing.T) {
	_, err := drift.Detect(nil, makeSnapshot("production", map[string]string{}))
	if err == nil {
		t.Fatal("expected error for nil baseline")
	}
}

func TestDetect_NilTarget(t *testing.T) {
	_, err := drift.Detect(makeSnapshot("staging", map[string]string{}), nil)
	if err == nil {
		t.Fatal("expected error for nil target")
	}
}
