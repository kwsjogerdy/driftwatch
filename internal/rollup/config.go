package rollup

import (
	"fmt"
	"time"
)

// Config controls how drift history is rolled up into aggregate reports.
type Config struct {
	// Period is the aggregation window: "daily", "weekly", or "monthly".
	Period string `json:"period"`

	// TopKeysLimit is the maximum number of frequently-drifting keys to surface.
	TopKeysLimit int `json:"top_keys_limit"`
}

var validPeriods = map[string]bool{
	"daily":   true,
	"weekly":  true,
	"monthly": true,
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Period:       "daily",
		TopKeysLimit: 10,
	}
}

// Validate returns an error if the Config contains invalid values.
func Validate(cfg Config) error {
	if !validPeriods[cfg.Period] {
		return fmt.Errorf("rollup: invalid period %q: must be daily, weekly, or monthly", cfg.Period)
	}
	if cfg.TopKeysLimit <= 0 {
		return fmt.Errorf("rollup: top_keys_limit must be positive, got %d", cfg.TopKeysLimit)
	}
	return nil
}

// PeriodFor returns the start and end times of the aggregation window that
// contains t, according to cfg.Period.
func PeriodFor(cfg Config, t time.Time) (start, end time.Time) {
	switch cfg.Period {
	case "weekly":
		// Align to Monday of the current week.
		weekday := int(t.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday → 7 so Monday is day 1
		}
		start = time.Date(t.Year(), t.Month(), t.Day()-weekday+1, 0, 0, 0, 0, t.Location())
		end = start.Add(7 * 24 * time.Hour)
	case "monthly":
		start = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
		end = start.AddDate(0, 1, 0)
	default: // "daily"
		start = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		end = start.Add(24 * time.Hour)
	}
	return start, end
}
