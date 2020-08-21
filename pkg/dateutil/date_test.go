package dateutil

import (
	"encoding/json"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"gotest.tools/assert"
)

func TestDate_ParseDate(t *testing.T) {
	dt, err := ParseDate("")
	assert.ErrorContains(t, err, "cannot parse")
	assert.Equal(t, dt, Date{})
}

func TestDate_MarshalText(t *testing.T) {
	dt := DateOf(time.Date(2018, time.June, 04, 0, 0, 0, 0, time.UTC))
	value, err := dt.MarshalText()
	assert.NilError(t, err)

	expectedText := "2018-06-04"
	assert.Equal(t, string(value), expectedText)
}

func TestDate_UnmarshalText(t *testing.T) {
	str := "2018-06-04"
	expected := DateOf(time.Date(2018, time.June, 04, 0, 0, 0, 0, time.UTC))
	dt := Date{}
	err := dt.UnmarshalText([]byte(str))
	assert.NilError(t, err)
	assert.Equal(t, dt, expected)
}

func TestDate_MarshalJSON(t *testing.T) {
	var tests = []struct {
		name     string
		from     Date
		expected string
	}{
		{
			name:     "valid time",
			from:     DateOf(time.Date(2018, time.June, 04, 0, 0, 0, 0, time.UTC)),
			expected: `"2018-06-04"`,
		},
		{
			name:     "default time",
			from:     DateOf(time.Time{}),
			expected: `"0001-01-01"`,
		},
		{
			name: "invalid date",
			from: Date{
				Date: civil.Date{},
			},
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

func TestDate_UnmarshalJSON(t *testing.T) {
	var err error
	var str string
	var expectedTime time.Time

	type TestType struct {
		Date Date `json:"date"`
	}
	expectedTime = time.Date(2018, time.June, 04, 0, 0, 0, 0, time.UTC)

	var dt *Date
	str = `"2018-06-04"`
	err = json.Unmarshal([]byte(str), &dt)
	assert.NilError(t, err)
	assert.Equal(t, dt.Date.In(time.UTC), expectedTime)

	var obj1 *TestType
	str = `{
		"date": "2018-06-04"
	}`
	err = json.Unmarshal([]byte(str), &obj1)
	assert.NilError(t, err)
	assert.Equal(t, obj1.Date.In(time.UTC), expectedTime)

	var obj2 *TestType
	str = `{
		"date": null
	}`
	err = json.Unmarshal([]byte(str), &obj2)
	assert.NilError(t, err)
	assert.Assert(t, !obj2.Date.IsValid())
}

func TestDate_Before(t *testing.T) {
	for _, test := range []struct {
		d1, d2 Date
		want   bool
	}{
		{Date{civil.Date{Year: 2016, Month: 12, Day: 31}},
			Date{civil.Date{Year: 2017, Month: 1, Day: 1}}, true},
		{Date{civil.Date{Year: 2016, Month: 1, Day: 1}},
			Date{civil.Date{Year: 2016, Month: 1, Day: 1}}, false},
		{Date{civil.Date{Year: 2016, Month: 12, Day: 30}},
			Date{civil.Date{Year: 2016, Month: 12, Day: 31}}, true},
	} {
		if got := test.d1.Before(test.d2); got != test.want {
			t.Errorf("%v.Before(%v): got %t, want %t", test.d1, test.d2, got, test.want)
		}
	}
}

func TestDate_After(t *testing.T) {
	for _, test := range []struct {
		d1, d2 Date
		want   bool
	}{
		{Date{civil.Date{Year: 2016, Month: 12, Day: 31}},
			Date{civil.Date{Year: 2017, Month: 1, Day: 1}}, false},
		{Date{civil.Date{Year: 2016, Month: 1, Day: 1}},
			Date{civil.Date{Year: 2016, Month: 1, Day: 1}}, false},
		{Date{civil.Date{Year: 2016, Month: 12, Day: 30}},
			Date{civil.Date{Year: 2016, Month: 12, Day: 31}}, false},
	} {
		if got := test.d1.After(test.d2); got != test.want {
			t.Errorf("%v.After(%v): got %t, want %t", test.d1, test.d2, got, test.want)
		}
	}
}
