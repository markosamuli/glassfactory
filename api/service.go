package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/markosamuli/glassfactory/model"
)

// NewService creates a new Service.
func NewService(ctx context.Context, settings *Settings) (*Service, error) {
	client, endpoint, err := NewClient(ctx, settings)
	if err != nil {
		return nil, err
	}
	s := &Service{client: client}
	s.Client = NewClientService(s)
	s.Member = NewMemberService(s)
	s.Project = NewProjectService(s)
	if endpoint != "" {
		s.BasePath = endpoint
	}
	s.settings = settings
	return s, nil
}

// Service is used for calling Glass Factory APIs
type Service struct {
	client        *http.Client
	settings      *Settings
	currentMember *model.Member
	BasePath      string // Base URL for the API

	Client  *ClientService
	Member  *MemberService
	Project *ProjectService
}

// GetCurrentMember returns a member matching the user email address in settings
func (s *Service) GetCurrentMember() (*model.Member, error) {
	// Get user email from settings
	email := s.settings.UserEmail
	if email == "" {
		return nil, errors.New("user email not defined")
	}
	// Get cached member if the email address matches
	if s.currentMember != nil && s.currentMember.Email == email {
		return s.currentMember, nil
	}
	// Get active members without caching them
	members, err := s.Member.Active(WithCache(false))
	if err != nil {
		return nil, err
	}
	// Find member with matching email address
	for _, member := range members {
		if member.Email == email {
			s.currentMember = member
			return member, nil
		}
	}
	return nil, fmt.Errorf("no users matching email %s found", email)
}

// DecodeResponse decodes the body of res into target. If there is no body,
// target is unchanged.
func DecodeResponse(target interface{}, res *http.Response) error {
	if res.StatusCode == http.StatusNoContent {
		return nil
	}
	return json.NewDecoder(res.Body).Decode(target)
}
