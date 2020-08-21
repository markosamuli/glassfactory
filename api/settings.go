package api

import (
	"errors"
)

// Settings for the HTTP client
type Settings struct {
	UserEmail        string
	UserToken        string
	AccountSubdomain string
}

func NewSettings(account string, email string, token string) (*Settings, error) {
	settings := &Settings{}
	settings.AccountSubdomain = account
	settings.UserEmail = email
	settings.UserToken = token
	err := settings.Validate()
	if err != nil {
		return nil, err
	}
	return settings, nil
}

// Validate returns an error if any settings are missing
func (s *Settings) Validate() error {
	if s.UserEmail == "" {
		return errors.New("user email missing")
	}
	if s.UserToken == "" {
		return errors.New("user API key missing")
	}
	if s.AccountSubdomain == "" {
		return errors.New("account subdomain missing")
	}
	return nil
}
