package reporting

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/markosamuli/glassfactory/dateutil"
	"github.com/markosamuli/glassfactory/model"
	"github.com/olekukonko/tablewriter"
)

type MonthlyMemberTimeReport struct {
	UserID        int
	CalendarMonth CalendarMonth
	Start         dateutil.Date
	End           dateutil.Date
	Reports       []*model.MemberTimeReport
}

func NewMonthlyMemberTimeReport(userID int, month CalendarMonth) *MonthlyMemberTimeReport {
	return &MonthlyMemberTimeReport{
		UserID:        userID,
		CalendarMonth: month,
		Reports:       make([]*model.MemberTimeReport, 0),
	}
}

func (tr *MonthlyMemberTimeReport) Append(r *model.MemberTimeReport) {
	if !tr.Start.IsValid() || r.Date.Before(tr.Start) {
		tr.Start = r.Date
	}
	if !tr.End.IsValid() || r.Date.Before(tr.End) {
		tr.End = r.Date
	}
	tr.Reports = append(tr.Reports, r)
}

func (tr *MonthlyMemberTimeReport) Planned() float64 {
	var planned float64
	for _, r := range tr.Reports {
		planned += r.Planned
	}
	return planned
}

func (tr *MonthlyMemberTimeReport) Actual() float64 {
	var actual float64
	for _, r := range tr.Reports {
		actual += r.Actual
	}
	return actual
}

func MonthlyMemberTimeReports(reports []*model.MemberTimeReport) []*MonthlyMemberTimeReport {
	months := make(map[CalendarMonth]*MonthlyMemberTimeReport, 0)
	for _, r := range reports {
		month := CalendarMonth{
			Year:  r.Date.Year,
			Month: r.Date.Month,
		}
		mr, ok := months[month]
		if !ok {
			mr = NewMonthlyMemberTimeReport(r.UserID, month)
			months[month] = mr
		}
		mr.Append(r)
	}
	mr := make([]*MonthlyMemberTimeReport, 0, len(months))
	for _, p := range months {
		mr = append(mr, p)
	}
	sort.Sort(ByCalendarMonth(mr))
	return mr
}

// ByCalendarMonth implements sort.Interface based on the CalendarMonth field.
type ByCalendarMonth []*MonthlyMemberTimeReport

func (a ByCalendarMonth) Len() int           { return len(a) }
func (a ByCalendarMonth) Less(i, j int) bool { return a[i].CalendarMonth.Before(a[j].CalendarMonth) }
func (a ByCalendarMonth) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }


type MonthlyTimeReport struct {
	CalendarMonth CalendarMonth
	Client        *model.Client
	Project       *model.Project
	Planned       float64
	Actual        float64
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

func (tr *MonthlyMemberTimeReport) RenderTable(writer io.Writer) {
	reportGroups := make(map[string][]*MonthlyTimeReport)
	projectReports := ProjectMemberTimeReports(tr.Reports)
	for _, pr := range projectReports {
		r := &MonthlyTimeReport{
			CalendarMonth: tr.CalendarMonth,
			Client:        pr.Client,
			Project:       pr.Project,
			Planned:       pr.Planned(),
			Actual:        pr.Actual(),
		}
		billableStatus := r.BillableStatus()
		br, ok := reportGroups[billableStatus]
		if !ok {
			br = make([]*MonthlyTimeReport, 0)
		}
		br = append(br, r)
		reportGroups[billableStatus] = br
	}

	table := NewMonthlyTimeReportTableWriter(writer)
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
