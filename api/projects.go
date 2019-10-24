package api

import (
	"errors"
	"fmt"
	"github.com/markosamuli/glassfactory/models"
	"net/http"
)

func NewProjectsService(s *Service) *ProjectsService {
	rs := &ProjectsService{s: s}
	rs.projects = make(map[int]*models.Project)
	return rs
}

type ProjectsService struct {
	s *Service
	projects map[int]*models.Project
}

func (r *ProjectsService) Get(projectID int) (*models.Project, error) {
	project, ok := r.projects[projectID]
	if ok {
		return project, nil
	}
	res, err := r.Details(projectID).Do()
	if err != nil {
		return nil, err
	}
	r.projects[projectID] = res.Project
	return res.Project, nil
}

func (r *ProjectsService) Details(projectID int) *ProjectDetailsCall {
	c := &ProjectDetailsCall{s: r.s}
	c.projectID = projectID
	return c
}

type ProjectDetailsCall struct {
	s *Service
	projectID int
}

type ProjectDetailsResponse struct {
	Project *models.Project
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

func (c *ProjectDetailsCall) Do() (*ProjectDetailsResponse, error) {
	res, err := c.doRequest()
	if err != nil {
		return nil, err
	}
	var target models.Project
	if err := DecodeResponse(&target, res); err != nil {
		return nil, err
	}
	ret := &ProjectDetailsResponse{}
	ret.Project = &target
	return ret, nil
}