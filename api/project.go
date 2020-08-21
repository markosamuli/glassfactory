package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/markosamuli/glassfactory/model"
)

// NewProjectService initialises a new ProjectService
func NewProjectService(s *Service) *ProjectService {
	rs := &ProjectService{s: s}
	rs.projects = model.NewProjectCollection()
	return rs
}

// ProjectService is used for calling the Glass Factory project APIs
type ProjectService struct {
	s        *Service
	projects *model.ProjectCollection
}

// All returns all projects in the Glass Factory account
func (r *ProjectService) All(opts ...Option) ([]*model.Project, error) {
	res, err := r.List(opts...).Do()
	if err != nil {
		return nil, err
	}
	return res.Projects, nil
}

// Get returns a project from Glass Factory
func (r *ProjectService) Get(projectID int, opts ...Option) (*model.Project, error) {
	options := Options(opts)
	if options.cache {
		project, ok := r.projects.Get(projectID)
		if ok {
			return project, nil
		}
	}
	res, err := r.Details(projectID).Do()
	if err != nil {
		return nil, err
	}
	if options.cache {
		r.projects.Add(res.Project)
	}
	return res.Project, nil
}

// Details returns project details from Glass Factory
func (r *ProjectService) Details(projectID int) *ProjectDetailsCall {
	c := &ProjectDetailsCall{s: r.s}
	c.projectID = projectID
	return c
}

// List returns a list of projects in the Glass Factory account
func (r *ProjectService) List(opts ...Option) *ProjectListCall {
	c := &ProjectListCall{s: r.s}
	c.options = opts
	return c
}

// ProjectDetailsCall represents a request to Project Details API
type ProjectDetailsCall struct {
	s         *Service
	projectID int
}

// ProjectDetailsResponse represents a response from Project Details API
type ProjectDetailsResponse struct {
	Project *model.Project
}

func (c *ProjectDetailsCall) doRequest() (*http.Response, error) {
	var urls string
	if c.projectID > 0 {
		urls = c.s.BasePath + fmt.Sprintf("projects/%d.json", c.projectID)
	} else {
		return nil, errors.New("project ID is required")
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
func (c *ProjectDetailsCall) Do() (*ProjectDetailsResponse, error) {
	res, err := c.doRequest()
	if err != nil {
		return nil, err
	}
	var target model.Project
	if err := DecodeResponse(&target, res); err != nil {
		return nil, err
	}
	ret := &ProjectDetailsResponse{}
	ret.Project = &target
	return ret, nil
}

// ProjectListCall represents a request to List Account's PRojects API
type ProjectListCall struct {
	s       *Service
	options []Option
}

// Options returns request options with defaults
func (c *ProjectListCall) Options() options {
	options := options{}
	options.apply(c.options)
	return options
}

// ProjectListResponse represents a response from List Account's PRojects API
type ProjectListResponse struct {
	Projects []*model.Project
}

func (c *ProjectListCall) doRequest() (*http.Response, error) {
	var urls string
	urls = c.s.BasePath + "projects.json"

	options := c.Options()

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
func (c *ProjectListCall) Do() (*ProjectListResponse, error) {
	res, err := c.doRequest()
	if err != nil {
		return nil, err
	}
	target := make([]*model.Project, 0)
	if err := DecodeResponse(&target, res); err != nil {
		return nil, err
	}
	ret := &ProjectListResponse{}
	ret.Projects = target
	return ret, nil
}
