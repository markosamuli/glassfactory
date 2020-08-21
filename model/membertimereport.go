package model

import (
	"github.com/markosamuli/glassfactory/pkg/dateutil"
)

// MemberTimeReport represents a single team member time report entry
type MemberTimeReport struct {
	UserID     int           `json:"user_id"`
	Date       dateutil.Date `json:"date"`
	Planned    float64       `json:"planned"`
	Actual     float64       `json:"time"`
	ClientID   int           `json:"client_id,omitempty"`
	ProjectID  int           `json:"project_id,omitempty"`
	JobID      string        `json:"job_id,omitempty"`
	ActivityID int           `json:"activity_id,omitempty"`
	RoleID     int           `json:"role_id,omitempty"`
	Client     *Client
	Project    *Project
}
