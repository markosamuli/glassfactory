package api

import (
	"cloud.google.com/go/civil"
	"errors"
	"fmt"
	"github.com/markosamuli/glassfactory/dateutils"
	"github.com/markosamuli/glassfactory/models"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func NewReportsService(s *Service) *ReportsService {
	rs := &ReportsService{s: s}
	return rs
}

type ReportsService struct {
	s *Service
}

//func (r *ReportsService) GetTimeReportsForCurrentMonth(userID int) ([]*MemberTimeReport, error) {
//	return r.GetTimeReportsForMonth(userID, time.Now())
//}
//
//func (r *ReportsService) GetTimeReportsForCurrentYear(userID int) ([]*MemberTimeReport, error) {
//	return r.GetTimeReportsForYear(userID, time.Now())
//}


//func (r *ReportsService) GetTimeReportsForMonth(userID int, t time.Time) ([]*MemberTimeReport, error) {
//	response, err := r.MemberTimeReportForMonth(userID, t).Do()
//	if err != nil {
//		return nil, err
//	}
//	return response.Reports, nil
//}
//
//func (r *ReportsService) GetTimeReportsForYear(userID int, t time.Time) ([]*MemberTimeReport, error) {
//	responses, err := r.MemberTimeReportsForYear(userID, t).Do()
//	if err != nil {
//		return nil, err
//	}
//	reports := make([]*MemberTimeReport, 0)
//	for _, response := range responses {
//		for _, report := range response.Reports {
//			reports = append(reports, report)
//		}
//	}
//	return reports, nil
//}

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

type MemberTimeReport struct {
	UserID int `json:"user_id"`
	Date civil.Date `json:"date"`
	Planned float32 `json:"planned"`
	Actual float32 `json:"time"`
	ClientID int `json:"client_id,omitempty"`
	ProjectID int `json:"project_id,omitempty"`
	JobID string `json:"job_id,omitempty"`
	ActivityID int `json:"activity_id,omitempty"`
	RoleID int `json:"role_id,omitempty"`
	Client *models.Client
	Project *models.Project
}

func (r *MemberTimeReport) CalendarMonth() dateutils.CalendarMonth {
	return dateutils.CalendarMonth{
		Year:  r.Date.Year,
		Month: r.Date.Month,
	}
}

type MemberTimeReportResponse struct {
	Reports []*MemberTimeReport
}

type MemberTimeReportCalls struct {
	s *Service
	userID int
	start civil.Date
	end civil.Date
	calls []*MemberTimeReportCall
}

func (m *MemberTimeReportCalls) Append(c *MemberTimeReportCall) {
	if !m.start.IsValid() || c.start.Before(m.start) {
		m.start = c.start
	}
	if !m.end.IsValid() || c.end.After(m.end) {
		m.end = c.end
	}
	m.calls = append(m.calls, c)
}

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

func (c *MemberTimeReportCall) Do() (*MemberTimeReportResponse, error) {
	res, err := c.doRequest()
	if err != nil {
		return nil, err
	}
	target := make([]*MemberTimeReport, 0)
	if err := DecodeResponse(&target, res); err != nil {
		return nil, err
	}
	ret := &MemberTimeReportResponse{}
	ret.Reports = target
	return ret, nil
}

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

func (r *ReportsService) MemberTimeReportForMonth(userID int, t time.Time) *MemberTimeReportCall {
	start := dateutils.BeginningOfMonth(t)
	end := dateutils.EndOfMonth(t)
	return r.MemberTimeReport(userID, start, end)
}

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

