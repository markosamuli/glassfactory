package reports

type TimeReportTotals struct {
	actual float32
	planned float32
}

func FormatBillableStatus(billableStatus string) string {
	switch {
	case billableStatus == "billable":
		return "Billable"
	case billableStatus == "non_billable":
		return "Non Billable"
	case billableStatus == "new_business":
		return "New Business"
	default:
		return billableStatus
	}
}