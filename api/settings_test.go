package api

import (
	"testing"

	"gotest.tools/assert"
)

func TestNewSettings(t *testing.T) {
	var tests = []struct {
		name string
		account   string
		email string
		token   string
		err string
	}{
		{
			name: "all missing",
			account: "",
			email: "",
			token: "",
			err: "user email missing",
		},
		{
			name: "user API key missing",
			account: "example",
			email: "test@example.com",
			token: "",
			err: "user API key missing",
		},
		{
			name: "account subdomain missing",
			account: "",
			email: "test@example.com",
			token: "abcdefg1234",
			err: "account subdomain missing",
		},
		{
			name: "valid",
			account: "example",
			email: "test@example.com",
			token: "abcdefg1234",
			err: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			settings, err := NewSettings(tt.account, tt.email, tt.token)
			if tt.err != "" {
				assert.Error(t, err, tt.err)
			} else {
				assert.NilError(t, err)
				assert.Equal(t, settings.AccountSubdomain, tt.account)
				assert.Equal(t, settings.UserEmail, tt.email)
				assert.Equal(t, settings.UserToken, tt.token)
			}
		})
	}
}