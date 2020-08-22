package api

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/101loops/clock"
	"github.com/markosamuli/glassfactory/model"
	"gopkg.in/h2non/gock.v1"
	"gotest.tools/assert"
)

func TestTimeReportOptions(t *testing.T) {
	defer gock.Off()

	domain := "https://example.glassfactory.io"
	apiPath := "/api/public/v1/"
	endpoint := domain + apiPath
	userID := 123

	gock.New(domain).
		Get(apiPath+fmt.Sprintf("members/%d/reports/time.json", userID)).
		MatchParam("start", "2019-09-01").
		MatchParam("end", "2019-09-30").
		MatchParam("client_id", "2079").
		Reply(200).
		BodyString(`[
		  {
			"client_id": 2079,
			"project_id": 123,
			"user_id": 123,
			"role_id": 1480,
			"date": "2019-09-01",
			"planned": 8,
			"time": 5.5
		  }
		]`)

	gock.New(domain).
		Get(apiPath+fmt.Sprintf("members/%d/reports/time.json", userID)).
		MatchParam("start", "2019-09-01").
		MatchParam("end", "2019-09-30").
		MatchParam("project_id", "14330").
		Reply(200).
		BodyString(`[
		  {
			"client_id": 123,
			"project_id": 14330,
			"user_id": 123,
			"role_id": 1480,
			"date": "2019-09-01",
			"planned": 8,
			"time": 5.5
		  }
		]`)

	gock.New(domain).
		Get(apiPath+fmt.Sprintf("members/%d/reports/time.json", userID)).
		MatchParam("start", "2019-09-01").
		MatchParam("end", "2019-09-30").
		MatchParam("project_id", "111,222").
		Reply(200).
		BodyString(`[
		  {
			"client_id": 123,
			"project_id": 111,
			"user_id": 123,
			"role_id": 1480,
			"date": "2019-09-01",
			"planned": 8,
			"time": 5.5
		  },
		  {
			"client_id": 123,
			"project_id": 222,
			"user_id": 123,
			"role_id": 1480,
			"date": "2019-09-01",
			"planned": 8,
			"time": 5.5
		  }
		]`)

	gock.New(domain).
		Get(apiPath+fmt.Sprintf("members/%d/reports/time.json", userID)).
		MatchParam("start", "2019-09-01").
		MatchParam("end", "2019-09-30").
		MatchParam("office_id", "123").
		Reply(200).
		BodyString(`[
		  {
			"client_id": 2079,
			"project_id": 123,
			"user_id": 123,
			"role_id": 1480,
			"date": "2019-09-01",
			"planned": 8,
			"time": 5.5
		  },
		  {
			"client_id": 123,
			"project_id": 14330,
			"user_id": 123,
			"role_id": 1480,
			"date": "2019-09-01",
			"planned": 8,
			"time": 5.5
		  }
		]`)

	s := &Service{}
	s.client = &http.Client{}
	s.BasePath = endpoint

	var res *MemberTimeReportResponse
	var reports []*model.MemberTimeReport
	var err error
	ms := NewMemberService(s)
	rs := NewMemberReportsService(ms)

	// Mock the current time
	today := time.Date(2019, time.October, 1, 0, 0, 0, 0, time.UTC)
	rs.clock = clock.NewMock().Set(today)

	// Reports from the past
	start := time.Date(2019, time.September, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2019, time.September, 30, 0, 0, 0, 0, time.UTC)

	res, err = rs.TimeReport(userID, start, end, WithClient(2079)).Do()
	assert.NilError(t, err)

	reports = res.Reports
	assert.Equal(t, len(reports), 1)
	for _, r := range reports {
		assert.Equal(t, r.ClientID, 2079)
	}

	res, err = rs.TimeReport(userID, start, end, WithProject(14330)).Do()
	assert.NilError(t, err)

	reports = res.Reports
	assert.Equal(t, len(reports), 1)
	for _, r := range reports {
		assert.Equal(t, r.ProjectID, 14330)
	}

	res, err = rs.TimeReport(userID, start, end, WithProjects([]int{111, 222})).Do()
	assert.NilError(t, err)

	reports = res.Reports
	assert.Equal(t, len(reports), 2)
	assert.Equal(t, reports[0].ProjectID, 111)
	assert.Equal(t, reports[1].ProjectID, 222)

	res, err = rs.TimeReport(userID, start, end, WithOffice(123)).Do()
	assert.NilError(t, err)

	reports = res.Reports
	assert.Equal(t, len(reports), 2)
	for _, r := range reports {
		assert.Equal(t, r.UserID, 123)
	}

	// Verify that we don't have pending mocks
	assert.Assert(t, gock.IsDone(), "all mocks should have been called")
}
