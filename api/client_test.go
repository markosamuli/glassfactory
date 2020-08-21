package api

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/markosamuli/glassfactory/model"
	"gopkg.in/h2non/gock.v1"
	"gotest.tools/assert"
)

func TestGetClient(t *testing.T) {
	defer gock.Off()

	domain := "https://example.glassfactory.io"
	apiPath := "/api/public/v1/"
	endpoint := domain + apiPath
	clientID := 1234

	gock.New(domain).
		Get(apiPath + fmt.Sprintf("clients/%d.json", clientID)).
		Reply(200).
		BodyString(`{
			"id": 1234,
			"name": "ACME Inc.",
			"archived_at": null,
			"owner_id": 567,
			"office_id": 789
		}`)

	s := &Service{}
	s.client = &http.Client{}
	s.BasePath = endpoint

	// Project without data
	invalid := &model.Client{
		ID: clientID,
	}

	var client *model.Client
	var err error
	rs := NewClientService(s)

	// Add client without data into the cache
	rs.clients.Add(invalid)

	client, err = rs.Get(clientID, WithCache(false))
	assert.NilError(t, err)

	assert.Equal(t, client.ID, clientID)
	assert.Equal(t, client.Name, "ACME Inc.")
	assert.Assert(t, client.ArchivedAt.IsZero())
	assert.Equal(t, client.OwnerID, 567)
	assert.Equal(t, client.OfficeID, 789)

	// Verify that we don't have pending mocks
	assert.Assert(t, gock.IsDone(), "all mocks should have been called")

}

func TestGetClientCached(t *testing.T) {
	defer gock.Off()

	domain := "https://example.glassfactory.io"
	apiPath := "/api/public/v1/"
	endpoint := domain + apiPath
	clientID := 1234

	gock.New(domain).
		Get(apiPath + fmt.Sprintf("clients/%d.json", clientID)).
		Reply(200).
		BodyString(`{
			"id": 1234
		}`)

	s := &Service{}
	s.client = &http.Client{}
	s.BasePath = endpoint

	// Project without data
	expected := &model.Client{
		ID:   clientID,
		Name: "Cached Limited",
	}

	var client *model.Client
	var err error
	rs := NewClientService(s)

	// Add client without data into the cache
	rs.clients.Add(expected)

	client, err = rs.Get(clientID, WithCache(true))
	assert.NilError(t, err)

	assert.Equal(t, client.ID, clientID)
	assert.Equal(t, client, expected)

	assert.Equal(t, len(gock.Pending()), 1, "mock request shouldn't be called")

}

func TestListClient(t *testing.T) {
	defer gock.Off()

	domain := "https://example.glassfactory.io"
	apiPath := "/api/public/v1/"
	endpoint := domain + apiPath

	gock.New(domain).
		Get(apiPath + "clients.json").
		Reply(200).
		BodyString(`[
			{
				"id": 1234,
				"name": "ACME Inc.",
				"archived_at": null,
				"owner_id": 567,
				"office_id": 789
			}
		]`)

	s := &Service{}
	s.client = &http.Client{}
	s.BasePath = endpoint

	var res *ClientListResponse
	var clients []*model.Client
	var err error
	rs := NewClientService(s)

	res, err = rs.List().Do()
	assert.NilError(t, err)

	clients = res.Clients
	assert.Equal(t, len(clients), 1)

	assert.Equal(t, clients[0].ID, 1234)
	assert.Equal(t, clients[0].Name, "ACME Inc.")

	// Verify that we don't have pending mocks
	assert.Assert(t, gock.IsDone(), "all mocks should have been called")
}
