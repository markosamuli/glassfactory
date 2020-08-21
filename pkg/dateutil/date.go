package dateutil

import (
	"fmt"
	"time"

	"cloud.google.com/go/civil"
)

const dateLayout = "2006-01-02"

// Date is a native date in YYYY-MM-DD format
type Date struct {
	civil.Date
}

// DateOf constructs a Date type from given time
func DateOf(t time.Time) Date {
	var d Date
	d.Year, d.Month, d.Day = t.Date()
	return d
}

// ParseDate returns a Date type parsed from a string
func ParseDate(s string) (Date, error) {
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return Date{}, err
	}
	return DateOf(t), nil
}

// Before reports whether dt occurs before dt2.
func (d Date) Before(d2 Date) bool {
	return d.In(time.UTC).Before(d2.In(time.UTC))
}

// After reports whether dt occurs after dt2.
func (d Date) After(d2 Date) bool {
	return d2.Before(d)
}

// MarshalText implements the encoding.TextMarshaler interface.
// The output is the result of d.String().
func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The date is expected to be a string in a format accepted by ParseDate.
func (d *Date) UnmarshalText(data []byte) error {
	var err error
	*d, err = ParseDate(string(data))
	return err
}

// MarshalJSON implements the json.Marshaler interface.
// The time is a quoted string.
func (d Date) MarshalJSON() ([]byte, error) {
	if !d.IsValid() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", d.String())), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string.
func (d *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	var t time.Time
	t, err = time.Parse(`"`+dateLayout+`"`, string(data))
	*d = DateOf(t)
	return err
}
