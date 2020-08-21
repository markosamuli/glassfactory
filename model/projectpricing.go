package model

// ProjectPricing contains pricing data for a project
type ProjectPricing struct {
	Strategy  string  `json:"strategy"`
	RatesType string  `json:"rates_type"`
	Budget    float64 `json:"budget"`
	Currency  string  `json:"currency"`
}
