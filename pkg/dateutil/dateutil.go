// Package dateutil provides tools and types for handling date and time information
package dateutil

import (
	"time"

	"github.com/jinzhu/now"
)

var nilTime = (time.Time{}).UnixNano()

// BeginningOfMonth returns the date of the first day of the month
//
// Deprecated: Use a toolkit like github.com/jinzhu/now instead
func BeginningOfMonth(t time.Time) time.Time {
	return now.With(t).BeginningOfMonth()
}

// EndOfMonth returns the date of the last day of the month
//
// Deprecated: Use a toolkit like github.com/jinzhu/now instead
func EndOfMonth(t time.Time) time.Time {
	return now.With(t).EndOfMonth()
}

// BeginningOfYear returns the date of the first day of the year
//
// Deprecated: Use a toolkit like github.com/jinzhu/now instead
func BeginningOfYear(t time.Time) time.Time {
	return now.With(t).BeginningOfYear()
}

// EndOfYear returns the date of the first day of the year
//
// Deprecated: Use a toolkit like github.com/jinzhu/now instead
func EndOfYear(t time.Time) time.Time {
	return now.With(t).EndOfYear()
}
