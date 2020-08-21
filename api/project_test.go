package api

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/markosamuli/glassfactory/model"
	"gopkg.in/h2non/gock.v1"
	"gotest.tools/assert"
)

func TestGetProject(t *testing.T) {
	defer gock.Off()

	domain := "https://example.glassfactory.io"
	apiPath := "/api/public/v1/"
	endpoint := domain + apiPath
	projectID := 12345

	gock.New(domain).
		Get(apiPath + fmt.Sprintf("projects/%d.json", projectID)).
		Reply(200).
		BodyString(`{
		  "id": 12345,
		  "name": "Luna Lights Phase 1",
		  "archived": false,
		  "manager_id": 423,
		  "url": "http://example.glassfactory.io/projects/14052",
		  "office_id": 103,
		  "client_id": 1974,
		  "job_id": "J12345",
		  "created_at": "2018-06-04 05:46:15",
		  "closed_at": null,
		  "billable_status": "billable",
		  "probability": 100,
		  "closed": false,
		  "probability_state": "live",
		  "enable_general_time": true,
		  "pricing": {
			"strategy": "material",
			"rates_type": "client",
			"budget": 10000,
			"currency": "USD"
		  }
		}`)

	s := &Service{}
	s.client = &http.Client{}
	s.BasePath = endpoint

	// Project without data
	invalid := &model.Project{
		ID: projectID,
	}

	var project *model.Project
	var err error
	rs := NewProjectService(s)

	// Add project without data into the cache
	rs.projects.Add(invalid)

	// Get project without using cache
	project, err = rs.Get(projectID, WithCache(false))
	assert.NilError(t, err)

	assert.Equal(t, project.ID, projectID)
	assert.Equal(t, project.Name, "Luna Lights Phase 1")
	assert.Assert(t, !project.Archived)
	assert.Equal(t, project.ManagerID, 423)
	assert.Equal(t, project.URL, "http://example.glassfactory.io/projects/14052")
	assert.Equal(t, project.OfficeID, 103)
	assert.Equal(t, project.ClientID, 1974)
	assert.Equal(t, project.JobID, "J12345")
	assert.Equal(t, project.CreatedAt.In(time.UTC), time.Date(2018, time.June, 4, 5, 46, 15, 0, time.UTC))
	assert.Assert(t, !project.ClosedAt.IsValid())
	assert.Equal(t, project.BillableStatus, model.Billable)
	assert.Equal(t, project.Probability, 100.0)
	assert.Equal(t, project.Closed, false)
	assert.Equal(t, project.ProbabilityState, "live")
	assert.Equal(t, project.EnableGeneralTime, true)
	assert.Equal(t, project.Pricing.Budget, 10000.0)
	assert.Equal(t, project.Pricing.RatesType, "client")
	assert.Equal(t, project.Pricing.Strategy, "material")
	assert.Equal(t, project.Pricing.Currency, "USD")

	// Verify that we don't have pending mocks
	assert.Assert(t, gock.IsDone(), "all mocks should have been called")

}

func TestGetProjectCached(t *testing.T) {
	defer gock.Off()

	domain := "https://example.glassfactory.io"
	apiPath := "/api/public/v1/"
	endpoint := domain + apiPath
	projectID := 678

	gock.New(domain).
		Get(apiPath + fmt.Sprintf("projects/%d.json", projectID)).
		Reply(200).
		BodyString(`{
		  "id": 678
		}`)

	s := &Service{}
	s.client = &http.Client{}
	s.BasePath = endpoint

	// Mock project
	expected := &model.Project{
		ID:   projectID,
		Name: "Cached Project",
	}

	var project *model.Project
	var err error
	rs := NewProjectService(s)

	// Add project into the service cache
	rs.projects.Add(expected)

	// Get project with caching enabled
	project, err = rs.Get(projectID, WithCache(true))
	assert.NilError(t, err)

	assert.Equal(t, project.ID, projectID)
	assert.Equal(t, project, expected)

	assert.Equal(t, len(gock.Pending()), 1, "mock request shouldn't be called")
}

func TestProjectServiceAll(t *testing.T) {
	defer gock.Off()

	domain := "https://example.glassfactory.io"
	apiPath := "/api/public/v1/"
	endpoint := domain + apiPath

	gock.New(domain).
		Get(apiPath + "projects.json").
		Reply(200).
		BodyString(`[
			{
			  "id": 12345,
			  "name": "Luna Lights Phase 1",
			  "archived": false,
			  "manager_id": 423,
			  "url": "http://example.glassfactory.io/projects/14052",
			  "office_id": 103,
			  "client_id": 1974,
			  "job_id": "J12345",
			  "created_at": "2018-06-04 05:46:15",
			  "closed_at": null,
			  "billable_status": "billable",
			  "probability": 100,
			  "closed": false,
			  "probability_state": "live",
			  "enable_general_time": true,
			  "pricing": {
				"strategy": "material",
				"rates_type": "client",
				"budget": 10000,
				"currency": "USD"
			  }
			}
		]`)

	s := &Service{}
	s.client = &http.Client{}
	s.BasePath = endpoint

	var projects []*model.Project
	var err error
	rs := NewProjectService(s)

	projects, err = rs.All()
	assert.NilError(t, err)
	assert.Equal(t, len(projects), 1)

	assert.Equal(t, projects[0].ID, 12345)
	assert.Equal(t, projects[0].Name, "Luna Lights Phase 1")

	// Verify that we don't have pending mocks
	assert.Assert(t, gock.IsDone(), "all mocks should have been called")
}
