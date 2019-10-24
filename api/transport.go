package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const GlassFactoryPublicAPI = "https://%s.glassfactory.io/api/public/v1/"

// GlassFactoryTransport for the HTTP client provides the custom authentication headers
type GlassFactoryAuthTransport struct {
	UserEmail string
	UserToken string
	AccountSubdomain string
}

// RoundTrip authorizes the request with HTTP Basic Auth
func (trans GlassFactoryAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-User-Email", trans.UserEmail)
	req.Header.Set("X-User-Token", trans.UserToken)
	req.Header.Set("X-Account-Subdomain", trans.AccountSubdomain)
	return http.DefaultTransport.RoundTrip(req)
}

func NewClient(ctx context.Context, gfs *GlassFactorySettings) (*http.Client, string, error) {
	if gfs.UserEmail == "" || gfs.UserToken == "" {
		return nil, "", errors.New("user email or API key missing")
	}
	if gfs.AccountSubdomain == "" {
		return nil, "", errors.New("account subdomain missing")
	}
	trans := &GlassFactoryAuthTransport{
		UserEmail: gfs.UserEmail,
		UserToken: gfs.UserToken,
		AccountSubdomain: gfs.AccountSubdomain}
	endpoint := fmt.Sprintf(GlassFactoryPublicAPI, gfs.AccountSubdomain)
	return &http.Client{Transport: trans}, endpoint, nil
}