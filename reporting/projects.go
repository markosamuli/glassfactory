package reporting

import (
	"sort"

	"github.com/markosamuli/glassfactory/model"
	"github.com/markosamuli/glassfactory/pkg/dateutil"
)

// ProjectMemberTimeReport represents time report data for a given project and team member
type ProjectMemberTimeReport struct {
	UserID  int
	Client  *model.Client
	Project *model.Project
	Start   dateutil.Date
	End     dateutil.Date
	Reports []*model.MemberTimeReport
}

// NewProjectMemberTimeReport creates a new ProjectMemberTimeReport for the given project and user
func NewProjectMemberTimeReport(userID int, client *model.Client, project *model.Project) *ProjectMemberTimeReport {
	return &ProjectMemberTimeReport{
		UserID:  userID,
		Client:  client,
		Project: project,
		Reports: make([]*model.MemberTimeReport, 0),
	}
}

// Append adds time report data to the report
func (tr *ProjectMemberTimeReport) Append(r *model.MemberTimeReport) {
	if !tr.Start.IsValid() || r.Date.Before(tr.Start) {
		tr.Start = r.Date
	}
	if !tr.End.IsValid() || r.Date.Before(tr.End) {
		tr.End = r.Date
	}
	tr.Reports = append(tr.Reports, r)
}

// Planned returns total planned hours
func (tr *ProjectMemberTimeReport) Planned() float64 {
	var planned float64
	for _, r := range tr.Reports {
		planned += r.Planned
	}
	return planned
}

// Actual returns total actual hours
func (tr *ProjectMemberTimeReport) Actual() float64 {
	var actual float64
	for _, r := range tr.Reports {
		actual += r.Actual
	}
	return actual
}

// ProjectMemberTimeReports converts MemberTimeReport to ProjectMemberTimeReport grouped by projects
func ProjectMemberTimeReports(reports []*model.MemberTimeReport) []*ProjectMemberTimeReport {
	projects := make(map[int]*ProjectMemberTimeReport, 0)
	for _, r := range reports {
		pr, ok := projects[r.Project.ID]
		if !ok {
			pr = NewProjectMemberTimeReport(r.UserID, r.Client, r.Project)
			projects[r.Project.ID] = pr
		}
		pr.Append(r)
	}
	pr := make([]*ProjectMemberTimeReport, 0, len(projects))
	for _, p := range projects {
		pr = append(pr, p)
	}
	sort.Sort(ByClient(pr))
	return pr
}

// ByClient implements sort.Interface based on the Client.ID field.
type ByClient []*ProjectMemberTimeReport

func (a ByClient) Len() int           { return len(a) }
func (a ByClient) Less(i, j int) bool { return a[i].Client.ID < a[j].Client.ID }
func (a ByClient) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
