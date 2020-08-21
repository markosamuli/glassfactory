package api

import (
	"context"
	"fmt"
	"net/http"
)

const publicAPI = "https://%s.glassfactory.io/api/public/v1/"

// Transport for the HTTP client provides the custom authentication headers
type AuthTransport struct {
	UserEmail        string
	UserToken        string
	AccountSubdomain string
}

// RoundTrip authorizes the request with HTTP Basic Auth
func (trans AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-User-Email", trans.UserEmail)
	req.Header.Set("X-User-Token", trans.UserToken)
	req.Header.Set("X-Account-Subdomain", trans.AccountSubdomain)
	return http.DefaultTransport.RoundTrip(req)
}

// NewClient returns a new HTTP client with Glass Factory authentication
func NewClient(ctx context.Context, settings *Settings) (*http.Client, string, error) {
	err := settings.Validate()
	if err != nil {
		return nil, "", err
	}
	trans := &AuthTransport{
		UserEmail:        settings.UserEmail,
		UserToken:        settings.UserToken,
		AccountSubdomain: settings.AccountSubdomain}
	endpoint := fmt.Sprintf(publicAPI, settings.AccountSubdomain)
	return &http.Client{Transport: trans}, endpoint, nil
}
