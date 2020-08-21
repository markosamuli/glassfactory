package dateutil

import (
	"encoding/json"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"gotest.tools/assert"
)

func TestDateTime_ParseDateTime(t *testing.T) {
	dt, err := ParseDateTime("")
	assert.ErrorContains(t, err, "cannot parse")
	assert.Equal(t, dt, DateTime{})
}

func TestDateTime_MarshalText(t *testing.T) {
	dt := DateTimeOf(time.Date(2018, time.June, 04, 5, 46, 15, 0, time.UTC))
	value, err := dt.MarshalText()
	assert.NilError(t, err)

	expectedText := "2018-06-04 05:46:15"
	assert.Equal(t, string(value), expectedText)
}

func TestDateTime_UnmarshalText(t *testing.T) {
	str := "2018-06-04 05:46:15"
	expected := DateTimeOf(time.Date(2018, time.June, 04, 5, 46, 15, 0, time.UTC))
	dt := DateTime{}
	err := dt.UnmarshalText([]byte(str))
	assert.NilError(t, err)
	assert.Equal(t, dt, expected)
}

func TestDateTime_MarshalJSON(t *testing.T) {
	var tests = []struct {
		name     string
		from     DateTime
		expected string
	}{
		{
			name:     "valid time",
			from:     DateTimeOf(time.Date(2018, time.June, 04, 5, 46, 15, 0, time.UTC)),
			expected: `"2018-06-04 05:46:15"`,
		},
		{
			name:     "default time",
			from:     DateTimeOf(time.Time{}),
			expected: `"0001-01-01 00:00:00"`,
		},
		{
			name: "invalid date part",
			from: DateTime{civil.DateTime{
				Date: civil.Date{},
				Time: civil.Time{Hour: 1},
			}},
			expected: `null`,
		},
		{
			name: "invalid time part",
			from: DateTime{civil.DateTime{
				Date: civil.Date{Year: 2019, Month: 1, Day: 1},
				Time: civil.Time{Hour: -1},
			}},
			expected: `null`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				value []byte
				err   error
			)
			value, err = json.Marshal(&tt.from)
			assert.NilError(t, err)
			assert.Equal(t, string(value), tt.expected)
		})
	}
}

func TestDateTime_UnmarshalJSON(t *testing.T) {
	var tests = []struct {
		name     string
		from     string
		expected time.Time
		err      string
		valid    bool
	}{
		{
			name:     "valid date and time",
			from:     `"2018-06-04 05:46:15"`,
			expected: time.Date(2018, time.June, 04, 5, 46, 15, 0, time.UTC),
			valid:    true,
		},
		{
			name: "missing date",
			from: `"05:46:15"`,
			err:  "cannot parse",
		},
		{
			name: "missing time",
			from: `"2018-06-04"`,
			err:  "cannot parse",
		},
		{
			name:  "null",
			from:  `null`,
			valid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dt DateTime
			var err error
			err = json.Unmarshal([]byte(tt.from), &dt)
			if tt.err != "" {
				assert.ErrorContains(t, err, tt.err)
			} else {
				assert.NilError(t, err)
				if tt.valid {
					assert.Equal(t, dt.In(tt.expected.Location()), tt.expected)
				} else {
					assert.Assert(t, !dt.IsValid())
				}
			}
		})
	}
}

func TestDateTime_Before(t *testing.T) {
	d1 := civil.Date{Year: 2016, Month: 12, Day: 31}
	d2 := civil.Date{Year: 2017, Month: 1, Day: 1}
	t1 := civil.Time{Hour: 5, Minute: 6, Second: 7, Nanosecond: 8}
	t2 := civil.Time{Hour: 5, Minute: 6, Second: 7, Nanosecond: 9}
	for _, test := range []struct {
		dt1, dt2 DateTime
		want     bool
	}{
		{DateTime{civil.DateTime{Date: d1, Time: t1}},
			DateTime{civil.DateTime{Date: d2, Time: t1}}, true},
		{DateTime{civil.DateTime{Date: d1, Time: t1}},
			DateTime{civil.DateTime{Date: d1, Time: t2}}, true},
		{DateTime{civil.DateTime{Date: d2, Time: t1}},
			DateTime{civil.DateTime{Date: d1, Time: t1}}, false},
		{DateTime{civil.DateTime{Date: d2, Time: t1}},
			DateTime{civil.DateTime{Date: d2, Time: t1}}, false},
	} {
		if got := test.dt1.Before(test.dt2); got != test.want {
			t.Errorf("%v.Before(%v): got %t, want %t", test.dt1, test.dt2, got, test.want)
		}
	}
}

func TestDateTime_After(t *testing.T) {
	d1 := civil.Date{Year: 2016, Month: 12, Day: 31}
	d2 := civil.Date{Year: 2017, Month: 1, Day: 1}
	t1 := civil.Time{Hour: 5, Minute: 6, Second: 7, Nanosecond: 8}
	t2 := civil.Time{Hour: 5, Minute: 6, Second: 7, Nanosecond: 9}
	for _, test := range []struct {
		dt1, dt2 DateTime
		want     bool
	}{
		{DateTime{civil.DateTime{Date: d1, Time: t1}},
			DateTime{civil.DateTime{Date: d2, Time: t1}}, false},
		{DateTime{civil.DateTime{Date: d1, Time: t1}},
			DateTime{civil.DateTime{Date: d1, Time: t2}}, false},
		{DateTime{civil.DateTime{Date: d2, Time: t1}},
			DateTime{civil.DateTime{Date: d1, Time: t1}}, true},
		{DateTime{civil.DateTime{Date: d2, Time: t1}},
			DateTime{civil.DateTime{Date: d2, Time: t1}}, false},
	} {
		if got := test.dt1.After(test.dt2); got != test.want {
			t.Errorf("%v.After(%v): got %t, want %t", test.dt1, test.dt2, got, test.want)
		}
	}
}
