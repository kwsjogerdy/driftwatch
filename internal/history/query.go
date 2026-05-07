package history

import (
	"sort"
	"time"
)

// QueryOptions controls filtering when retrieving history entries.
type QueryOptions struct {
	Since    time.Time
	Until    time.Time
	EnvPair  string // e.g. "staging->production"
	MaxItems int
}

// Query returns entries from the store filtered by the provided options.
func Query(entries []Entry, opts QueryOptions) []Entry {
	var result []Entry

	for _, e := range entries {
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		if !opts.Until.IsZero() && e.Timestamp.After(opts.Until) {
			continue
		}
		if opts.EnvPair != "" && e.EnvPair != opts.EnvPair {
			continue
		}
		result = append(result, e)
	}

	// Sort descending by timestamp (most recent first).
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.After(result[j].Timestamp)
	})

	if opts.MaxItems > 0 && len(result) > opts.MaxItems {
		result = result[:opts.MaxItems]
	}

	return result
}
