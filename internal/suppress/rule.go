package suppress

import (
	"errors"
	"strings"
	"time"
)

// Rule defines a suppression rule that silences drift alerts for a specific
// key pattern across an optional environment pair for a bounded duration.
type Rule struct {
	ID        string    `json:"id"`
	KeyPrefix string    `json:"key_prefix"`
	SourceEnv string    `json:"source_env"`
	TargetEnv string    `json:"target_env"`
	ExpiresAt time.Time `json:"expires_at"`
	Reason    string    `json:"reason"`
}

// Validate returns an error if the rule is not well-formed.
func (r Rule) Validate() error {
	if strings.TrimSpace(r.ID) == "" {
		return errors.New("rule id must not be empty")
	}
	if strings.TrimSpace(r.KeyPrefix) == "" {
		return errors.New("key_prefix must not be empty")
	}
	if r.ExpiresAt.IsZero() {
		return errors.New("expires_at must be set")
	}
	return nil
}

// IsExpired reports whether the rule has passed its expiry time.
func (r Rule) IsExpired(now time.Time) bool {
	return now.After(r.ExpiresAt)
}

// Matches reports whether the rule suppresses the given key for the given
// environment pair. An empty SourceEnv or TargetEnv in the rule matches any
// value.
func (r Rule) Matches(key, sourceEnv, targetEnv string) bool {
	if !strings.HasPrefix(key, r.KeyPrefix) {
		return false
	}
	if r.SourceEnv != "" && !strings.EqualFold(r.SourceEnv, sourceEnv) {
		return false
	}
	if r.TargetEnv != "" && !strings.EqualFold(r.TargetEnv, targetEnv) {
		return false
	}
	return true
}
