package reporting

import (
	"context"
	"errors"
	"time"

	"github.com/jinzhu/now"
	"github.com/markosamuli/glassfactory/api"
)

// NewService creates a new Service for reporting
func NewService(ctx context.Context, apiService *api.Service) (*Service, error) {
	if apiService == nil {
		return nil, errors.New("apiService is nil")
	}
	s := &Service{}
	s.api = apiService
	return s, nil
}

// Service provides methods for fetching time report data from Glass Factory
type Service struct {
	api *api.Service
}

// MonthlyMemberTimeReports queries Glass Factory and returns time reports for a full calendar year matching the given time
func (s *Service) MonthlyMemberTimeReports(userID int, t time.Time) ([]*MonthlyMemberTimeReport, error) {
	start := now.With(t).BeginningOfYear()
	end := now.With(t).EndOfYear()
	reports, err := s.api.Member.Reports.GetTimeReportsBetweenDates(userID, start, end, api.FetchRelated())
	if err != nil {
		return nil, err
	}
	return MonthlyMemberTimeReports(reports), nil
}

// FiscalYearMemberTimeReports queries Glass Factory and returns time reports for the given fiscal year
func (s *Service) FiscalYearMemberTimeReports(userID int, fiscalYear *FiscalYear) ([]*FiscalYearMemberTimeReport, error) {
	start := fiscalYear.Start
	end := fiscalYear.End
	reports, err := s.api.Member.Reports.GetTimeReportsBetweenDates(userID, start, end, api.FetchRelated())
	if err != nil {
		return nil, err
	}
	return FiscalYearMemberTimeReports(reports, fiscalYear.End.Month()), nil
}
