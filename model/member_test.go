package model

import (
	"encoding/json"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestMember(t *testing.T) {
	jsonString := `{
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
	var member Member
	err := json.Unmarshal([]byte(jsonString), &member)
	assert.NilError(t, err)

	assert.Equal(t, member.ID, 1401)
	assert.Equal(t, member.Name, "First Last")
	assert.Equal(t, member.Email, "first.last@example.com")
	assert.Assert(t, member.JoinedAt.IsValid())
	assert.Assert(t, !member.ArchivedAt.IsValid())
	assert.Assert(t, !member.Freelancer)
	assert.Assert(t, !member.Archived)
	assert.Equal(t, member.Capacity, 8.0)
	assert.Equal(t, member.RoleID, 1019)
	assert.Equal(t, member.OfficeID, 152)
}

func TestMemberArchived(t *testing.T) {
	jsonString := `{
	  "id": 1401,
	  "user_id": 422,
	  "name": "First Last",
	  "email": "first.last@example.com",
	  "joined_at": "2016-02-17",
      "archived_at": "2018-06-07",
	  "freelancer": false,
	  "role_id": 1019,
	  "capacity": 8,
	  "archived": true,
	  "office_id": 152
	}`
	var member Member
	err := json.Unmarshal([]byte(jsonString), &member)
	assert.NilError(t, err)

	assert.Equal(t, member.ArchivedAt.In(time.UTC), time.Date(2018, time.Month(6), 7, 0, 0, 0, 0, time.UTC))
	assert.Assert(t, member.Archived)
}

func TestMemberFreelancer(t *testing.T) {
	jsonString := `{
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
	var member Member
	err := json.Unmarshal([]byte(jsonString), &member)
	assert.NilError(t, err)

	assert.Assert(t, member.Freelancer)
}

func TestMemberCapacity(t *testing.T) {
	jsonString := `{
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
	var member Member
	err := json.Unmarshal([]byte(jsonString), &member)
	assert.NilError(t, err)

	assert.Equal(t, member.Capacity, 7.5)
}
