package auth

import (
	"fmt"
	"testing"

	"gotest.tools/assert"
)

func TestNewService(t *testing.T) {
	domain := "example"
	gfAuth := NewAuth()
	gfAuth.Account = domain
	gfAuth.Email = "test@example.com"
	gfAuth.APIKey = "abcdefg1234"
	api, err := gfAuth.NewService()
	assert.NilError(t, err)
	assert.Equal(t, api.BasePath, fmt.Sprintf("https://%s.glassfactory.io/api/public/v1/", domain))
}
