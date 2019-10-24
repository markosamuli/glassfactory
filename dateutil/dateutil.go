// Package dateutil provides tools and types for handling date and time information
package dateutil

import (
	"fmt"
	"time"
)

type DateRange struct {
	Start time.Time
	End time.Time
}

type CalendarMonth struct {
	Year  int        // Year (e.g., 2014).
	Month time.Month // Month of the year (January = 1, ...).
}

func CalendarMonthOf(d time.Time) CalendarMonth {
	return CalendarMonth{
		Year:  d.Year(),
		Month: d.Month(),
	}
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

// String returns the date in RFC3339 full-date format.
func (m CalendarMonth) String() string {
	return fmt.Sprintf("%04d-%02d", m.Year, m.Month)
}

type FiscalYear struct {
	Start time.Time
	End time.Time
}

func (fy FiscalYear) String() string {
	return fmt.Sprintf("FY %04d", fy.End.Year())
}

// Before reports whether fy occurs before fy2.
func (fy FiscalYear) Before(fy2 FiscalYear) bool {
	return fy.End.Before(fy2.End)
}

// After reports whether fy occurs after fy2.
func (fy FiscalYear) After(fy2 FiscalYear) bool {
	return fy2.Before(fy)
}

func NewFiscalYear(d time.Time, finalMonth time.Month) *FiscalYear {
	var start time.Time
	var end time.Time
	if finalMonth < time.December {
		start = time.Date(d.Year() - 1, finalMonth + 1, 1, 0, 0, 0, 0, time.Local)
		end = time.Date(d.Year(), finalMonth, 1, 0, 0, 0, 0, time.Local)
	} else {
		start = time.Date(d.Year(), time.January, 1, 0, 0, 0, 0, time.Local)
		end = time.Date(d.Year(), time.December, 1, 0, 0, 0, 0, time.Local)
	}
	if d.Before(start) {
		start = start.AddDate(-1, 0 ,0)
		end = end.AddDate(-1,0 ,0)
	} else if d.After(end) {
		start = start.AddDate(1, 0 ,0)
		end = end.AddDate(1, 0 ,0)
	}
	return &FiscalYear{
		Start: start,
		End: end,
	}
}

func MonthsBetweenDates(start time.Time, end time.Time) []DateRange {
	var months []DateRange
	start = BeginningOfMonth(start)
	endOfLastMonth := end
	for ok := true; ok; ok = start.Before(endOfLastMonth) {
		end := EndOfMonth(start)
		months = append(months, DateRange{
			Start: start,
			End: end,
		})
		start = BeginningOfMonth(start.AddDate(0, 1, 0))
	}
	return months
}

func BeginningOfMonth(t time.Time) time.Time {
	y, m, _ := t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
}

func EndOfMonth(t time.Time) time.Time {
	return BeginningOfMonth(t).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

func BeginningOfYear(t time.Time) time.Time {
	y, _, _ := t.Date()
	return time.Date(y, time.January, 1, 0, 0, 0, 0, t.Location())
}

func EndOfYear(t time.Time) time.Time {
	return BeginningOfYear(t).AddDate(1, 0, 0).Add(-time.Nanosecond)
}