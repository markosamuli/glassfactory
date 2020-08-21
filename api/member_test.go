package api

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/markosamuli/glassfactory/model"
	"gopkg.in/h2non/gock.v1"
	"gotest.tools/assert"
)

func TestMemberService_Get(t *testing.T) {
	defer gock.Off()

	domain := "https://example.glassfactory.io"
	apiPath := "/api/public/v1/"
	endpoint := domain + apiPath
	userID := 1401

	gock.New(domain).
		Get(apiPath + fmt.Sprintf("members/%d.json", userID)).
		Reply(200).
		BodyString(`{
		  "id": 1401,
		  "name": "First Last",
		  "email": "first.last@example.com",
		  "joined_at": "2016-02-17",
		  "archived_at": null,
		  "freelancer": false,
		  "role_id": 1019,
		  "capacity": 8,
		  "archived": false,
		  "office_id": 152
		}`)

	s := &Service{}
	s.client = &http.Client{}
	s.BasePath = endpoint

	// Member without data
	invalid := &model.Member{
		ID: userID,
	}

	var member *model.Member
	var err error
	rs := NewMemberService(s)

	// Add member without data into the cache
	rs.members.Add(invalid)

	member, err = rs.Get(userID, WithCache(false))
	assert.NilError(t, err)

	assert.Equal(t, member.ID, userID)
	assert.Equal(t, member.Name, "First Last")
	assert.Equal(t, member.Email, "first.last@example.com")
	assert.Equal(t, member.JoinedAt.In(time.UTC), time.Date(2016, time.February, 17, 0, 0, 0, 0, time.UTC))
	assert.Assert(t, !member.ArchivedAt.IsValid())
	assert.Assert(t, !member.Freelancer)
	assert.Equal(t, member.RoleID, 1019)
	assert.Equal(t, member.Capacity, 8.0)
	assert.Equal(t, member.Archived, false)
	assert.Equal(t, member.OfficeID, 152)

	// Verify that we don't have pending mocks
	assert.Assert(t, gock.IsDone(), "all mocks should have been called")
}

func TestMemberService_Get_WithCache(t *testing.T) {
	defer gock.Off()

	domain := "https://example.glassfactory.io"
	apiPath := "/api/public/v1/"
	endpoint := domain + apiPath
	userID := 1401

	gock.New(domain).
		Get(apiPath + fmt.Sprintf("members/%d.json", userID)).
		Reply(200).
		BodyString(`{
		  "id": 1401
		}`)

	s := &Service{}
	s.client = &http.Client{}
	s.BasePath = endpoint

	// Member without data
	expected := &model.Member{
		ID:    userID,
		Name:  "Cached Member",
		Email: "cached@email.com",
	}

	var member *model.Member
	var err error
	rs := NewMemberService(s)

	// Add member without data into the cache
	rs.members.Add(expected)

	// Get project with caching enabled
	member, err = rs.Get(userID, WithCache(true))
	assert.NilError(t, err)

	assert.Equal(t, member.ID, userID)
	assert.Equal(t, member, expected)

	assert.Equal(t, len(gock.Pending()), 1, "mock request shouldn't be called")
}

func TestMemberService_All(t *testing.T) {
	defer gock.Off()

	domain := "https://example.glassfactory.io"
	apiPath := "/api/public/v1/"
	endpoint := domain + apiPath

	gock.New(domain).
		Get(apiPath + "members.json").
		Reply(200).
		BodyString(`[
			{
			  "id": 123,
			  "name": "First",
			  "email": "first@example.com",
			  "joined_at": "2016-02-17",
			  "archived_at": null,
			  "freelancer": false,
			  "role_id": 1019,
			  "capacity": 8,
			  "archived": false,
			  "office_id": 152
			},
			{
			  "id": 456,
			  "name": "Second",
			  "email": "second@example.com",
			  "joined_at": "2016-02-17",
			  "archived_at": null,
			  "freelancer": false,
			  "role_id": 1019,
			  "capacity": 8,
			  "archived": false,
			  "office_id": 152
			},
			{
			  "id": 789,
			  "name": "Third",
			  "email": "third@example.com",
			  "joined_at": "2016-02-17",
			  "archived_at": "2019-02-17",
			  "freelancer": false,
			  "role_id": 1019,
			  "capacity": 8,
			  "archived": true,
			  "office_id": 152
			}
		]`)

	s := &Service{}
	s.client = &http.Client{}
	s.BasePath = endpoint

	var members []*model.Member
	var err error
	rs := NewMemberService(s)

	members, err = rs.All()
	assert.NilError(t, err)
	assert.Equal(t, len(members), 3)

	assert.Equal(t, members[0].ID, 123)
	assert.Equal(t, members[0].Name, "First")

	assert.Equal(t, members[1].ID, 456)
	assert.Equal(t, members[1].Name, "Second")

	assert.Equal(t, members[2].ID, 789)
	assert.Equal(t, members[2].Name, "Third")

	// Verify that we don't have pending mocks
	assert.Assert(t, gock.IsDone(), "all mocks should have been called")
}

func TestMemberService_List_WithStatusActive(t *testing.T) {
	defer gock.Off()

	domain := "https://example.glassfactory.io"
	apiPath := "/api/public/v1/"
	endpoint := domain + apiPath

	gock.New(domain).
		Get(apiPath + "members/active.json").
		Reply(200).
		BodyString(`[
			{
			  "id": 123,
			  "name": "First",
			  "email": "first@example.com",
			  "joined_at": "2016-02-17",
			  "archived_at": null,
			  "freelancer": false,
			  "role_id": 1019,
			  "capacity": 8,
			  "archived": false,
			  "office_id": 152
			},
			{
			  "id": 456,
			  "name": "Second",
			  "email": "second@example.com",
			  "joined_at": "2016-02-17",
			  "archived_at": null,
			  "freelancer": false,
			  "role_id": 1019,
			  "capacity": 8,
			  "archived": false,
			  "office_id": 152
			}
		]`)

	s := &Service{}
	s.client = &http.Client{}
	s.BasePath = endpoint

	var res *MemberListResponse
	var members []*model.Member
	var err error
	rs := NewMemberService(s)

	res, err = rs.List(WithStatus(memberStatusActive)).Do()
	assert.NilError(t, err)

	members = res.Members
	assert.Equal(t, len(members), 2)

	assert.Equal(t, members[0].ID, 123)
	assert.Equal(t, members[0].Name, "First")

	assert.Equal(t, members[1].ID, 456)
	assert.Equal(t, members[1].Name, "Second")

	// Verify that we don't have pending mocks
	assert.Assert(t, gock.IsDone(), "all mocks should have been called")
}

func TestMemberService_List_WithStatusArchived(t *testing.T) {
	defer gock.Off()

	domain := "https://example.glassfactory.io"
	apiPath := "/api/public/v1/"
	endpoint := domain + apiPath

	gock.New(domain).
		Get(apiPath + "members/archived.json").
		Reply(200).
		BodyString(`[
			{
			  "id": 789,
			  "name": "Third",
			  "email": "third@example.com",
			  "joined_at": "2016-02-17",
			  "archived_at": "2019-02-17",
			  "freelancer": false,
			  "role_id": 1019,
			  "capacity": 8,
			  "archived": true,
			  "office_id": 152
			}
		]`)

	s := &Service{}
	s.client = &http.Client{}
	s.BasePath = endpoint

	var res *MemberListResponse
	var members []*model.Member
	var err error
	rs := NewMemberService(s)

	res, err = rs.List(WithStatus(memberStatusArchived)).Do()
	assert.NilError(t, err)

	members = res.Members
	assert.Equal(t, len(members), 1)

	assert.Equal(t, members[0].ID, 789)
	assert.Equal(t, members[0].Name, "Third")
	assert.Equal(t, members[0].ArchivedAt.In(time.UTC), time.Date(2019, time.February, 17, 0, 0, 0, 0, time.UTC))
	assert.Assert(t, members[0].Archived)

	// Verify that we don't have pending mocks
	assert.Assert(t, gock.IsDone(), "all mocks should have been called")
}

func TestMemberService_List_WithTerm(t *testing.T) {
	defer gock.Off()

	domain := "https://example.glassfactory.io"
	apiPath := "/api/public/v1/"
	endpoint := domain + apiPath

	gock.New(domain).
		Get(apiPath+"members.json").
		MatchParam("term", "First").
		Reply(200).
		BodyString(`[
			{
			  "id": 123,
			  "name": "First",
			  "email": "first@example.com",
			  "joined_at": "2016-02-17",
			  "archived_at": null,
			  "freelancer": false,
			  "role_id": 1019,
			  "capacity": 8,
			  "archived": false,
			  "office_id": 152
			}
		]`)

	s := &Service{}
	s.client = &http.Client{}
	s.BasePath = endpoint

	var res *MemberListResponse
	var members []*model.Member
	var err error
	rs := NewMemberService(s)

	res, err = rs.List(WithTerm("First")).Do()
	assert.NilError(t, err)

	members = res.Members
	assert.Equal(t, len(members), 1)

	assert.Equal(t, members[0].ID, 123)
	assert.Equal(t, members[0].Name, "First")

	// Verify that we don't have pending mocks
	assert.Assert(t, gock.IsDone(), "all mocks should have been called")
}

func TestMemberService_Active(t *testing.T) {
	defer gock.Off()

	domain := "https://example.glassfactory.io"
	apiPath := "/api/public/v1/"
	endpoint := domain + apiPath

	gock.New(domain).
		Get(apiPath + "members/active.json").
		Reply(200).
		BodyString(`[
			{
			  "id": 123,
			  "name": "First",
			  "email": "first@example.com",
			  "joined_at": "2016-02-17",
			  "archived_at": null,
			  "freelancer": false,
			  "role_id": 1019,
			  "capacity": 8,
			  "archived": false,
			  "office_id": 152
			},
			{
			  "id": 456,
			  "name": "Second",
			  "email": "second@example.com",
			  "joined_at": "2016-02-17",
			  "archived_at": null,
			  "freelancer": false,
			  "role_id": 1019,
			  "capacity": 8,
			  "archived": false,
			  "office_id": 152
			}
		]`)

	s := &Service{}
	s.client = &http.Client{}
	s.BasePath = endpoint

	var members []*model.Member
	var err error
	rs := NewMemberService(s)

	members, err = rs.Active(WithCache(false))
	assert.NilError(t, err)
	assert.Equal(t, len(members), 2)

	assert.Equal(t, members[0].ID, 123)
	assert.Equal(t, members[0].Name, "First")

	assert.Equal(t, members[1].ID, 456)
	assert.Equal(t, members[1].Name, "Second")

	// Verify that we don't have pending mocks
	assert.Assert(t, gock.IsDone(), "all mocks should have been called")
}
