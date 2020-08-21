package dateutil

import (
	"time"

	"github.com/jinzhu/now"
)

// DateRange represents a range between two dates
type DateRange struct {
	Start time.Time
	End   time.Time
}

// MonthsBetweenDates returns full calendar months matching the given start and end dates
func MonthsBetweenDates(start time.Time, end time.Time) []DateRange {
	var months []DateRange
	start = now.With(start).BeginningOfMonth()
	endOfLastMonth := end
	for ok := true; ok; ok = start.Before(endOfLastMonth) {
		end := now.With(start).EndOfMonth()
		months = append(months, DateRange{
			Start: start,
			End:   end,
		})
		start = now.With(start.AddDate(0, 1, 0)).BeginningOfMonth()
	}
	return months
}
