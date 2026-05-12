package rollup

import (
	"fmt"
	"sort"
	"time"

	"github.com/driftwatch/internal/diff"
)

// Period defines a time window for aggregation.
type Period struct {
	Start time.Time
	End   time.Time
}

// EnvPairKey identifies a source/target environment pair.
type EnvPairKey struct {
	Source string
	Target string
}

// AggregatedResult holds rolled-up drift stats for a single env pair within a period.
type AggregatedResult struct {
	EnvPair      EnvPairKey
	Period       Period
	TotalRuns    int
	DriftedRuns  int
	CleanRuns    int
	DriftRate    float64
	TopDriftKeys []string
}

// Entry is a minimal drift record consumed by the aggregator.
type Entry struct {
	Timestamp time.Time
	Source    string
	Target    string
	Diffs     []diff.Difference
}

// Aggregate groups entries by env pair and computes rollup statistics
// for those whose timestamps fall within the given period.
func Aggregate(entries []Entry, period Period) []AggregatedResult {
	groups := make(map[EnvPairKey][]Entry)
	for _, e := range entries {
		if e.Timestamp.Before(period.Start) || e.Timestamp.After(period.End) {
			continue
		}
		key := EnvPairKey{Source: e.Source, Target: e.Target}
		groups[key] = append(groups[key], e)
	}

	results := make([]AggregatedResult, 0, len(groups))
	for key, runs := range groups {
		results = append(results, computeResult(key, runs, period))
	}
	sort.Slice(results, func(i, j int) bool {
		return fmt.Sprintf("%s/%s", results[i].EnvPair.Source, results[i].EnvPair.Target) <
			fmt.Sprintf("%s/%s", results[j].EnvPair.Source, results[j].EnvPair.Target)
	})
	return results
}

func computeResult(key EnvPairKey, runs []Entry, period Period) AggregatedResult {
	keyCounts := make(map[string]int)
	drifted := 0
	for _, run := range runs {
		if len(run.Diffs) > 0 {
			drifted++
		}
		for _, d := range run.Diffs {
			keyCounts[d.Key]++
		}
	}
	total := len(runs)
	rate := 0.0
	if total > 0 {
		rate = float64(drifted) / float64(total)
	}
	return AggregatedResult{
		EnvPair:      key,
		Period:       period,
		TotalRuns:    total,
		DriftedRuns:  drifted,
		CleanRuns:    total - drifted,
		DriftRate:    rate,
		TopDriftKeys: topKeys(keyCounts, 5),
	}
}

func topKeys(counts map[string]int, n int) []string {
	type kv struct {
		key   string
		count int
	}
	pairs := make([]kv, 0, len(counts))
	for k, c := range counts {
		pairs = append(pairs, kv{k, c})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count != pairs[j].count {
			return pairs[i].count > pairs[j].count
		}
		return pairs[i].key < pairs[j].key
	})
	keys := make([]string, 0, n)
	for i, p := range pairs {
		if i >= n {
			break
		}
		keys = append(keys, p.key)
	}
	return keys
}
