package reports

import (
	"fmt"
	"github.com/markosamuli/glassfactory/dateutils"
	"github.com/markosamuli/glassfactory/models"
	"github.com/olekukonko/tablewriter"
	"io"
	"os"
)

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

type TimeReportTotals struct {
	actual float32
	planned float32
}

// Monthly time reports

type MonthlyTimeReport struct {
	CalendarMonth dateutils.CalendarMonth
	Client *models.Client
	Project *models.Project
	Planned float32
	Actual float32
}

func (r *MonthlyTimeReport) BillableStatus() string {
	return FormatBillableStatus(r.Project.BillableStatus)
}

type MonthlyTimeReportTableWriter struct {
	table *tablewriter.Table
	totals map[string]*TimeReportTotals
}

func NewMonthlyTimeReportTableWriter(writer io.Writer) *MonthlyTimeReportTableWriter {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Month",
		"Billable",
		"Client",
		"Project",
		"Actual",
		"Planned",
		"Diff",
	})
	table.SetAutoMergeCells(false)
	table.SetRowLine(true)
	return &MonthlyTimeReportTableWriter{
		table: table,
		totals: make(map[string]*TimeReportTotals),
	}
}

func (t *MonthlyTimeReportTableWriter) Append(r *MonthlyTimeReport) {
	billable := r.BillableStatus()
	t.table.Append([]string{
		fmt.Sprintf("%s", r.CalendarMonth),
		billable,
		r.Client.Name,
		r.Project.Name,
		fmt.Sprintf("%6.2f ", r.Actual),
		fmt.Sprintf("%6.2f ", r.Planned),
		fmt.Sprintf("%6.2f ", r.Actual - r.Planned),
	})
	totals, ok := t.totals[billable]
	if !ok {
		totals = &TimeReportTotals{planned: 0.0, actual: 0.0}
	}
	totals.planned += r.Planned
	totals.actual += r.Actual
	t.totals[billable] = totals
}

func (t *MonthlyTimeReportTableWriter) Render() {
	var planned float32
	var actual float32
	for billable, totals := range t.totals {
		totalHeader := fmt.Sprintf("Total %s", billable)
		t.table.Append([]string{
			"",
			"",
			"",
			totalHeader,
			fmt.Sprintf("%6.2f ", totals.actual),
			fmt.Sprintf("%6.2f ", totals.planned),
			fmt.Sprintf("%6.2f ", totals.actual - totals.planned),
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
		fmt.Sprintf("%6.2f ", actual - planned),
	})
	t.table.Render()
}

// Fiscal year time reports

type FiscalYearTimeReport struct {
	FiscalYear dateutils.FiscalYear
	Client *models.Client
	Project *models.Project
	Planned float32
	Actual float32
}

func (r *FiscalYearTimeReport) BillableStatus() string {
	return FormatBillableStatus(r.Project.BillableStatus)
}

type FiscalYearTimeReportTableWriter struct {
	table *tablewriter.Table
	totals map[string]*TimeReportTotals
}

func NewFiscalYearTimeReportTableWriter(writer io.Writer) *FiscalYearTimeReportTableWriter {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Fiscal Year",
		"Billable",
		"Client",
		"Project",
		"Actual",
		"Planned",
		"Diff",
	})
	table.SetAutoMergeCells(false)
	table.SetRowLine(true)
	return &FiscalYearTimeReportTableWriter{
		table: table,
		totals: make(map[string]*TimeReportTotals),
	}
}

func (t *FiscalYearTimeReportTableWriter) Append(r *FiscalYearTimeReport) {
	billable := r.BillableStatus()
	t.table.Append([]string{
		fmt.Sprintf("%s", r.FiscalYear),
		billable,
		r.Client.Name,
		r.Project.Name,
		fmt.Sprintf("%6.2f ", r.Actual),
		fmt.Sprintf("%6.2f ", r.Planned),
		fmt.Sprintf("%6.2f ", r.Actual - r.Planned),
	})
	totals, ok := t.totals[billable]
	if !ok {
		totals = &TimeReportTotals{planned: 0.0, actual: 0.0}
	}
	totals.planned += r.Planned
	totals.actual += r.Actual
	t.totals[billable] = totals
}

func (t *FiscalYearTimeReportTableWriter) Render() {
	var planned float32
	var actual float32
	for billable, totals := range t.totals {
		totalHeader := fmt.Sprintf("Total %s", billable)
		t.table.Append([]string{
			"",
			"",
			"",
			totalHeader,
			fmt.Sprintf("%6.2f ", totals.actual),
			fmt.Sprintf("%6.2f ", totals.planned),
			fmt.Sprintf("%6.2f ", totals.actual - totals.planned),
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
		fmt.Sprintf("%6.2f ", actual - planned),
	})
	t.table.Render()
}

// Annual time reports

type AnnualTimeReport struct {
	Year int
	Client *models.Client
	Project *models.Project
	Planned float32
	Actual float32
}

func (r *AnnualTimeReport) BillableStatus() string {
	return FormatBillableStatus(r.Project.BillableStatus)
}

type AnnualTimeReportTableWriter struct {
	table *tablewriter.Table
	totals map[string]*TimeReportTotals
}

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
		table: table,
		totals: make(map[string]*TimeReportTotals),
	}
}

func (t *AnnualTimeReportTableWriter) Append(r *AnnualTimeReport) {
	billable := r.BillableStatus()
	t.table.Append([]string{
		fmt.Sprintf("%d", r.Year),
		billable,
		r.Client.Name,
		r.Project.Name,
		fmt.Sprintf("%6.2f ", r.Actual),
		fmt.Sprintf("%6.2f ", r.Planned),
		fmt.Sprintf("%6.2f ", r.Actual - r.Planned),
	})
	totals, ok := t.totals[billable]
	if !ok {
		totals = &TimeReportTotals{planned: 0.0, actual: 0.0}
	}
	totals.planned += r.Planned
	totals.actual += r.Actual
	t.totals[billable] = totals
}

func (t *AnnualTimeReportTableWriter) Render() {
	var planned float32
	var actual float32
	for billable, totals := range t.totals {
		totalHeader := fmt.Sprintf("Total %s", billable)
		t.table.Append([]string{
			"",
			"",
			"",
			totalHeader,
			fmt.Sprintf("%6.2f ", totals.actual),
			fmt.Sprintf("%6.2f ", totals.planned),
			fmt.Sprintf("%6.2f ", totals.actual - totals.planned),
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
		fmt.Sprintf("%6.2f ", actual - planned),
	})
	t.table.Render()
}
