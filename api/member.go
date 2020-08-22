package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/markosamuli/glassfactory/model"
)

const memberStatusAll = "all"
const memberStatusActive = "active"
const memberStatusArchived = "archived"

// NewMemberService initialises a new MemberService
func NewMemberService(s *Service) *MemberService {
	rs := &MemberService{s: s}
	rs.Reports = NewMemberReportsService(rs)
	rs.members = model.NewMemberCollection()
	return rs
}

// MemberService is used for calling the Glass Factory member APIs
type MemberService struct {
	s       *Service
	members *model.MemberCollection
	Reports *MemberReportsService
}

// GetActive returns all active members in the Glass Factory account
//
// Deprecated: Use MemberService.Active()
func (r *MemberService) GetActive(opts ...RequestOption) ([]*model.Member, error) {
	return r.Active(opts...)
}

// Active returns all active members in the Glass Factory account
func (r *MemberService) Active(opts ...RequestOption) ([]*model.Member, error) {
	opts = append(opts, WithStatus(memberStatusActive))
	res, err := r.List(opts...).Do()
	if err != nil {
		return nil, err
	}
	return res.Members, nil
}

// All returns all members in the Glass Factory account
func (r *MemberService) All(opts ...RequestOption) ([]*model.Member, error) {
	opts = append(opts, WithStatus(memberStatusAll))
	res, err := r.List(opts...).Do()
	if err != nil {
		return nil, err
	}
	return res.Members, nil
}

// Get returns a member from Glass Factory by their user ID
func (r *MemberService) Get(userID int, opts ...RequestOption) (*model.Member, error) {
	options := NewRequestOptions(opts)
	if options.cache {
		member, ok := r.members.Get(userID)
		if ok {
			return member, nil
		}
	}
	res, err := r.Details(userID).Do()
	if err != nil {
		return nil, err
	}
	if options.cache {
		r.members.Add(res.Member)
	}
	return res.Member, nil
}

// Details returns member details from Glass Factory
func (r *MemberService) Details(userID int) *MemberDetailsCall {
	c := &MemberDetailsCall{s: r.s}
	c.userID = userID
	return c
}

// FindByEmail lists all members in the Glass Factory account and looks for a member
// with matching email address.
//
// Deprecated: Filter members outside the service.
func (r *MemberService) FindByEmail(email string, opts ...RequestOption) (*model.Member, error) {
	res, err := r.List(opts...).Do()
	if err != nil {
		return nil, err
	}
	if len(res.Members) == 0 {
		return nil, errors.New("no users found")
	}
	for _, member := range res.Members {
		if member.Email == email {
			return member, nil
		}
	}
	return nil, fmt.Errorf("no users matching email %s found", email)
}

// SearchActive returns a list of active users matching the search term.
//
// Deprecated: You should call MemberService.List() directly with WithTerm(term)
// and WithStatus("active")) options instead.
func (r *MemberService) SearchActive(term string) *MemberListCall {
	return r.List(WithTerm(term), WithStatus(memberStatusActive))
}

// List returns a list of staff members in the Glass Factory account
func (r *MemberService) List(opts ...RequestOption) *MemberListCall {
	c := &MemberListCall{s: r.s}
	c.options = opts
	return c
}

// ListActive returns a list of active staff members in the Glass Factory account
//
// Deprecated: You should call MemberService.List(WithStatus("active")) instead.
func (r *MemberService) ListActive() *MemberListCall {
	return r.List(WithStatus(memberStatusActive))
}

// ListArchived returns a list of archived staff members in the Glass Factory account
//
// Deprecated: You should call MemberService.List(WithStatus("archived")) instead.
func (r *MemberService) ListArchived() *MemberListCall {
	return r.List(WithStatus(memberStatusArchived))
}

// MemberDetailsCall represents a request to Member Details API
type MemberDetailsCall struct {
	s      *Service
	userID int
}

// MemberDetailsResponse represents a response from Member Details API
type MemberDetailsResponse struct {
	Member *model.Member
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

// Do executes the request and parses results
func (c *MemberDetailsCall) Do() (*MemberDetailsResponse, error) {
	res, err := c.doRequest()
	if err != nil {
		return nil, err
	}
	var target model.Member
	if err := DecodeResponse(&target, res); err != nil {
		return nil, err
	}
	ret := &MemberDetailsResponse{}
	ret.Member = &target
	return ret, nil
}

// MemberListCall represents a request to List Staff Members API
type MemberListCall struct {
	s       *Service
	options []RequestOption
}

// Options returns request options with defaults
func (c *MemberListCall) Options() RequestOptions {
	options := RequestOptions{
		status: memberStatusAll,
	}
	options.apply(c.options)
	return options
}

// MemberListResponse represents a response from List Staff Members API
type MemberListResponse struct {
	Members []*model.Member
	Status  string
}

func (c *MemberListCall) doRequest() (*http.Response, error) {
	var urls string

	options := c.Options()
	switch options.status {
	case memberStatusActive:
		urls = c.s.BasePath + "members/active.json"
	case memberStatusArchived:
		urls = c.s.BasePath + "members/archived.json"
	case memberStatusAll:
		urls = c.s.BasePath + "members.json"
	}

	urlParams := url.Values{}
	if options.term != "" {
		urlParams.Add("term", options.term)
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

// Do executes the request and parses results
func (c *MemberListCall) Do() (*MemberListResponse, error) {
	res, err := c.doRequest()
	if err != nil {
		return nil, err
	}
	target := make([]*model.Member, 0)
	if err := DecodeResponse(&target, res); err != nil {
		return nil, err
	}
	ret := &MemberListResponse{}
	ret.Members = target
	//ret.Status = c.status
	return ret, nil
}
