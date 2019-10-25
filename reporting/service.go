package reporting

import (
	"context"
	"errors"
	"time"

	"github.com/markosamuli/glassfactory/dateutil"
	"github.com/markosamuli/glassfactory/glassfactory"
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

// MonthlyMemberTimeReports queries Glass Factory and returns time reports for a full calendar year matching the given time
func (s *Service) MonthlyMemberTimeReports(userID int, t time.Time) ([]*MonthlyMemberTimeReport, error) {
	start := dateutil.BeginningOfYear(t)
	end := dateutil.EndOfYear(t)
	reports, err := s.api.Member.Reports.GetTimeReportsBetweenDates(userID, start, end, glassfactory.FetchRelated())
	if err != nil {
		return nil, err
	}
	return MonthlyMemberTimeReports(reports), nil
}

// MonthlyMemberTimeReports queries Glass Factory and returns time reports for the given fiscal year
func (s *Service) FiscalYearMemberTimeReports(userID int, fiscalYear *FiscalYear) ([]*FiscalYearMemberTimeReport, error) {
	start := fiscalYear.Start
	end := fiscalYear.End
	reports, err := s.api.Member.Reports.GetTimeReportsBetweenDates(userID, start, end, glassfactory.FetchRelated())
	if err != nil {
		return nil, err
	}
	return FiscalYearMemberTimeReports(reports, fiscalYear.End.Month()), nil
}