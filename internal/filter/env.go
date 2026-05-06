package filter

import (
	"fmt"
	"strings"
)

// EnvFilter holds a set of allowed environment names for filtering comparisons.
type EnvFilter struct {
	allowed map[string]struct{}
}

// NewEnvFilter creates an EnvFilter from a slice of environment names.
// Names are normalized to lowercase. Returns an error if the slice is empty.
func NewEnvFilter(envs []string) (*EnvFilter, error) {
	if len(envs) == 0 {
		return nil, fmt.Errorf("filter: at least one environment name required")
	}
	allowed := make(map[string]struct{}, len(envs))
	for _, e := range envs {
		name := strings.TrimSpace(strings.ToLower(e))
		if name == "" {
			return nil, fmt.Errorf("filter: environment name must not be blank")
		}
		allowed[name] = struct{}{}
	}
	return &EnvFilter{allowed: allowed}, nil
}

// Allow reports whether the given environment name passes the filter.
func (f *EnvFilter) Allow(env string) bool {
	_, ok := f.allowed[strings.ToLower(env)]
	return ok
}

// FilterKeys returns only the keys from the provided map whose keys pass the filter.
func (f *EnvFilter) FilterKeys(envMap map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	for k, v := range envMap {
		if f.Allow(k) {
			out[k] = v
		}
	}
	return out
}

// Names returns the sorted list of allowed environment names.
func (f *EnvFilter) Names() []string {
	names := make([]string, 0, len(f.allowed))
	for k := range f.allowed {
		names = append(names, k)
	}
	return names
}
