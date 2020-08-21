package model

import "time"

// Client represents client details in Glass Factory
type Client struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	ArchivedAt time.Time `json:"archived_at,omitempty"`
	OwnerID    int       `json:"owner_id"`
	OfficeID   int       `json:"office_id"`
}

// IsArchived returns true if archive date has been defined
func (c *Client) IsArchived() bool {
	return !c.ArchivedAt.IsZero()
}