package dateutil

import (
	"testing"
	"time"
)

func TestBeginningOfMonth(t *testing.T) {
	for _, test := range []struct {
		given time.Time
		want  time.Time
	}{
		{
			time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			time.Date(2019, time.February, 28, 1, 2, 3, 0, time.UTC),
			time.Date(2019, time.February, 1, 0, 0, 0, 0, time.UTC),
		},
	} {
		if got := BeginningOfMonth(test.given); got != test.want {
			t.Errorf("BeginningOfMonth(%v): got %v, want %v", test.given, got, test.want)
		}
	}
}

func TestEndOfMonth(t *testing.T) {
	for _, test := range []struct {
		given time.Time
		want  time.Time
	}{
		{
			time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2019, time.January, 31, 23, 59, 59, 999999999, time.UTC),
		},
		{
			time.Date(2019, time.February, 1, 1, 2, 3, 0, time.UTC),
			time.Date(2019, time.February, 28, 23, 59, 59, 999999999, time.UTC),
		},
		{
			time.Date(2019, time.April, 15, 0, 0, 0, 0, time.UTC),
			time.Date(2019, time.April, 30, 23, 59, 59, 999999999, time.UTC),
		},
	} {
		if got := EndOfMonth(test.given); got != test.want {
			t.Errorf("EndOfMonth(%v): got %v, want %v", test.given, got, test.want)
		}
	}
}

func TestBeginningOfYear(t *testing.T) {
	for _, test := range []struct {
		given time.Time
		want  time.Time
	}{
		{
			time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			time.Date(2018, time.February, 1, 1, 2, 3, 0, time.UTC),
			time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			time.Date(2019, time.December, 31, 23, 59, 59, 999999999, time.UTC),
			time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	} {
		if got := BeginningOfYear(test.given); got != test.want {
			t.Errorf("BeginningOfYear(%v): got %v, want %v", test.given, got, test.want)
		}
	}
}

func TestEndOfYear(t *testing.T) {
	for _, test := range []struct {
		given time.Time
		want  time.Time
	}{
		{
			time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2017, time.December, 31, 23, 59, 59, 999999999, time.UTC),
		},
		{
			time.Date(2018, time.February, 1, 1, 2, 3, 0, time.UTC),
			time.Date(2018, time.December, 31, 23, 59, 59, 999999999, time.UTC),
		},
		{
			time.Date(2019, time.April, 15, 0, 0, 0, 0, time.UTC),
			time.Date(2019, time.December, 31, 23, 59, 59, 999999999, time.UTC),
		},
	} {
		if got := EndOfYear(test.given); got != test.want {
			t.Errorf("EndOfYear(%v): got %v, want %v", test.given, got, test.want)
		}
	}
}
