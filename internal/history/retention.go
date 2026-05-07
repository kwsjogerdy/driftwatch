package history

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// RetentionPolicy defines how long or how many history entries to keep.
type RetentionPolicy struct {
	// MaxAge removes entries older than this duration. Zero means no age limit.
	MaxAge time.Duration
	// MaxEntries keeps only the N most recent entries. Zero means no limit.
	MaxEntries int
}

// DefaultRetentionPolicy returns a sensible default policy.
func DefaultRetentionPolicy() RetentionPolicy {
	return RetentionPolicy{
		MaxAge:     30 * 24 * time.Hour, // 30 days
		MaxEntries: 100,
	}
}

// Apply enforces the retention policy against the given store directory,
// deleting files that exceed age or count limits. Returns the number of
// files removed.
func Apply(dir string, policy RetentionPolicy) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("read history dir: %w", err)
	}

	type fileInfo struct {
		name    string
		modTime time.Time
	}

	var files []fileInfo
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		files = append(files, fileInfo{name: e.Name(), modTime: info.ModTime()})
	}

	// Sort newest first.
	sort.Slice(files, func(i, j int) bool {
		return files[i].modTime.After(files[j].modTime)
	})

	now := time.Now()
	removed := 0

	for i, f := range files {
		shouldRemove := false

		if policy.MaxAge > 0 && now.Sub(f.modTime) > policy.MaxAge {
			shouldRemove = true
		}
		if policy.MaxEntries > 0 && i >= policy.MaxEntries {
			shouldRemove = true
		}

		if shouldRemove {
			path := filepath.Join(dir, f.name)
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				return removed, fmt.Errorf("remove %s: %w", path, err)
			}
			removed++
		}
	}

	return removed, nil
}
