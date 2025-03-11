package api

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
)

type Repository interface {
	GetRepositories() ([]Repository, error)
	GetContexts() ([]Context, error)
}

type Context interface {
	GetName() string
	GetURL() string
}

type repository struct {
	Name string
	URL  string
}

func (r *repository) GetName() string {
	return r.Name
}

func (r *repository) GetURL() string {
	return r.URL
}

type ClientAPI interface {
	Init(apiUrl string)
	GetRepositories() ([]Repository, error)
	GetContexts() ([]Context, error)
	// Add more methods as necessary
}

type NewClientAPI struct {
	RestyClient *resty.Client
}

func (c *NewClientAPI) Init(apiUrl string) {
	c.RestyClient.SetHostURL(apiUrl)
}

func (c *NewClientAPI) GetRepositories() ([]Repository, error) {
	resp, err := c.RestyClient.R().Get("/api/repositories")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, errors.New("failed to fetch repositories")
	}
	var repos []Repository
	err = json.Unmarshal(resp.Body(), &repos)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func (c *NewClientAPI) GetContexts() ([]Context, error) {
	resp, err := c.RestyClient.R().Get("/api/contexts")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, errors.New("failed to fetch contexts")
	}
	var contexts []Context
	err = json.Unmarshal(resp.Body(), &contexts)
	if err != nil {
		return nil, err
	}
	return contexts, nil
}
