package policy

import (
	"errors"
	"time"
)

// Severity represents the alert level for a policy rule.
type Severity string

const (
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

// Rule defines a policy constraint applied to a specific key or key prefix.
type Rule struct {
	ID        string    `json:"id"`
	KeyPrefix string    `json:"key_prefix"`
	Severity  Severity  `json:"severity"`
	Message   string    `json:"message,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// IsExpired reports whether the rule has passed its expiry time.
func (r Rule) IsExpired(now time.Time) bool {
	if r.ExpiresAt.IsZero() {
		return false
	}
	return now.After(r.ExpiresAt)
}

// Validate returns an error if the rule is not well-formed.
func (r Rule) Validate() error {
	if r.ID == "" {
		return errors.New("rule id must not be empty")
	}
	if r.KeyPrefix == "" {
		return errors.New("rule key_prefix must not be empty")
	}
	if r.Severity != SeverityWarning && r.Severity != SeverityCritical {
		return errors.New("rule severity must be 'warning' or 'critical'")
	}
	return nil
}
