package api

import (
	"errors"
	"fmt"
	"github.com/markosamuli/glassfactory/models"
	"net/http"
)

func NewClientsService(s *Service) *ClientsService {
	rs := &ClientsService{s: s}
	rs.clients = make(map[int]*models.Client)
	return rs
}

type ClientsService struct {
	s *Service
	clients map[int]*models.Client
}

func (r *ClientsService) Get(clientID int) (*models.Client, error) {
	client, ok := r.clients[clientID]
	if ok {
		return client, nil
	}
	res, err := r.Details(clientID).Do()
	if err != nil {
		return nil, err
	}
	r.clients[clientID] = res.Client
	return res.Client, nil
}

func (r *ClientsService) Details(clientID int) *ClientDetailsCall {
	c := &ClientDetailsCall{s: r.s}
	c.clientID = clientID
	return c
}

type ClientDetailsCall struct {
	s *Service
	clientID int
}

type ClientDetailsResponse struct {
	Client *models.Client
}

func (c *ClientDetailsCall) doRequest() (*http.Response, error) {
	var urls string
	if c.clientID > 0 {
		urls = c.s.BasePath + fmt.Sprintf("clients/%d.json", c.clientID)
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

func (c *ClientDetailsCall) Do() (*ClientDetailsResponse, error) {
	res, err := c.doRequest()
	if err != nil {
		return nil, err
	}
	var target models.Client
	if err := DecodeResponse(&target, res); err != nil {
		return nil, err
	}
	ret := &ClientDetailsResponse{}
	ret.Client = &target
	return ret, nil
}