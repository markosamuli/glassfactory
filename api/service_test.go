package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/markosamuli/glassfactory/model"
	"gopkg.in/h2non/gock.v1"
	"gotest.tools/assert"
)

func TestNewService(t *testing.T) {
	ctx := context.Background()

	settings := &Settings{}
	settings.UserEmail = "test@example.com"
	settings.UserToken = "abcdefg1234"
	settings.AccountSubdomain = "example"

	t.Run("valid settings", func(t *testing.T) {
		service, err := NewService(ctx, settings)
		assert.NilError(t, err)
		assert.Equal(t, service.BasePath, fmt.Sprintf(publicAPI, "example"))
	})

	t.Run("missing settings", func(t *testing.T) {
		settings.AccountSubdomain = ""

		_, err := NewService(ctx, settings)
		assert.Error(t, err, "account subdomain missing")
	})
}

func TestGetCurrentMember(t *testing.T) {
	defer gock.Off()

	ctx := context.Background()

	domain := "https://example.glassfactory.io"
	apiPath := "/api/public/v1/"
	email := "test@example.com"

	models := []*model.Member{
		{
			ID:    1401,
			Email: "test@example.com",
		},
	}
	data, err := json.Marshal(&models)
	assert.NilError(t, err)

	t.Run("user found", func(t *testing.T) {
		gock.New(domain).
			Get(apiPath + "members/active.json").
			Reply(200).
			BodyString(string(data))

		settings := &Settings{
			UserEmail:        email,
			UserToken:        "abcdefg1234",
			AccountSubdomain: "example",
		}
		service, err := NewService(ctx, settings)
		assert.NilError(t, err)

		member, err := service.GetCurrentMember()
		assert.NilError(t, err)
		assert.Equal(t, member.ID, 1401)
		assert.Equal(t, member.Email, email)

		member, err = service.GetCurrentMember()
		assert.NilError(t, err)
		assert.Equal(t, member.ID, 1401)
		assert.Equal(t, member.Email, email)
	})

	t.Run("user not found", func(t *testing.T) {
		gock.New(domain).
			Get(apiPath + "members/active.json").
			Reply(200).
			BodyString(string(data))

		settings := &Settings{
			UserEmail:        "invalid@email.com",
			UserToken:        "abcdefg1234",
			AccountSubdomain: "example",
		}
		service, err := NewService(ctx, settings)
		assert.NilError(t, err)

		_, err = service.GetCurrentMember()
		assert.Error(t, err, "no users matching email invalid@email.com found")
	})

	// Verify that we don't have pending mocks
	assert.Assert(t, gock.IsDone(), "all mocks should have been called")
}

func newHttpResponseWithJsonBody(body string) *http.Response {
	header := http.Header{}
	header.Add("Content-Type", "application/json, charset=UTF-8")
	return &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.0",
		ProtoMajor:    1,
		ProtoMinor:    0,
		Header:        header,
		Uncompressed:  true,
		ContentLength: -1,
		Body:          ioutil.NopCloser(strings.NewReader(body)),
	}
}

func TestDecodeResponse(t *testing.T) {
	var tests = []struct {
		name   string
		target interface{}
		body   string
	}{
		{
			name:   "Client",
			target: model.Client{},
			body: `{
				"id": 1234,
				"name": "Google",
				"archived_at": "2018-06-07T07:27:54.563Z",
				"owner_id": 567,
				"office_id": 789
			}`,
		},
		{
			name:   "Member",
			target: model.Member{},
			body: `{
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
			}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := newHttpResponseWithJsonBody(tt.body)
			err := DecodeResponse(&tt.target, res)
			assert.NilError(t, err)
		})
	}
}
