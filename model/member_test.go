package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMember(t *testing.T) {
	jsonString := `
	{
	  "id": 1401,
	  "user_id": 422,
	  "name": "First Last",
	  "email": "first.last@example.com",
	  "joined_at": "2016-02-17",
	  "archived_at": null,
	  "freelancer": false,
	  "role_id": 1019,
	  "capacity": 8,
	  "archived": false,
	  "office_id": 152
	}`
	var target Member
	err := json.Unmarshal([]byte(jsonString), &target)
	if err != nil {
		t.Errorf("Failed to decode JSON, got error %v", err)
	}
	if target.ID != 1401 {
		t.Errorf("ID = %d; want 1401", target.ID)
	}
	if target.Name != "First Last" {
		t.Errorf("Name = %s; want First Last", target.Name)
	}
	if target.Email != "first.last@example.com" {
		t.Errorf("Email = %s; want first.last@example.com", target.Email)
	}
	if !target.JoinedAt.IsValid() {
		t.Errorf("JoinedAt = %s; want valid date", target.JoinedAt)
	}
	if !target.ArchivedAt.IsZero() {
		t.Errorf("ArchivedAt = %s; want zero date", target.ArchivedAt)
	}
	if target.Freelancer {
		t.Errorf("Freelancer = true; want false")
	}
	if target.Archived {
		t.Errorf("Archived = true; want false")
	}
	if target.Capacity != 8.0 {
		t.Errorf("Capacity = %f; want 8.0", target.Capacity)
	}
	if target.RoleID != 1019 {
		t.Errorf("RoleID = %d; want 1019", target.RoleID)
	}
	if target.OfficeID != 152 {
		t.Errorf("OfficeID = %d; want 152", target.OfficeID)
	}
}

func TestMemberArchived(t *testing.T) {
	jsonString := `
	{
	  "id": 1401,
	  "user_id": 422,
	  "name": "First Last",
	  "email": "first.last@example.com",
	  "joined_at": "2016-02-17",
      "archived_at": "2018-06-07T07:27:54.563Z",
	  "freelancer": false,
	  "role_id": 1019,
	  "capacity": 8,
	  "archived": true,
	  "office_id": 152
	}`
	var target Member
	err := json.Unmarshal([]byte(jsonString), &target)
	if err != nil {
		t.Errorf("Failed to decode JSON, got error %v", err)
	}
	archivedAt := time.Date(2018, time.Month(6), 7, 7, 27, 54, 563 * 1000000, time.UTC)
	if target.ArchivedAt != archivedAt {
		t.Errorf("ArchivedAt = %s; want %s", target.ArchivedAt, archivedAt)
	}
	if !target.Archived {
		t.Errorf("Archived = false; want true")
	}
}

func TestMemberFreelancer(t *testing.T) {
	jsonString := `
	{
	  "id": 1401,
	  "user_id": 422,
	  "name": "First Last",
	  "email": "first.last@example.com",
	  "joined_at": "2016-02-17",
      "archived_at": null,
	  "freelancer": true,
	  "role_id": 1019,
	  "capacity": 8,
	  "archived": false,
	  "office_id": 152
	}`
	var target Member
	err := json.Unmarshal([]byte(jsonString), &target)
	if err != nil {
		t.Errorf("Failed to decode JSON, got error %v", err)
	}
	if !target.Freelancer {
		t.Errorf("Freelancer = false; want true")
	}
}

func TestMemberCapacity(t *testing.T) {
	jsonString := `
	{
	  "id": 1401,
	  "user_id": 422,
	  "name": "First Last",
	  "email": "first.last@example.com",
	  "joined_at": "2016-02-17",
      "archived_at": null,
	  "freelancer": false,
	  "role_id": 1019,
	  "capacity": 7.5,
	  "archived": false,
	  "office_id": 152
	}`
	var target Member
	err := json.Unmarshal([]byte(jsonString), &target)
	if err != nil {
		t.Errorf("Failed to decode JSON, got error %v", err)
	}
	if target.Capacity != 7.5 {
		t.Errorf("Capacity = %f; want 7.5", target.Capacity)
	}
}