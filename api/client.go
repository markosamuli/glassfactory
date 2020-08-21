package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/markosamuli/glassfactory/model"
)

// NewClientService initialises a new ClientService
func NewClientService(s *Service) *ClientService {
	rs := &ClientService{s: s}
	rs.clients = model.NewClientCollection()
	return rs
}

// ClientService is used for calling the Glass Factory client APIs
type ClientService struct {
	s       *Service
	clients *model.ClientCollection
}

// All returns all clients in the Glass Factory account
func (r *ClientService) All(opts ...Option) ([]*model.Client, error) {
	res, err := r.List(opts...).Do()
	if err != nil {
		return nil, err
	}
	return res.Clients, nil
}

// Get returns a client from Glass Factory
func (r *ClientService) Get(clientID int, opts ...Option) (*model.Client, error) {
	options := Options(opts)
	if options.cache {
		client, ok := r.clients.Get(clientID)
		if ok {
			return client, nil
		}
	}
	res, err := r.Details(clientID).Do()
	if err != nil {
		return nil, err
	}
	if options.cache {
		r.clients.Add(res.Client)
	}
	return res.Client, nil
}

// Details returns client details from Glass Factory
func (r *ClientService) Details(clientID int) *ClientDetailsCall {
	c := &ClientDetailsCall{s: r.s}
	c.clientID = clientID
	return c
}

// List returns a list of clients in the Glass Factory account
func (r *ClientService) List(opts ...Option) *ClientListCall {
	c := &ClientListCall{s: r.s}
	c.options = opts
	return c
}

// ClientDetailsCall represents a request to Client Details API
type ClientDetailsCall struct {
	s        *Service
	clientID int
}

// ClientDetailsResponse represents a response from Client Details API
type ClientDetailsResponse struct {
	Client *model.Client
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

// Do executes the request and parses results
func (c *ClientDetailsCall) Do() (*ClientDetailsResponse, error) {
	res, err := c.doRequest()
	if err != nil {
		return nil, err
	}
	var target model.Client
	if err := DecodeResponse(&target, res); err != nil {
		return nil, err
	}
	ret := &ClientDetailsResponse{}
	ret.Client = &target
	return ret, nil
}

// ClientListCall represents a request to List Account's Clients API
type ClientListCall struct {
	s       *Service
	options []Option
}

// Options returns request options with defaults
func (c *ClientListCall) Options() options {
	options := options{}
	options.apply(c.options)
	return options
}

// ClientListResponse represents a response from List Account's Clients API
type ClientListResponse struct {
	Clients []*model.Client
}

func (c *ClientListCall) doRequest() (*http.Response, error) {
	var urls string
	urls = c.s.BasePath + "clients.json"

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
func (c *ClientListCall) Do() (*ClientListResponse, error) {
	res, err := c.doRequest()
	if err != nil {
		return nil, err
	}
	target := make([]*model.Client, 0)
	if err := DecodeResponse(&target, res); err != nil {
		return nil, err
	}
	ret := &ClientListResponse{}
	ret.Clients = target
	return ret, nil
}
