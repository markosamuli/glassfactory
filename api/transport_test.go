package api

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"gopkg.in/h2non/gock.v1"
	"gotest.tools/assert"
)

func newTestSettings() *Settings {
	settings := &Settings{}
	settings.UserEmail = "test@example.com"
	settings.UserToken = "abcdefg1234"
	settings.AccountSubdomain = "example"
	return settings
}

func TestNewClient(t *testing.T) {
	ctx := context.Background()
	settings := newTestSettings()
	_, endpoint, err := NewClient(ctx, settings)
	assert.NilError(t, err)
	assert.Equal(t, endpoint, fmt.Sprintf(publicAPI, settings.AccountSubdomain))
}

func TestAuthTransport(t *testing.T) {
	defer gock.Off()

	settings := newTestSettings()

	gock.New("https://example.glassfactory.io").
		Get("/api/public/v1/dummy").
		MatchHeaders(map[string]string{
			"X-User-Email":        settings.UserEmail,
			"X-User-Token":        settings.UserToken,
			"X-Account-Subdomain": settings.AccountSubdomain,
		}).
		Reply(200).
		BodyString(`{
		  "ok": true
		}`)

	var client *http.Client
	var req *http.Request
	var res *http.Response
	var endpoint string
	var err error

	ctx := context.Background()
	client, endpoint, err = NewClient(ctx, settings)
	assert.NilError(t, err)
	assert.Equal(t, endpoint, "https://example.glassfactory.io/api/public/v1/")

	urls := "https://example.glassfactory.io/api/public/v1/dummy"
	req, err = http.NewRequest(http.MethodGet, urls, nil)
	assert.NilError(t, err)

	res, err = client.Do(req)
	assert.NilError(t, err)
	assert.Equal(t, res.Request.Header.Get("X-User-Email"), settings.UserEmail)
	assert.Equal(t, res.Request.Header.Get("X-User-Token"), settings.UserToken)
	assert.Equal(t, res.Request.Header.Get("X-Account-Subdomain"), settings.AccountSubdomain)

	// Verify that we don't have pending mocks
	assert.Assert(t, gock.IsDone(), "all mocks should have been called")
}
