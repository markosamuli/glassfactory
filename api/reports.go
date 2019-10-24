package api

import (
	"cloud.google.com/go/civil"
	"errors"
	"fmt"
	"github.com/markosamuli/glassfactory/dateutils"
	"github.com/markosamuli/glassfactory/models"
	"github.com/markosamuli/glassfactory/reports"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// NewReportsService creates a new ReportsService
func NewReportsService(s *Service) *ReportsService {
	rs := &ReportsService{s: s}
	return rs
}

// ReportsService provides methods for fetching time report data from Glass Factory
type ReportsService struct {
	s *Service
}

// MonthlyMemberTimeReports queries Glass Factory and returns time reports for a full calendar year matching the given time
func (r *ReportsService) MonthlyMemberTimeReports(userID int, t time.Time) ([]*reports.MonthlyMemberTimeReport, error) {
	responses, err := r.MemberTimeReportsForYear(userID, t).Do()
	if err != nil {
		return nil, err
	}
	mtr := make([]*models.MemberTimeReport, 0)
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
			mtr = append(mtr, report)
		}
	}
	return reports.MonthlyMemberTimeReports(mtr), nil
}

// MonthlyMemberTimeReports queries Glass Factory and returns time reports for the given fiscal year
func (r *ReportsService) FiscalYearMemberTimeReports(userID int, fiscalYear *dateutils.FiscalYear) ([]*reports.FiscalYearMemberTimeReport, error) {
	responses, err := r.MemberTimeReportsForFiscalYear(userID, fiscalYear).Do()
	if err != nil {
		return nil, err
	}
	mtr := make([]*models.MemberTimeReport, 0)
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
			mtr = append(mtr, report)
		}
	}
	return reports.FiscalYearMemberTimeReports(mtr, fiscalYear.End.Month()), nil
}

// MemberTimeReport queries Glass Factory and returns time reports for the given time period
func (r *ReportsService) MemberTimeReport(userID int, start time.Time, end time.Time) *MemberTimeReportCall {
	today := time.Now()
	if today.Before(end) {
		end = today
	}
	c := &MemberTimeReportCall{s: r.s}
	c.userID = userID
	c.start = civil.DateOf(start)
	c.end = civil.DateOf(end)
	return c
}

// MemberTimeReportsForFiscalYear queries Glass Factory and returns time reports for the given fiscal year
func (r *ReportsService) MemberTimeReportsForFiscalYear(userID int, fiscalYear *dateutils.FiscalYear) *MemberTimeReportCalls {
	now := time.Now()
	start := fiscalYear.Start
	end := fiscalYear.End
	if end.After(now) {
		end = now
	}
	calls := &MemberTimeReportCalls{s: r.s}
	calls.userID = userID
	for _, m := range dateutils.MonthsBetweenDates(start, end) {
		c := r.MemberTimeReport(userID, m.Start, m.End)
		calls.Append(c)
	}
	return calls
}

// MemberTimeReportsForYear creates multiple calls for fetching time reports for a full calendar year matching the given date
func (r *ReportsService) MemberTimeReportsForYear(userID int, t time.Time) *MemberTimeReportCalls {
	now := time.Now()
	start := dateutils.BeginningOfYear(t)
	end := dateutils.EndOfYear(t)
	if end.After(now) {
		end = now
	}
	calls := &MemberTimeReportCalls{s: r.s}
	calls.userID = userID
	for _, m := range dateutils.MonthsBetweenDates(start, end) {
		c := r.MemberTimeReport(userID, m.Start, m.End)
		calls.Append(c)
	}
	return calls
}

// MemberTimeReportForMonth creates a call for fetching time reports for a single month
func (r *ReportsService) MemberTimeReportForMonth(userID int, t time.Time) *MemberTimeReportCall {
	start := dateutils.BeginningOfMonth(t)
	end := dateutils.EndOfMonth(t)
	return r.MemberTimeReport(userID, start, end)
}

// MemberTimeReportCall is used for fetching time reports for given time period from Glass Factory
type MemberTimeReportCall struct {
	s *Service
	userID int // User ID
	clientID int // Client ID
	projectIDs []int // Project IDs separated by comma
	officeID int // Office ID
	start civil.Date // Range start date
	end civil.Date // Range end date
	date civil.Date // Date of report data for single date. If date is set, start and end params will be ignored.
}

// MemberTimeReportResponse contains the time reports returned from Glass Factory
type MemberTimeReportResponse struct {
	Reports []*models.MemberTimeReport
}

func (c *MemberTimeReportCall) doRequest() (*http.Response, error) {
	var urls string
	if c.userID > 0 {
		urls = c.s.BasePath + fmt.Sprintf("members/%d/reports/time.json", c.userID)
	} else {
		return nil, errors.New("user ID is required")
	}

	urlParams := url.Values{}

	// Mandatory date range for the report
	switch {
	case c.date.IsValid():
		urlParams.Add("date", c.start.String())
	case c.start.IsValid() && c.end.IsValid():
		urlParams.Add("start", c.start.String())
		urlParams.Add("end", c.end.String())
	default:
		return nil, errors.New("start and end or date parameter is required")
	}

	// Optional parameters
	if len(c.projectIDs) > 0 {
		var projectIDs []string
		for i := range c.projectIDs {
			projectIDs = append(projectIDs, strconv.Itoa(c.projectIDs[i]))
		}
		urlParams.Add("project_id", strings.Join(projectIDs, ","))
	}
	if c.clientID > 0 {
		urlParams.Add("client_id", strconv.Itoa(c.clientID))
	}
	if c.clientID > 0 {
		urlParams.Add("office_id", strconv.Itoa(c.clientID))
	}

	urls += "?" + urlParams.Encode()
	req, err := http.NewRequest(http.MethodGet, urls, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Do executes the MemberTimeReportCall and returns MemberTimeReportResponse with the time report data
func (c *MemberTimeReportCall) Do() (*MemberTimeReportResponse, error) {
	res, err := c.doRequest()
	if err != nil {
		return nil, err
	}
	target := make([]*models.MemberTimeReport, 0)
	if err := DecodeResponse(&target, res); err != nil {
		return nil, err
	}
	ret := &MemberTimeReportResponse{}
	ret.Reports = target
	return ret, nil
}

// MemberTimeReportCall is used for fetching multiple time reports from Glass Factory
type MemberTimeReportCalls struct {
	s *Service
	userID int
	start civil.Date
	end civil.Date
	calls []*MemberTimeReportCall
}

// Append additional MemberTimeReportCall to the list of calls
func (m *MemberTimeReportCalls) Append(c *MemberTimeReportCall) {
	if !m.start.IsValid() || c.start.Before(m.start) {
		m.start = c.start
	}
	if !m.end.IsValid() || c.end.After(m.end) {
		m.end = c.end
	}
	m.calls = append(m.calls, c)
}

// Do executes all the queries in MemberTimeReportCalls
func (m *MemberTimeReportCalls) Do() ([]*MemberTimeReportResponse, error) {
	responses := make([]*MemberTimeReportResponse, len(m.calls))
	for i, c := range m.calls {
		res, err := c.Do()
		if err != nil {
			return nil, err
		}
		responses[i] = res
	}
	return responses, nil
}