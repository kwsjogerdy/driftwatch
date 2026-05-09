package metrics

import "time"

// Summary holds aggregate statistics derived from a set of RunMetrics.
type Summary struct {
	TotalRuns      int
	ErrorRuns      int
	DriftRuns      int
	TotalDriftKeys int
	AvgDuration    time.Duration
}

func summarise(entries []RunMetrics) Summary {
	if len(entries) == 0 {
		return Summary{}
	}

	var s Summary
	var totalDur time.Duration

	for _, e := range entries {
		s.TotalRuns++
		totalDur += e.Duration
		s.TotalDriftKeys += e.DriftCount
		if e.HadError {
			s.ErrorRuns++
		}
		if e.DriftCount > 0 {
			s.DriftRuns++
		}
	}

	if s.TotalRuns > 0 {
		s.AvgDuration = totalDur / time.Duration(s.TotalRuns)
	}
	return s
}
