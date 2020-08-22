package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"github.com/101loops/clock"
	"github.com/markosamuli/glassfactory/model"
	"github.com/markosamuli/glassfactory/pkg/dateutil"
)

// NewMemberReportsService creates a new MemberReportsService
func NewMemberReportsService(m *MemberService) *MemberReportsService {
	rs := &MemberReportsService{m: m}
	rs.clock = clock.New()
	return rs
}

// MemberReportsService provides methods for fetching time report data from Glass Factory
type MemberReportsService struct {
	m     *MemberService
	clock clock.Clock
}

// GetTimeReportsBetweenDates returns Glass Factory member time reports between given dates
func (r *MemberReportsService) GetTimeReportsBetweenDates(userID int, start time.Time, end time.Time, opts ...TimeReportOption) ([]*model.MemberTimeReport, error) {
	responses, err := r.TimeReportsBetweenDates(userID, start, end, opts...).Do()
	if err != nil {
		return nil, err
	}

	options := NewTimeReportOptions(opts)

	reports := make([]*model.MemberTimeReport, 0)
	for _, response := range responses {
		for _, report := range response.Reports {
			// Fetch related data if FetchRelated() option was enabled
			if options.fetchRelated {
				client, err := r.m.s.Client.Get(report.ClientID)
				if err != nil {
					return nil, err
				}
				project, err := r.m.s.Project.Get(report.ProjectID)
				if err != nil {
					return nil, err
				}
				report.Client = client
				report.Project = project
			}
			reports = append(reports, report)
		}
	}
	return reports, nil
}

// TimeReport queries Glass Factory and returns member time reports for the given time period
func (r *MemberReportsService) TimeReport(userID int, start time.Time, end time.Time, opts ...TimeReportOption) *MemberTimeReportCall {
	today := time.Now()
	if start.After(today) {
		start = today // Make sure we're not getting reports from the future
	}
	if end.After(today) {
		end = today // Make sure we're not getting reports from the future
	}
	c := &MemberTimeReportCall{s: r.m.s}
	c.userID = userID
	c.start = civil.DateOf(start)
	c.end = civil.DateOf(end)
	c.options = opts
	return c
}

// TimeReportsBetweenDates creates MemberTimeReportCalls to be used for fetching member time reports between the given dates
func (r *MemberReportsService) TimeReportsBetweenDates(userID int, start time.Time, end time.Time, opts ...TimeReportOption) *MemberTimeReportCalls {
	today := time.Now()
	if start.After(today) {
		start = today // Make sure we're not getting reports from the future
	}
	if end.After(today) {
		end = today // Make sure we're not getting reports from the future
	}
	calls := &MemberTimeReportCalls{s: r.m.s}
	calls.userID = userID
	calls.options = opts
	// Split calls into months for better performance
	for _, m := range dateutil.MonthsBetweenDates(start, end) {
		// Make sure we're not getting reports outside the given dates
		if m.Start.Before(start) {
			m.Start = start
		}
		if m.End.After(end) {
			m.End = end
		}
		c := r.TimeReport(userID, m.Start, m.End, opts...)
		calls.Append(c)
	}
	return calls
}

// MemberTimeReportCall is used for fetching time reports for given time period from Glass Factory
type MemberTimeReportCall struct {
	s       *Service
	userID  int        // User ID
	start   civil.Date // Range start date
	end     civil.Date // Range end date
	date    civil.Date // Date of report data for single date. If date is set, start and end params will be ignored.
	options []TimeReportOption
}

// Options returns request options with defaults
func (c *MemberTimeReportCall) Options() TimeReportOptions {
	options := TimeReportOptions{}
	options.apply(c.options)
	return options
}

// MemberTimeReportResponse contains the time reports returned from Glass Factory
type MemberTimeReportResponse struct {
	Reports []*model.MemberTimeReport
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
	options := c.Options()
	if len(options.projectIDs) > 0 {
		var projectIDs []string
		for i := range options.projectIDs {
			projectIDs = append(projectIDs, strconv.Itoa(options.projectIDs[i]))
		}
		urlParams.Add("project_id", strings.Join(projectIDs, ","))
	}
	if options.clientID > 0 {
		urlParams.Add("client_id", strconv.Itoa(options.clientID))
	}
	if options.officeID > 0 {
		urlParams.Add("office_id", strconv.Itoa(options.officeID))
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
	target := make([]*model.MemberTimeReport, 0)
	if err := DecodeResponse(&target, res); err != nil {
		return nil, err
	}
	ret := &MemberTimeReportResponse{}
	ret.Reports = target
	return ret, nil
}

// MemberTimeReportCalls is used for fetching multiple time reports from Glass Factory
type MemberTimeReportCalls struct {
	s       *Service
	userID  int
	start   civil.Date
	end     civil.Date
	options []TimeReportOption
	calls   []*MemberTimeReportCall
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
