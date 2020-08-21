package model

import (
	"github.com/markosamuli/glassfactory/pkg/dateutil"
)

// Member represents staff member details in Glass Factory
type Member struct {
	ID         int           `json:"id"`
	Name       string        `json:"name"`
	Email      string        `json:"email"`
	JoinedAt   dateutil.Date `json:"joined_at,omitempty"`
	ArchivedAt dateutil.Date `json:"archived_at,omitempty"`
	Freelancer bool          `json:"freelancer"`
	RoleID     int           `json:"role_id"`
	Capacity   float64       `json:"capacity"`
	Archived   bool          `json:"archived"`
	OfficeID   int           `json:"office_id"`
	Avatar     *MemberAvatar `json:"avatar"`
}
