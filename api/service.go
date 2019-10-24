package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/markosamuli/glassfactory/settings"
	"github.com/markosamuli/glassfactory/transport"
	"net/http"
)

// NewService creates a new Service.
func NewService(ctx context.Context, gfs *settings.GlassFactorySettings) (*Service, error) {
	client, endpoint, err := transport.NewClient(ctx, gfs)
	if err != nil {
		return nil, err
	}
	s, err := New(client)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		s.BasePath = endpoint
	}
	s.settings = gfs
	return s, nil
}

func New(client *http.Client) (*Service, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	s := &Service{client: client}
	s.Clients = NewClientsService(s)
	s.Members = NewMembersService(s)
	s.Projects = NewProjectsService(s)
	s.Reports = NewReportsService(s)
	return s, nil
}

type Service struct {
	client *http.Client
	settings *settings.GlassFactorySettings
	BasePath string // Base URL for the API

	Clients *ClientsService
	Members *MembersService
	Projects *ProjectsService
	Reports *ReportsService
}

// DecodeResponse decodes the body of res into target. If there is no body,
// target is unchanged.
func DecodeResponse(target interface{}, res *http.Response) error {
	if res.StatusCode == http.StatusNoContent {
		return nil
	}
	return json.NewDecoder(res.Body).Decode(target)
}