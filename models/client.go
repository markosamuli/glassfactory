package models

type Client struct {
	ID int `json:"id"`
	Name string `json:"name"`
	ArchivedAt string `json:"archived_at,omitempty"`
	OwnerID int `json:"owner_id"`
	OfficeID int `json:"office_id"`
}
