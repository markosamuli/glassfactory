package reporting

import (
	"fmt"
	"time"
)

// CalendarMonth represents a calendar month
type CalendarMonth struct {
	Year  int        // Year (e.g., 2014).
	Month time.Month // Month of the year (January = 1, ...).
}

// Before reports whether m occurs before m2.
func (m CalendarMonth) Before(m2 CalendarMonth) bool {
	if m.Year != m2.Year {
		return m.Year < m2.Year
	}
	return m.Month < m2.Month
}

// After reports whether m occurs after m2.
func (m CalendarMonth) After(m2 CalendarMonth) bool {
	return m2.Before(m)
}

// String returns the month in YYYY-MM format.
func (m CalendarMonth) String() string {
	return fmt.Sprintf("%04d-%02d", m.Year, m.Month)
}