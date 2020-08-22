package reporting

import (
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/jinzhu/now"
	"github.com/markosamuli/glassfactory/model"
	"github.com/markosamuli/glassfactory/pkg/dateutil"
	"github.com/olekukonko/tablewriter"
)

// FiscalYear represents a time range for a fiscal year
type FiscalYear struct {
	Start time.Time
	End   time.Time
}

// String returns the fiscal year in FY YYYY format.
func (fy FiscalYear) String() string {
	return fmt.Sprintf("FY %04d", fy.End.Year())
}

// Before reports whether fy occurs before fy2.
func (fy FiscalYear) Before(fy2 FiscalYear) bool {
	return fy.End.Before(fy2.End)
}

// After reports whether fy occurs after fy2.
func (fy FiscalYear) After(fy2 FiscalYear) bool {
	return fy2.Before(fy)
}

// NewFiscalYear returns new FiscalYear for the given date ending at the given month
func NewFiscalYear(d time.Time, finalMonth time.Month) *FiscalYear {
	var start time.Time
	var end time.Time
	if finalMonth < time.December {
		start = time.Date(d.Year()-1, finalMonth+1, 1, 0, 0, 0, 0, time.Local)
		end = now.With(time.Date(d.Year(), finalMonth, 1, 23, 59, 59, 999999999, time.Local)).EndOfMonth()
	} else {
		start = time.Date(d.Year(), time.January, 1, 0, 0, 0, 0, time.Local)
		end = now.With(time.Date(d.Year(), time.December, 1, 23, 59, 59, 999999999, time.Local)).EndOfMonth()
	}
	if d.Before(start) {
		start = start.AddDate(-1, 0, 0)
		end = end.AddDate(-1, 0, 0)
	} else if d.After(end) {
		start = start.AddDate(1, 0, 0)
		end = end.AddDate(1, 0, 0)
	}
	return &FiscalYear{
		Start: start,
		End:   end,
	}
}

// FiscalYearMemberTimeReports convers MemberTimeReport data into FiscalYearMemberTimeReport
func FiscalYearMemberTimeReports(reports []*model.MemberTimeReport, finalMonth time.Month) []*FiscalYearMemberTimeReport {
	periods := make(map[FiscalYear]*FiscalYearMemberTimeReport, 0)
	for _, r := range reports {
		fy := *NewFiscalYear(r.Date.In(time.Local), finalMonth)
		p, ok := periods[fy]
		if !ok {
			p = NewFiscalYearMemberTimeReport(r.UserID, fy)
			periods[fy] = p
		}
		p.Append(r)
	}
	fyr := make([]*FiscalYearMemberTimeReport, 0, len(periods))
	for _, p := range periods {
		fyr = append(fyr, p)
	}
	sort.Sort(ByFiscalYear(fyr))
	return fyr
}

// ByFiscalYear implements sort.Interface based on the FiscalYear field.
type ByFiscalYear []*FiscalYearMemberTimeReport

func (a ByFiscalYear) Len() int           { return len(a) }
func (a ByFiscalYear) Less(i, j int) bool { return a[i].FiscalYear.Before(a[j].FiscalYear) }
func (a ByFiscalYear) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// FiscalYearMemberTimeReport represents MemberTimeReport data for a given fiscal year
type FiscalYearMemberTimeReport struct {
	UserID     int
	FiscalYear FiscalYear
	Start      dateutil.Date
	End        dateutil.Date
	Reports    []*model.MemberTimeReport
}

// NewFiscalYearMemberTimeReport creates FiscalYearMemberTimeReport for a user and given fiscal year
func NewFiscalYearMemberTimeReport(userID int, fy FiscalYear) *FiscalYearMemberTimeReport {
	return &FiscalYearMemberTimeReport{
		UserID:     userID,
		FiscalYear: fy,
		Reports:    make([]*model.MemberTimeReport, 0),
	}
}

// RenderTable displays FiscalYearMemberTimeReport in using NewFiscalYearTimeReportTableWriter
func (tr *FiscalYearMemberTimeReport) RenderTable(writer io.Writer) {
	reportGroups := make(map[string][]*FiscalYearTimeReport)

	projectReports := ProjectMemberTimeReports(tr.Reports)
	for _, pr := range projectReports {
		r := &FiscalYearTimeReport{
			FiscalYear: tr.FiscalYear,
			Client:     pr.Client,
			Project:    pr.Project,
			Planned:    pr.Planned(),
			Actual:     pr.Actual(),
		}
		billableStatus := r.BillableStatus()
		br, ok := reportGroups[billableStatus]
		if !ok {
			br = make([]*FiscalYearTimeReport, 0)
		}
		br = append(br, r)
		reportGroups[billableStatus] = br
	}

	table := NewFiscalYearTimeReportTableWriter(writer)
	for _, r := range reportGroups {
		sort.SliceStable(r, func(i, j int) bool {
			if r[i].Client.ID != r[j].Client.ID {
				return r[i].Client.ID < r[i].Client.ID
			}
			if r[i].Client.OfficeID != r[j].Client.OfficeID {
				return r[i].Client.OfficeID < r[i].Client.OfficeID
			}
			if r[i].Project.ID != r[j].Project.ID {
				return r[i].Project.ID < r[i].Project.ID
			}
			return true
		})
		for _, tr := range r {
			table.Append(tr)
		}
	}
	table.Render()
}

// Append a MemberTimeReport to the FiscalYearMemberTimeReport
func (tr *FiscalYearMemberTimeReport) Append(r *model.MemberTimeReport) {
	if !tr.Start.IsValid() || r.Date.Before(tr.Start) {
		tr.Start = r.Date
	}
	if !tr.End.IsValid() || r.Date.Before(tr.End) {
		tr.End = r.Date
	}
	tr.Reports = append(tr.Reports, r)
}

// Planned returns total planned hours
func (tr *FiscalYearMemberTimeReport) Planned() float64 {
	var planned float64
	for _, r := range tr.Reports {
		planned += r.Planned
	}
	return planned
}

// Actual returns total actual hours
func (tr *FiscalYearMemberTimeReport) Actual() float64 {
	var actual float64
	for _, r := range tr.Reports {
		actual += r.Actual
	}
	return actual
}

// FiscalYearTimeReport represents fiscal year totals for a given client and project
type FiscalYearTimeReport struct {
	FiscalYear FiscalYear
	Client     *model.Client
	Project    *model.Project
	Planned    float64
	Actual     float64
}

// BillableStatus returns project's billable status
func (r *FiscalYearTimeReport) BillableStatus() string {
	return FormatBillableStatus(r.Project.BillableStatus)
}

// FiscalYearTimeReportTableWriter is used for displaying reports in table format
type FiscalYearTimeReportTableWriter struct {
	table  *tablewriter.Table
	totals map[string]*TimeReportTotals
}

// NewFiscalYearTimeReportTableWriter creates a new FiscalYearTimeReportTableWriter
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
		table:  table,
		totals: make(map[string]*TimeReportTotals),
	}
}

// Append adds FiscalYearTimeReport data to the table and updates the total hours
func (t *FiscalYearTimeReportTableWriter) Append(r *FiscalYearTimeReport) {
	billable := r.BillableStatus()
	t.table.Append([]string{
		fmt.Sprintf("%s", r.FiscalYear),
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

// Render displays the report data in a table format
func (t *FiscalYearTimeReportTableWriter) Render() {
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
