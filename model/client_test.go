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
		"name": "ACME Inc.",
		"archived_at": null,
		"owner_id": 567,
		"office_id": 789
	}`
	var client Client
	err := json.Unmarshal([]byte(jsonString), &client)
	assert.NilError(t, err)

	assert.Equal(t, client.ID, 1234)
	assert.Equal(t, client.Name, "ACME Inc.")
	assert.Assert(t, client.ArchivedAt.IsZero())
	assert.Assert(t, !client.IsArchived())
	assert.Equal(t, client.OwnerID, 567)
	assert.Equal(t, client.OfficeID, 789)
}

func TestClientArchived(t *testing.T) {
	jsonString := `
	{
		"id": 1234,
		"name": "ACME Inc.",
		"archived_at": "2018-06-07T07:27:54.563Z",
		"owner_id": 567,
		"office_id": 789
	}`
	var client Client
	err := json.Unmarshal([]byte(jsonString), &client)
	assert.NilError(t, err)

	archivedAt := time.Date(2018, time.Month(6), 7, 7, 27, 54, 563*1000000, time.UTC)
	assert.Equal(t, client.ArchivedAt, archivedAt)
	assert.Assert(t, client.IsArchived())
}
