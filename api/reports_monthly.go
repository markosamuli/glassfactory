package api

import (
	"cloud.google.com/go/civil"
	"github.com/markosamuli/glassfactory/dateutils"
	"github.com/markosamuli/glassfactory/reports"
	"io"
	"sort"
	"time"
)

func (r *ReportsService) MonthlyMemberTimeReports(userID int, t time.Time) ([]*MonthlyMemberTimeReport, error) {
	responses, err := r.MemberTimeReportsForYear(userID, t).Do()
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
	return MonthlyMemberTimeReports(reports), nil
}

func (tr *MonthlyMemberTimeReport) RenderTable(writer io.Writer) {
	reportGroups := make(map[string][]*reports.MonthlyTimeReport)
	projectReports := ProjectMemberTimeReports(tr.Reports)
	for _, pr := range projectReports {
		r := &reports.MonthlyTimeReport{
			CalendarMonth: tr.CalendarMonth,
			Client:        pr.Client,
			Project:       pr.Project,
			Planned:       pr.Planned(),
			Actual:        pr.Actual(),
		}
		billableStatus := r.BillableStatus()
		br, ok := reportGroups[billableStatus]
		if !ok {
			br = make([]*reports.MonthlyTimeReport, 0)
		}
		br = append(br, r)
		reportGroups[billableStatus] = br
	}

	table := reports.NewMonthlyTimeReportTableWriter(writer)
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

type MonthlyMemberTimeReport struct {
	UserID int
	CalendarMonth dateutils.CalendarMonth
	Start civil.Date
	End civil.Date
	Reports []*MemberTimeReport
}

func NewMonthlyMemberTimeReport(userID int, month dateutils.CalendarMonth) *MonthlyMemberTimeReport {
	return &MonthlyMemberTimeReport{
		UserID: userID,
		CalendarMonth: month,
		Reports: make([]*MemberTimeReport, 0),
	}
}

func (tr *MonthlyMemberTimeReport) Append(r *MemberTimeReport) {
	if !tr.Start.IsValid() || r.Date.Before(tr.Start) {
		tr.Start = r.Date
	}
	if !tr.End.IsValid() || r.Date.Before(tr.End) {
		tr.End = r.Date
	}
	tr.Reports = append(tr.Reports, r)
}

func (tr *MonthlyMemberTimeReport) Planned() float32 {
	var planned float32
	for  _, r := range tr.Reports {
		planned += r.Planned
	}
	return planned
}

func (tr *MonthlyMemberTimeReport) Actual() float32 {
	var actual float32
	for  _, r := range tr.Reports {
		actual += r.Actual
	}
	return actual
}

func MonthlyMemberTimeReports(reports []*MemberTimeReport) []*MonthlyMemberTimeReport {
	months := make(map[dateutils.CalendarMonth]*MonthlyMemberTimeReport, 0)
	for _, r := range reports {
		month := r.CalendarMonth()
		mr, ok := months[month]
		if !ok {
			mr = NewMonthlyMemberTimeReport(r.UserID, month)
			months[month] = mr
		}
		mr.Append(r)
	}
	mr := make([]*MonthlyMemberTimeReport, 0, len(months))
	for  _, p := range months {
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