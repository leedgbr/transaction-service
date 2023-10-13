package date

import "time"

// NewInUTC is a convenience function for creating a date for testing purposes.
func NewInUTC(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
