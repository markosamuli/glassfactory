package models

import "time"

type Client struct {
	ID int `json:"id"`
	Name string `json:"name"`
	ArchivedAt time.Time `json:"archived_at,omitempty"`
	OwnerID int `json:"owner_id"`
	OfficeID int `json:"office_id"`
}

func (c *Client) IsArchived() bool {
	return !c.ArchivedAt.IsZero()
}