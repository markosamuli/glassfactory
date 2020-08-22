package reporting

import "github.com/markosamuli/glassfactory/model"

// TimeReportTotals represents the total actual and planned hours
type TimeReportTotals struct {
	actual  float64
	planned float64
}

// FormatBillableStatus returns the BillableStatus field as a string
func FormatBillableStatus(billableStatus model.BillableStatus) string {
	return billableStatus.String()
}
