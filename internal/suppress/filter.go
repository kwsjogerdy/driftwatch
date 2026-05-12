package suppress

import (
	"time"

	"github.com/driftwatch/internal/diff"
)

// Filter removes diff entries that are covered by at least one active
// (non-expired) suppression rule.
type Filter struct {
	rules []Rule
	now   func() time.Time
}

// NewFilter constructs a Filter from the provided rules.
func NewFilter(rules []Rule) *Filter {
	return &Filter{
		rules: rules,
		now:   time.Now,
	}
}

// Apply returns a copy of diffs with suppressed entries removed.
// sourceEnv and targetEnv are the environment names being compared.
func (f *Filter) Apply(diffs []diff.Difference, sourceEnv, targetEnv string) []diff.Difference {
	now := f.now()
	var out []diff.Difference
	for _, d := range diffs {
		if !f.isSuppressed(d.Key, sourceEnv, targetEnv, now) {
			out = append(out, d)
		}
	}
	return out
}

// ActiveCount returns the number of rules that are not yet expired.
func (f *Filter) ActiveCount() int {
	now := f.now()
	count := 0
	for _, r := range f.rules {
		if !r.IsExpired(now) {
			count++
		}
	}
	return count
}

func (f *Filter) isSuppressed(key, sourceEnv, targetEnv string, now time.Time) bool {
	for _, r := range f.rules {
		if r.IsExpired(now) {
			continue
		}
		if r.Matches(key, sourceEnv, targetEnv) {
			return true
		}
	}
	return false
}
