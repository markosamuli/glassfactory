package dateutil

import (
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestMonthsBetweenDates(t *testing.T) {
	start := time.Date(2017, time.January, 30, 0, 0, 0, 0, time.UTC)
	end := time.Date(2017, time.April, 2, 0, 0, 0, 0, time.UTC)
	expected := []DateRange{
		{
			Start: time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2017, time.January, 31, 23, 59, 59, 999999999, time.UTC),
		},
		{
			Start: time.Date(2017, time.February, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2017, time.February, 28, 23, 59, 59, 999999999, time.UTC),
		},
		{
			Start: time.Date(2017, time.March, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2017, time.March, 31, 23, 59, 59, 999999999, time.UTC),
		},
		{
			Start: time.Date(2017, time.April, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2017, time.April, 30, 23, 59, 59, 999999999, time.UTC),
		},
	}
	months := MonthsBetweenDates(start, end)
	assert.DeepEqual(t, months, expected)
}
