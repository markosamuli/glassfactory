package reporting

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/markosamuli/glassfactory/dateutil"
	"github.com/markosamuli/glassfactory/model"
	"gotest.tools/assert"
)

func TestProjectMemberTimeReport(t *testing.T) {
	userID := 123
	clientID := 111
	projectID := 222

	client := &model.Client{ID: clientID, Name: "Test Client"}
	project := &model.Project{ID: projectID, Name: "Test Project", BillableStatus: model.Billable}
	report := NewProjectMemberTimeReport(userID, client, project)

	assert.Equal(t, report.UserID, userID)
	assert.Equal(t, report.Client, client)
	assert.Equal(t, report.Project, project)
	assert.Equal(t, len(report.Reports), 0)

	today := time.Date(2018, time.January, 1, 0,0,0, 0, time.UTC)
	var expected []*model.MemberTimeReport
	for i := 0; i < 10; i++ {
		d := today.AddDate(0, 0, 1)
		tr := &model.MemberTimeReport{
			UserID:    userID,
			Client:    client,
			Project:   project,
			ClientID:  client.ID,
			ProjectID: project.ID,
			Date:      dateutil.DateOf(d),
			Planned:   8.0,
			Actual:    7.5,
		}
		expected = append(expected, tr)
		report.Append(tr)
	}

	assert.Equal(t, len(report.Reports), len(expected))

	assert.Equal(t, report.Start, expected[0].Date)
	assert.Equal(t, report.End, expected[len(expected)-1].Date)

	assert.Equal(t, report.Planned(), float64(len(expected)) * 8.0)
	assert.Equal(t, report.Actual(), float64(len(expected)) * 7.5)
}

func TestProjectMemberTimeReports(t *testing.T) {

	today := time.Date(2018, time.January, 1, 0,0,0, 0, time.UTC)

	userID := 123

	var reports []*model.MemberTimeReport

	billableClient := &model.Client{ID: 200, Name: "Test Client with multiple projects"}
	for i := 0; i < 10; i++ {
		projectID := billableClient.ID + i
		project := &model.Project{
			ID: projectID,
			Name: fmt.Sprintf("Billable project %d", projectID),
			BillableStatus: model.Billable,
		}
		tr := &model.MemberTimeReport{
			UserID:    userID,
			Client:    billableClient,
			Project:   project,
			ClientID:  billableClient.ID,
			ProjectID: project.ID,
			Date:      dateutil.DateOf(today),
			Planned:   8.0,
			Actual:    7.5,
		}
		reports = append(reports, tr)
	}

	nonBillableClient := &model.Client{ID: 300, Name: "Test Client with multiple projects"}
	for i := 0; i < 10; i++ {
		projectID := nonBillableClient.ID + i
		project := &model.Project{
			ID: projectID,
			Name: fmt.Sprintf("Non billable project %d", projectID),
			BillableStatus: model.NonBillable,
		}
		tr := &model.MemberTimeReport{
			UserID:    userID,
			Client:    nonBillableClient,
			Project:   project,
			ClientID:  nonBillableClient.ID,
			ProjectID: project.ID,
			Date:      dateutil.DateOf(today),
			Planned:   1.0,
			Actual:    1.5,
		}
		reports = append(reports, tr)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(reports), func(i, j int) { reports[i], reports[j] = reports[j], reports[i] })

	projectReports := ProjectMemberTimeReports(reports)
	assert.Equal(t, len(projectReports), 20)
	for i := 0; i < 10; i++ {
		assert.Equal(t, projectReports[i].Client.ID, billableClient.ID)
		assert.Equal(t, projectReports[i].Planned(), 8.0)
		assert.Equal(t, projectReports[i].Actual(), 7.5)
	}
	for i := 0; i < 10; i++ {
		assert.Equal(t, projectReports[10+i].Client.ID, nonBillableClient.ID)
		assert.Equal(t, projectReports[10+i].Planned(), 1.0)
		assert.Equal(t, projectReports[10+i].Actual(), 1.5)
	}
}