package model

import (
	"cloud.google.com/go/civil"
	"time"
)

const MemberStatusActive = "active"
const MemberStatusArchived = "archived"

type Member struct {
	ID         int           `json:"id"`
	Name       string        `json:"name"`
	Email      string        `json:"email"`
	JoinedAt   civil.Date 	 `json:"joined_at,omitempty"`
	ArchivedAt time.Time     `json:"archived_at,omitempty"`
	Freelancer bool          `json:"freelancer"`
	RoleID     int           `json:"role_id"`
	Capacity   float32       `json:"capacity"`
	Archived   bool          `json:"archived"`
	OfficeID   int           `json:"office_id"`
	Avatar     *MemberAvatar `json:"avatar"`
}

type MemberAvatar struct {
	URL string `json:"url"`
}