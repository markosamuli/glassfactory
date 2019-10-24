package models

import (
	"cloud.google.com/go/civil"
	"github.com/markosamuli/glassfactory/dateutils"
)

type MemberTimeReport struct {
	UserID     int        `json:"user_id"`
	Date       civil.Date `json:"date"`
	Planned    float32    `json:"planned"`
	Actual     float32    `json:"time"`
	ClientID   int        `json:"client_id,omitempty"`
	ProjectID  int        `json:"project_id,omitempty"`
	JobID      string     `json:"job_id,omitempty"`
	ActivityID int        `json:"activity_id,omitempty"`
	RoleID     int        `json:"role_id,omitempty"`
	Client     *Client
	Project    *Project
}

func (r *MemberTimeReport) CalendarMonth() dateutils.CalendarMonth {
	return dateutils.CalendarMonth{
		Year:  r.Date.Year,
		Month: r.Date.Month,
	}
}

