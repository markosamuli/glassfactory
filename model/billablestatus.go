package model

import "fmt"

//go:generate stringer -type=BillableStatus -linecomment

type BillableStatus int

const (
	Unknown BillableStatus = iota
	Billable
	NonBillable // Non Billable
	NewBusiness // New Business
)

// MarshalJSON implements the json.Marshaler interface.
func (i BillableStatus) MarshalJSON() ([]byte, error) {
	var s string
	switch i {
	case Unknown:
		s = "null"
	case Billable:
		s = `"billable"`
	case NonBillable:
		s = `"non_billable"`
	case NewBusiness:
		s = `"new_business"`
	default:
		return nil, fmt.Errorf("invalid BillableStatus %d", i)
	}
	return []byte(s), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (i *BillableStatus) UnmarshalJSON(data []byte) error {
	s := string(data)
	switch s {
	case "null":
		*i = Unknown
	case `"billable"`:
		*i = Billable
	case `"non_billable"`:
		*i = NonBillable
	case `"new_business"`:
		*i = NewBusiness
	default:
		return fmt.Errorf("invalid BillableStatus %s", s)
	}
	return nil
}
