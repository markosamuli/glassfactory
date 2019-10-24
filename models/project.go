package models

type Project struct {
	ID                int             `json:"id"`
	Name              string          `json:"name"`
	Archived          bool            `json:"archived,omitempty"`
	ManagerID         int             `json:"manager_id"`
	URL               string          `json:"url,omitempty"`
	OfficeID          int             `json:"office_id"`
	ClientID          int             `json:"client_id"`
	JobID             string          `json:"job_id"`
	CreatedAt         string          `json:created_at`
	ClosedAt          string          `json:closed_at`
	BillableStatus    string          `json:"billable_status"`
	Probability       int             `json:"probability"`
	Closed            bool            `json:"closed,omitempty"`
	ProbabilityState  string          `json:"probability_state"`
	EnableGeneralTime bool            `json:"enable_general_time,omitempty"`
	Pricing           *ProjectPricing `json:"pricing"`
}

type ProjectPricing struct {
	Strategy  string  `json:strategy`
	RatesType string  `json:rates_type`
	Budget    float32 `json:budget`
	Currency  string  `json:currency`
}