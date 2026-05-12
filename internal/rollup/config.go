package rollup

import (
	"errors"
	"time"
)

// Granularity describes the time bucket size used when building periods.
type Granularity string

const (
	GranularityHourly  Granularity = "hourly"
	GranularityDaily   Granularity = "daily"
	GranularityWeekly  Granularity = "weekly"
)

// Config controls rollup aggregation behaviour.
type Config struct {
	// Granularity sets the aggregation window size.
	Granularity Granularity `json:"granularity"`
	// MaxTopKeys limits how many frequently-drifted keys are reported.
	MaxTopKeys int `json:"max_top_keys"`
	// RetentionDays controls how far back entries are considered.
	RetentionDays int `json:"retention_days"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Granularity:   GranularityDaily,
		MaxTopKeys:    5,
		RetentionDays: 30,
	}
}

// Validate checks that the Config fields are acceptable.
func Validate(c Config) error {
	switch c.Granularity {
	case GranularityHourly, GranularityDaily, GranularityWeekly:
		// valid
	default:
		return errors.New("rollup: granularity must be one of hourly, daily, weekly")
	}
	if c.MaxTopKeys < 1 {
		return errors.New("rollup: max_top_keys must be at least 1")
	}
	if c.RetentionDays < 1 {
		return errors.New("rollup: retention_days must be at least 1")
	}
	return nil
}

// PeriodFor returns a Period aligned to the given granularity that contains t.
func PeriodFor(t time.Time, g Granularity) Period {
	switch g {
	case GranularityHourly:
		start := t.Truncate(time.Hour)
		return Period{Start: start, End: start.Add(time.Hour)}
	case GranularityWeekly:
		// Align to Monday.
		weekday := int(t.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start := time.Date(t.Year(), t.Month(), t.Day()-weekday+1, 0, 0, 0, 0, t.Location())
		return Period{Start: start, End: start.Add(7 * 24 * time.Hour)}
	default: // daily
		start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		return Period{Start: start, End: start.Add(24 * time.Hour)}
	}
}
