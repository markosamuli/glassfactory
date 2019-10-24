package api

import (
	"errors"
	"fmt"
	"github.com/markosamuli/glassfactory/models"
	"net/http"
	"net/url"
)

func NewMembersService(s *Service) *MembersService {
	rs := &MembersService{s: s}
	rs.members = make(map[int]*models.Member)
	return rs
}

type MembersService struct {
	s *Service
	members map[int]*models.Member
}

func (r *MembersService) GetCurrentMember() (*models.Member, error) {
	return r.FindByEmail(r.s.settings.UserEmail)
}

func (r *MembersService) GetActive() ([]*models.Member, error) {
	res, err := r.ListActive().Do()
	if err != nil {
		return nil, err
	}
	for _, m := range res.Members {
		r.members[m.ID] = m
	}
	return res.Members, nil
}

func (r *MembersService) Get(userID int) (*models.Member, error) {
	client, ok := r.members[userID]
	if ok {
		return client, nil
	}
	res, err := r.Details(userID).Do()
	if err != nil {
		return nil, err
	}
	r.members[userID] = res.Member
	return res.Member, nil
}

func (r *MembersService) Details(userID int) *MemberDetailsCall {
	c := &MemberDetailsCall{s: r.s}
	c.userID = userID
	return c
}

func (r *MembersService) FindByEmail(email string) (*models.Member, error) {
	for _, m := range r.members {
		if m.Email == email {
			return m, nil
		}
	}
	res, err := r.ListActive().Do()
	if err != nil {
		return nil, err
	}
	if len(res.Members) == 0 {
		return nil, errors.New("no active users found")
	}
	for _, m := range res.Members {
		if m.Email == email {
			r.members[m.ID] = m
			return m, nil
		}
	}
	return nil, fmt.Errorf("no users matching email %s found", email)
}

func (r *MembersService) SearchActive(term string) *MembersListCall {
	c := &MembersListCall{s: r.s}
	c.term = term
	c.status = models.MemberStatusActive
	return c
}

func (r *MembersService) List() *MembersListCall {
	c := &MembersListCall{s: r.s}
	return c
}

func (r *MembersService) ListActive() *MembersListCall {
	c := &MembersListCall{s: r.s}
	c.status = models.MemberStatusActive
	return c
}

func (r *MembersService) ListArchived() *MembersListCall {
	c := &MembersListCall{s: r.s}
	c.status = models.MemberStatusArchived
	return c
}

type MemberDetailsCall struct {
	s *Service
	userID int
}

type MemberDetailsResponse struct {
	Member *models.Member
}

func (c *MemberDetailsCall) doRequest() (*http.Response, error) {
	var urls string
	if c.userID > 0 {
		urls = c.s.BasePath + fmt.Sprintf("members/%d.json", c.userID)
	} else {
		return nil, errors.New("client ID is required")
	}

	req, err := http.NewRequest(http.MethodGet, urls, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *MemberDetailsCall) Do() (*MemberDetailsResponse, error) {
	res, err := c.doRequest()
	if err != nil {
		return nil, err
	}
	var target models.Member
	if err := DecodeResponse(&target, res); err != nil {
		return nil, err
	}
	ret := &MemberDetailsResponse{}
	ret.Member = &target
	return ret, nil
}

type MembersListCall struct {
	s *Service
	status string
	term string
}

type MembersListResponse struct {
	Members []*models.Member
	Status string
}

func (c *MembersListCall) doRequest() (*http.Response, error) {
	var urls string
	switch c.status {
	case models.MemberStatusActive:
		urls = c.s.BasePath + "members/active.json"
	case models.MemberStatusArchived:
		urls = c.s.BasePath + "members/archived.json"
	default:
		urls = c.s.BasePath + "members.json"
	}

	urlParams := url.Values{}
	if c.term != "" {
		urlParams.Add("term", c.term)
	}
	urls += "?" + urlParams.Encode()

	req, err := http.NewRequest(http.MethodGet, urls, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := c.s.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *MembersListCall) Do() (*MembersListResponse, error) {
	res, err := c.doRequest()
	if err != nil {
		return nil, err
	}
	target := make([]*models.Member, 0)
	if err := DecodeResponse(&target, res); err != nil {
		return nil, err
	}
	ret := &MembersListResponse{}
	ret.Members = target
	ret.Status = c.status
	return ret, nil
}