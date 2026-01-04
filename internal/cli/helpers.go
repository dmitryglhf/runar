package cli

import (
	"time"
)

func formatDuration(start time.Time, end *time.Time) string {
	if end == nil {
		return time.Since(start).Round(time.Second).String()
	}
	return end.Sub(start).Round(time.Second).String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
