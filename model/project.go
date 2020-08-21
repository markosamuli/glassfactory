package model

import (
	"github.com/markosamuli/glassfactory/pkg/dateutil"
)

// Project represents an internal or client project in Glass Factory
type Project struct {
	ID                int               `json:"id"`
	Name              string            `json:"name"`
	Archived          bool              `json:"archived,omitempty"`
	ManagerID         int               `json:"manager_id"`
	URL               string            `json:"url,omitempty"`
	OfficeID          int               `json:"office_id"`
	ClientID          int               `json:"client_id"`
	JobID             string            `json:"job_id"`
	CreatedAt         dateutil.DateTime `json:"created_at,omitempty"`
	ClosedAt          dateutil.Date     `json:"closed_at,omitempty"`
	BillableStatus    BillableStatus    `json:"billable_status"`
	Probability       float64           `json:"probability"`
	Closed            bool              `json:"closed,omitempty"`
	ProbabilityState  string            `json:"probability_state"`
	EnableGeneralTime bool              `json:"enable_general_time,omitempty"`
	Pricing           *ProjectPricing   `json:"pricing"`
}
