package dateutil

import (
	"fmt"
	"time"

	"cloud.google.com/go/civil"
)

const dateTimeLayout = "2006-01-02 15:04:05"

// DateTime is a native date time in YYYY-MM-DD HH:MM:SS format
type DateTime struct {
	civil.DateTime
}

// DateTimeOf constructs a DateTime type from given time
func DateTimeOf(t time.Time) DateTime {
	return DateTime{
		DateTime: civil.DateTimeOf(t),
	}
}

// ParseDateTime returns a DateTime type parsed from a string
func ParseDateTime(str string) (DateTime, error) {
	t, err := time.Parse(dateTimeLayout, str)
	if err != nil {
		return DateTime{}, err
	}
	return DateTimeOf(t), nil
}

// String returns the date in the format described in ParseDate.
func (dt DateTime) String() string {
	return dt.Date.String() + " " + dt.Time.String()
}

// Before reports whether dt occurs before dt2.
func (dt DateTime) Before(dt2 DateTime) bool {
	return dt.In(time.UTC).Before(dt2.In(time.UTC))
}

// After reports whether dt occurs after dt2.
func (dt DateTime) After(dt2 DateTime) bool {
	return dt2.Before(dt)
}

// MarshalText implements the encoding.TextMarshaler interface.
// The output is the result of dt.String().
func (dt DateTime) MarshalText() ([]byte, error) {
	return []byte(dt.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The date is expected to be a string in a format accepted by ParseDate.
func (dt *DateTime) UnmarshalText(data []byte) error {
	var err error
	*dt, err = ParseDateTime(string(data))
	return err
}

// MarshalJSON implements the json.Marshaler interface.
// The time is a quoted string.
func (dt DateTime) MarshalJSON() ([]byte, error) {
	if !dt.Date.IsValid() || !dt.Time.IsValid() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", dt.String())), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string.
func (dt *DateTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	var t time.Time
	t, err = time.Parse(`"`+dateTimeLayout+`"`, string(data))
	*dt = DateTimeOf(t)
	return err
}
