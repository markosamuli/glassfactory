package model

import (
	"encoding/json"
	"testing"
	"time"

	"gotest.tools/assert"
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
	assert.NilError(t, err)

	assert.Equal(t, target.ID, 1401)
	assert.Equal(t, target.Name, "First Last")
	assert.Equal(t, target.Email, "first.last@example.com")
	assert.Assert(t, target.JoinedAt.IsValid())
	assert.Assert(t, target.ArchivedAt.IsZero())
	assert.Assert(t, !target.Freelancer)
	assert.Assert(t, !target.Archived)
	assert.Equal(t, target.Capacity, float32(8.0))
	assert.Equal(t, target.RoleID, 1019)
	assert.Equal(t, target.OfficeID, 152)
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
	assert.NilError(t, err)

	archivedAt := time.Date(2018, time.Month(6), 7, 7, 27, 54, 563 * 1000000, time.UTC)
	assert.Equal(t, target.ArchivedAt, archivedAt)
	assert.Assert(t, target.Archived)
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
	assert.NilError(t, err)

	assert.Assert(t, target.Freelancer)
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
	assert.NilError(t, err)

	assert.Equal(t, target.Capacity, float32(7.5))
}