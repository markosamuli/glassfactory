package reporting

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/jinzhu/now"
	"github.com/markosamuli/glassfactory/api"
	"github.com/markosamuli/glassfactory/model"
	"github.com/markosamuli/glassfactory/pkg/dateutil"
	"gopkg.in/h2non/gock.v1"
	"gotest.tools/assert"
)

func newTestSettings() *api.Settings {
	settings := &api.Settings{}
	settings.UserEmail = "test@example.com"
	settings.UserToken = "abcdefg1234"
	settings.AccountSubdomain = "example"
	return settings
}

func TestNewService(t *testing.T) {
	var ctx context.Context
	var apiService *api.Service
	var s *Service
	var ns *Service
	var err error

	ctx = context.Background()

	s, err = NewService(ctx, apiService)
	assert.ErrorContains(t, err, "apiService is nil")
	assert.Equal(t, s, ns)

	apiService = &api.Service{}
	s, err = NewService(ctx, apiService)
	assert.NilError(t, err)
	assert.Equal(t, s.api, apiService)
}

func TestService_MonthlyMemberTimeReports(t *testing.T) {
	defer gock.Off()

	var body []byte
	var err error

	domain := "https://example.glassfactory.io"
	apiPath := "/api/public/v1/"

	today := time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
	userID := 123
	clientID := 111
	projectID := 222

	client := model.Client{ID: clientID, Name: "Test Client"}
	body, err = json.Marshal(client)
	assert.NilError(t, err)
	gock.New(domain).
		Get(apiPath + fmt.Sprintf("clients/%d.json", clientID)).
		Reply(200).
		BodyString(string(body))

	project := model.Project{ID: projectID, Name: "Test Project", BillableStatus: model.Billable}
	body, err = json.Marshal(project)
	assert.NilError(t, err)
	gock.New(domain).
		Get(apiPath + fmt.Sprintf("projects/%d.json", projectID)).
		Reply(200).
		BodyString(string(body))

	var expected [][]model.MemberTimeReport
	for i := 0; i < 12; i++ {
		start := today.AddDate(0, i, 0)
		end := now.With(start).EndOfMonth()
		data := []model.MemberTimeReport{
			{
				UserID:    userID,
				ClientID:  clientID,
				ProjectID: projectID,
				Date:      dateutil.DateOf(start),
				Planned:   8.0,
				Actual:    7.5,
			},
		}
		expected = append(expected, data)
		body, err = json.Marshal(data)
		assert.NilError(t, err)
		gock.New(domain).
			Get(apiPath+fmt.Sprintf("members/%d/reports/time.json", userID)).
			MatchParam("start", start.Format("2006-01-02")).
			MatchParam("end", end.Format("2006-01-02")).
			Reply(200).
			BodyString(string(body))
	}

	var ctx context.Context
	var apiService *api.Service
	var settings *api.Settings
	var s *Service
	var reports []*MonthlyMemberTimeReport

	ctx = context.Background()
	settings = newTestSettings()

	apiService, err = api.NewService(ctx, settings)
	assert.NilError(t, err)

	s, err = NewService(ctx, apiService)
	assert.NilError(t, err)

	reports, err = s.MonthlyMemberTimeReports(userID, today)
	assert.NilError(t, err)
	assert.Equal(t, len(reports), 12)

	for i, r := range reports {
		assert.Equal(t, r.UserID, userID)

		assert.Equal(t, r.CalendarMonth.Year, today.Year())
		assert.Equal(t, int(r.CalendarMonth.Month), int(today.Month())+i)

		assert.Assert(t, r.Start == expected[i][0].Date || r.Start.Before(expected[i][0].Date))
		assert.Assert(t, r.End == expected[i][0].Date || r.End.After(expected[i][0].Date))

		assert.Equal(t, len(r.Reports), len(expected[i]))

		assert.Equal(t, r.Planned(), float64(len(expected[i]))*8.0)
		assert.Equal(t, r.Actual(), float64(len(expected[i]))*7.5)

		assert.Equal(t, r.Reports[0].Client.ID, client.ID)
		assert.Equal(t, r.Reports[0].Project.ID, project.ID)
	}

}

func TestService_FiscalYearMemberTimeReports(t *testing.T) {
	defer gock.Off()

	var body []byte
	var err error

	domain := "https://example.glassfactory.io"
	apiPath := "/api/public/v1/"

	today := time.Date(2017, time.January, 1, 0, 0, 0, 0, time.Local)
	userID := 123
	clientID := 111
	projectID := 222
	fiscalYear := NewFiscalYear(today, time.January)

	client := model.Client{ID: clientID, Name: "Test Client"}
	body, err = json.Marshal(client)
	assert.NilError(t, err)
	gock.New(domain).
		Get(apiPath + fmt.Sprintf("clients/%d.json", clientID)).
		Reply(200).
		BodyString(string(body))

	project := model.Client{ID: projectID, Name: "Test Project"}
	body, err = json.Marshal(project)
	assert.NilError(t, err)
	gock.New(domain).
		Get(apiPath + fmt.Sprintf("projects/%d.json", projectID)).
		Reply(200).
		BodyString(string(body))

	var expected []model.MemberTimeReport
	for i := 0; i < 12; i++ {
		start := now.With(fiscalYear.Start.AddDate(0, i, 0)).BeginningOfMonth()
		end := now.With(start).EndOfMonth()
		data := []model.MemberTimeReport{
			{
				UserID:    userID,
				ClientID:  clientID,
				ProjectID: projectID,
				Date:      dateutil.DateOf(start),
				Planned:   8.0,
				Actual:    7.5,
			},
		}
		for _, r := range data {
			expected = append(expected, r)
		}
		body, err = json.Marshal(data)
		assert.NilError(t, err)
		gock.New(domain).
			Get(apiPath+fmt.Sprintf("members/%d/reports/time.json", userID)).
			MatchParam("start", start.Format("2006-01-02")).
			MatchParam("end", end.Format("2006-01-02")).
			Reply(200).
			BodyString(string(body))
	}

	var ctx context.Context
	var apiService *api.Service
	var settings *api.Settings
	var s *Service
	var reports []*FiscalYearMemberTimeReport

	ctx = context.Background()
	settings = newTestSettings()

	apiService, err = api.NewService(ctx, settings)
	assert.NilError(t, err)

	s, err = NewService(ctx, apiService)
	assert.NilError(t, err)

	reports, err = s.FiscalYearMemberTimeReports(userID, fiscalYear)
	assert.NilError(t, err)
	assert.Equal(t, len(reports), 1)

	for _, r := range reports {
		assert.Equal(t, r.UserID, userID)
		assert.Equal(t, r.FiscalYear.String(), fiscalYear.String())
		assert.Equal(t, r.FiscalYear.End, fiscalYear.End)

		assert.Assert(t, r.Start == expected[0].Date || r.Start.Before(expected[0].Date))
		assert.Assert(t, r.End == expected[0].Date || r.End.After(expected[0].Date))

		assert.Equal(t, len(r.Reports), len(expected))

		assert.Equal(t, r.Planned(), float64(len(expected))*8.0)
		assert.Equal(t, r.Actual(), float64(len(expected))*7.5)

		assert.Equal(t, r.Reports[0].Client.ID, client.ID)
		assert.Equal(t, r.Reports[0].Project.ID, project.ID)
	}
}
