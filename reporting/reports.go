package reporting

import "github.com/markosamuli/glassfactory/model"

type TimeReportTotals struct {
	actual  float64
	planned float64
}

func FormatBillableStatus(billableStatus model.BillableStatus) string {
	return billableStatus.String()
}
