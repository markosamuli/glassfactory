package model

import (
	"encoding/json"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestProject(t *testing.T) {
	jsonString := `{
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
	}`
	var project Project
	err := json.Unmarshal([]byte(jsonString), &project)
	assert.NilError(t, err)

	assert.Equal(t, project.ID, 12345)
	assert.Equal(t, project.Name, "Luna Lights Phase 1")
	assert.Assert(t, !project.Archived)
	assert.Equal(t, project.ManagerID, 423)
	assert.Equal(t, project.URL, "http://example.glassfactory.io/projects/14052")
	assert.Equal(t, project.OfficeID, 103)
	assert.Equal(t, project.ClientID, 1974)
	assert.Equal(t, project.JobID, "J12345")

	dt := time.Date(2018, time.June, 4, 5, 46, 15, 0, time.UTC)
	assert.Equal(t, project.CreatedAt.In(time.UTC), dt)
	assert.Assert(t, !project.ClosedAt.IsValid())

	assert.Equal(t, project.BillableStatus, Billable)
	assert.Equal(t, project.Probability, 100.0)
	assert.Equal(t, project.Closed, false)
	assert.Equal(t, project.ProbabilityState, "live")
	assert.Equal(t, project.EnableGeneralTime, true)

	assert.Equal(t, project.Pricing.Budget, 10000.0)
	assert.Equal(t, project.Pricing.RatesType, "client")
	assert.Equal(t, project.Pricing.Strategy, "material")
	assert.Equal(t, project.Pricing.Currency, "USD")
}
