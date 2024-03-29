package reporting

import (
	"fmt"
	"io"
	"os"

	"github.com/markosamuli/glassfactory/model"
	"github.com/olekukonko/tablewriter"
)

// AnnualTimeReport represents a project time report for a calendar year
type AnnualTimeReport struct {
	Year    int
	Client  *model.Client
	Project *model.Project
	Planned float64
	Actual  float64
}

// BillableStatus formats project billable status as a string
func (r *AnnualTimeReport) BillableStatus() string {
	return FormatBillableStatus(r.Project.BillableStatus)
}

// AnnualTimeReportTableWriter is used for printing annual time reports
type AnnualTimeReportTableWriter struct {
	table  *tablewriter.Table
	totals map[string]*TimeReportTotals
}

// NewAnnualTimeReportTableWriter creates writer for annual time reports
func NewAnnualTimeReportTableWriter(writer io.Writer) *AnnualTimeReportTableWriter {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Year",
		"Billable",
		"Client",
		"Project",
		"Actual",
		"Planned",
		"Diff",
	})
	table.SetAutoMergeCells(false)
	table.SetRowLine(true)
	return &AnnualTimeReportTableWriter{
		table:  table,
		totals: make(map[string]*TimeReportTotals),
	}
}

// Append adds annual time report data to the report
func (t *AnnualTimeReportTableWriter) Append(r *AnnualTimeReport) {
	billable := r.BillableStatus()
	t.table.Append([]string{
		fmt.Sprintf("%d", r.Year),
		billable,
		r.Client.Name,
		r.Project.Name,
		fmt.Sprintf("%6.2f ", r.Actual),
		fmt.Sprintf("%6.2f ", r.Planned),
		fmt.Sprintf("%6.2f ", r.Actual-r.Planned),
	})
	totals, ok := t.totals[billable]
	if !ok {
		totals = &TimeReportTotals{planned: 0.0, actual: 0.0}
	}
	totals.planned += r.Planned
	totals.actual += r.Actual
	t.totals[billable] = totals
}

// Render the annual time report data
func (t *AnnualTimeReportTableWriter) Render() {
	var planned float64
	var actual float64
	for billable, totals := range t.totals {
		totalHeader := fmt.Sprintf("Total %s", billable)
		t.table.Append([]string{
			"",
			"",
			"",
			totalHeader,
			fmt.Sprintf("%6.2f ", totals.actual),
			fmt.Sprintf("%6.2f ", totals.planned),
			fmt.Sprintf("%6.2f ", totals.actual-totals.planned),
		})
		planned += totals.planned
		actual += totals.actual
	}
	t.table.SetFooter([]string{
		"",
		"",
		"",
		"Total",
		fmt.Sprintf("%6.2f ", actual),
		fmt.Sprintf("%6.2f ", planned),
		fmt.Sprintf("%6.2f ", actual-planned),
	})
	t.table.Render()
}
