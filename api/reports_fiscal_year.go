package api

import (
	"cloud.google.com/go/civil"
	"github.com/markosamuli/glassfactory/dateutils"
	"github.com/markosamuli/glassfactory/reports"
	"io"
	"sort"
	"time"
)

func (r *ReportsService) FiscalYearMemberTimeReports(userID int, fiscalYear *dateutils.FiscalYear) ([]*FiscalYearMemberTimeReport, error) {
	responses, err := r.MemberTimeReportsForFiscalYear(userID, fiscalYear).Do()
	if err != nil {
		return nil, err
	}
	reports := make([]*MemberTimeReport, 0)
	for _, response := range responses {
		for _, report := range response.Reports {
			client, err := r.s.Clients.Get(report.ClientID)
			if err != nil {
				return nil, err
			}
			project, err := r.s.Projects.Get(report.ProjectID)
			if err != nil {
				return nil, err
			}
			report.Client = client
			report.Project = project
			reports = append(reports, report)
		}
	}
	return FiscalYearMemberTimeReports(reports, fiscalYear.End.Month()), nil
}

func (tr *FiscalYearMemberTimeReport) RenderTable(writer io.Writer) {
	reportGroups := make(map[string][]*reports.FiscalYearTimeReport)

	projectReports := ProjectMemberTimeReports(tr.Reports)
	for _, pr := range projectReports {
		r := &reports.FiscalYearTimeReport{
			FiscalYear: tr.FiscalYear,
			Client:     pr.Client,
			Project:    pr.Project,
			Planned:    pr.Planned(),
			Actual:     pr.Actual(),
		}
		billableStatus := r.BillableStatus()
		br, ok := reportGroups[billableStatus]
		if !ok {
			br = make([]*reports.FiscalYearTimeReport, 0)
		}
		br = append(br, r)
		reportGroups[billableStatus] = br
	}

	table := reports.NewFiscalYearTimeReportTableWriter(writer)
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

type FiscalYearMemberTimeReport struct {
	UserID     int
	FiscalYear dateutils.FiscalYear
	Start      civil.Date
	End        civil.Date
	Reports    []*MemberTimeReport
}

func NewFiscalYearMemberTimeReport(userID int, fy dateutils.FiscalYear) *FiscalYearMemberTimeReport {
	return &FiscalYearMemberTimeReport{
		UserID:     userID,
		FiscalYear: fy,
		Reports:    make([]*MemberTimeReport, 0),
	}
}

func (tr *FiscalYearMemberTimeReport) Append(r *MemberTimeReport) {
	if !tr.Start.IsValid() || r.Date.Before(tr.Start) {
		tr.Start = r.Date
	}
	if !tr.End.IsValid() || r.Date.Before(tr.End) {
		tr.End = r.Date
	}
	tr.Reports = append(tr.Reports, r)
}

func (tr *FiscalYearMemberTimeReport) Planned() float32 {
	var planned float32
	for _, r := range tr.Reports {
		planned += r.Planned
	}
	return planned
}

func (tr *FiscalYearMemberTimeReport) Actual() float32 {
	var actual float32
	for _, r := range tr.Reports {
		actual += r.Actual
	}
	return actual
}

func FiscalYearMemberTimeReports(reports []*MemberTimeReport, finalMonth time.Month) []*FiscalYearMemberTimeReport {
	periods := make(map[dateutils.FiscalYear]*FiscalYearMemberTimeReport, 0)
	for _, r := range reports {
		fy := *dateutils.NewFiscalYear(r.Date.In(time.Local), finalMonth)
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

// ByCalendarMonth implements sort.Interface based on the CalendarMonth field.
type ByFiscalYear []*FiscalYearMemberTimeReport

func (a ByFiscalYear) Len() int           { return len(a) }
func (a ByFiscalYear) Less(i, j int) bool { return a[i].FiscalYear.Before(a[j].FiscalYear) }
func (a ByFiscalYear) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

