package suppress

import (
	"testing"
	"time"
)

func baseRule() Rule {
	return Rule{
		ID:        "r1",
		KeyPrefix: "db.",
		SourceEnv: "staging",
		TargetEnv: "production",
		ExpiresAt: time.Now().Add(time.Hour),
		Reason:    "planned migration",
	}
}

func TestValidate_Valid(t *testing.T) {
	if err := baseRule().Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidate_MissingID(t *testing.T) {
	r := baseRule()
	r.ID = ""
	if err := r.Validate(); err == nil {
		t.Fatal("expected error for missing id")
	}
}

func TestValidate_MissingKeyPrefix(t *testing.T) {
	r := baseRule()
	r.KeyPrefix = ""
	if err := r.Validate(); err == nil {
		t.Fatal("expected error for empty key_prefix")
	}
}

func TestValidate_ZeroExpiry(t *testing.T) {
	r := baseRule()
	r.ExpiresAt = time.Time{}
	if err := r.Validate(); err == nil {
		t.Fatal("expected error for zero expires_at")
	}
}

func TestIsExpired(t *testing.T) {
	r := baseRule()
	if r.IsExpired(time.Now()) {
		t.Fatal("rule should not be expired")
	}
	if !r.IsExpired(r.ExpiresAt.Add(time.Second)) {
		t.Fatal("rule should be expired after ExpiresAt")
	}
}

func TestMatches_ExactEnvs(t *testing.T) {
	r := baseRule()
	if !r.Matches("db.host", "staging", "production") {
		t.Error("expected match")
	}
	if r.Matches("db.host", "dev", "production") {
		t.Error("unexpected match on wrong source env")
	}
}

func TestMatches_WildcardEnvs(t *testing.T) {
	r := baseRule()
	r.SourceEnv = ""
	r.TargetEnv = ""
	if !r.Matches("db.port", "any", "other") {
		t.Error("expected wildcard match")
	}
}

func TestMatches_PrefixMismatch(t *testing.T) {
	r := baseRule()
	if r.Matches("cache.host", "staging", "production") {
		t.Error("unexpected match on wrong key prefix")
	}
}
