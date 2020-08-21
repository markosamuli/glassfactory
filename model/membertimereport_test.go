package model

import (
	"encoding/json"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestMemberTimeReport(t *testing.T) {
	jsonString := `{
		"client_id": 2079,
		"project_id": 14330,
		"job_id": null,
		"activity_id": null,
		"user_id": 3,
		"role_id": 1480,
		"date": "2018-07-10",
		"planned": 8,
		"time": 0
	}`
	var report MemberTimeReport
	err := json.Unmarshal([]byte(jsonString), &report)
	assert.NilError(t, err)

	assert.Equal(t, report.ClientID, 2079)
	assert.Equal(t, report.ProjectID, 14330)
	assert.Equal(t, report.JobID, "")
	assert.Equal(t, report.ActivityID, 0)
	assert.Equal(t, report.UserID, 3)
	assert.Equal(t, report.RoleID, 1480)

	dt := time.Date(2018, time.July, 10, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, report.Date.In(time.UTC), dt)
	assert.Assert(t, report.Date.IsValid())

	assert.Equal(t, report.Planned, 8.0)
	assert.Equal(t, report.Actual, 0.0)
}
