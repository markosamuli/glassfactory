package reporting

import (
	"context"
	"errors"
	"time"

	"github.com/markosamuli/glassfactory"
	"github.com/markosamuli/glassfactory/dateutil"
	"github.com/markosamuli/glassfactory/model"
)

// NewService creates a new Service for reporting
func NewService(ctx context.Context, api *glassfactory.Service) (*Service, error) {
	if api == nil {
		return nil, errors.New("api is nil")
	}
	s := &Service{}
	s.api = api
	return s, nil
}

// ReportsService provides methods for fetching time report data from Glass Factory
type Service struct {
	api *glassfactory.Service
}

// MemberTimeReportsForFiscalYear queries Glass Factory and returns time reports for the given fiscal year
//func (s *Service) MemberTimeReportsForFiscalYear(userID int, fiscalYear *FiscalYear) *glassfactory.MemberTimeReportCalls {
//	r := s.api.Reports
//	now := time.Now()
//	start := fiscalYear.Start
//	end := fiscalYear.End
//	if end.After(now) {
//		end = now
//	}
//	calls := &glassfactory.MemberTimeReportCalls{s: r.s}
//	calls.userID = userID
//	for _, m := range dateutil.MonthsBetweenDates(start, end) {
//		c := r.MemberTimeReport(userID, m.Start, m.End)
//		calls.Append(c)
//	}
//	return calls
//}

// MonthlyMemberTimeReports queries Glass Factory and returns time reports for a full calendar year matching the given time
func (s *Service) MonthlyMemberTimeReports(userID int, t time.Time) ([]*MonthlyMemberTimeReport, error) {
	now := time.Now()
	start := dateutil.BeginningOfYear(t)
	end := dateutil.EndOfYear(t)
	if end.After(now) {
		end = now
	}
	responses, err := s.api.Reports.MemberTimeReportsBetweenDates(userID, start, end).Do()
	if err != nil {
		return nil, err
	}

	mtr := make([]*model.MemberTimeReport, 0)
	for _, response := range responses {
		for _, report := range response.Reports {
			client, err := s.api.Clients.Get(report.ClientID)
			if err != nil {
				return nil, err
			}
			project, err := s.api.Projects.Get(report.ProjectID)
			if err != nil {
				return nil, err
			}
			report.Client = client
			report.Project = project
			mtr = append(mtr, report)
		}
	}
	return MonthlyMemberTimeReports(mtr), nil
}

// MonthlyMemberTimeReports queries Glass Factory and returns time reports for the given fiscal year
func (s *Service) FiscalYearMemberTimeReports(userID int, fiscalYear *FiscalYear) ([]*FiscalYearMemberTimeReport, error) {
	now := time.Now()
	start := fiscalYear.Start
	end := fiscalYear.End
	if end.After(now) {
		end = now
	}

	responses, err := s.api.Reports.MemberTimeReportsBetweenDates(userID, start, end).Do()
	if err != nil {
		return nil, err
	}

	mtr := make([]*model.MemberTimeReport, 0)
	for _, response := range responses {
		for _, report := range response.Reports {
			client, err := s.api.Clients.Get(report.ClientID)
			if err != nil {
				return nil, err
			}
			project, err := s.api.Projects.Get(report.ProjectID)
			if err != nil {
				return nil, err
			}
			report.Client = client
			report.Project = project
			mtr = append(mtr, report)
		}
	}
	return FiscalYearMemberTimeReports(mtr, fiscalYear.End.Month()), nil
}