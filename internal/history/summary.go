package history

import "time"

// Summary aggregates statistics over a slice of history entries.
type Summary struct {
	TotalRuns    int
	DriftingRuns int
	CleanRuns    int
	LastChecked  time.Time
	LastDrift    time.Time
}

// Summarise computes a Summary from the provided entries.
func Summarise(entries []Entry) Summary {
	var s Summary
	s.TotalRuns = len(entries)
	for _, e := range entries {
		if e.DriftCount > 0 {
			s.DriftingRuns++
			if e.Timestamp.After(s.LastDrift) {
				s.LastDrift = e.Timestamp
			}
		} else {
			s.CleanRuns++
		}
		if e.Timestamp.After(s.LastChecked) {
			s.LastChecked = e.Timestamp
		}
	}
	return s
}
