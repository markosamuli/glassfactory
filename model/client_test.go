package model

import (
	"encoding/json"
	"testing"
	"time"

	"gotest.tools/assert"
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
	assert.NilError(t, err)

	assert.Equal(t, target.ID, 1234)
	assert.Equal(t, target.Name, "Google")
	assert.Assert(t, target.ArchivedAt.IsZero())
	assert.Assert(t, !target.IsArchived())
	assert.Equal(t, target.OwnerID, 567)
	assert.Equal(t, target.OfficeID, 789)
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
	assert.NilError(t, err)

	archivedAt := time.Date(2018, time.Month(6), 7, 7, 27, 54, 563 * 1000000, time.UTC)
	assert.Equal(t, target.ArchivedAt, archivedAt)
	assert.Assert(t, target.IsArchived())
}