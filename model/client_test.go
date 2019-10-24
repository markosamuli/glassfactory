package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	jsonString := `
	{
		"id": 1234,
		"name": "Google",
		"archived_at": null,
		"owner_id": 567,
		"office_id": 789
	}`
	var target Client
	err := json.Unmarshal([]byte(jsonString), &target)
	if err != nil {
		t.Errorf("Failed to decode JSON, got error %v", err)
	}
	if target.ID != 1234 {
		t.Errorf("ID = %d; want 1234", target.ID)
	}
	if target.Name != "Google" {
		t.Errorf("Name = %s; want Google", target.Name)
	}
	if !target.ArchivedAt.IsZero() {
		t.Errorf("ArchivedAt = %s; want zero date", target.ArchivedAt)
	}
	if target.IsArchived() {
		t.Errorf("IsArchived() returned true; want false")
	}
	if target.OwnerID != 567 {
		t.Errorf("OwnerID = %d; want 567", target.OwnerID)
	}
	if target.OfficeID != 789 {
		t.Errorf("OfficeID = %d; want 789", target.OfficeID)
	}
}

func TestClientArchived(t *testing.T) {
	jsonString := `
	{
		"id": 1234,
		"name": "Google",
		"archived_at": "2018-06-07T07:27:54.563Z",
		"owner_id": 567,
		"office_id": 789
	}`
	var target Client
	err := json.Unmarshal([]byte(jsonString), &target)
	if err != nil {
		t.Errorf("Failed to decode JSON, got error %v", err)
	}
	archivedAt := time.Date(2018, time.Month(6), 7, 7, 27, 54, 563 * 1000000, time.UTC)
	if target.ArchivedAt != archivedAt {
		t.Errorf("ArchivedAt = %s; want %s", target.ArchivedAt, archivedAt)
	}
	if !target.IsArchived() {
		t.Errorf("IsArchived() returned false; want true")
	}
}