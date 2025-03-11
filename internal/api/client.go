package api

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

type Client struct {
	RestyClient *resty.Client
}

func NewClient(apiUrl, token string) *Client {
	client := resty.New().SetBaseURL(apiUrl).SetAuthToken(token)
	return &Client{RestyClient: client}
}

func (c *Client) Init(apiUrl string) {
	c.RestyClient.SetHostURL(apiUrl)
	c.RestyClient.SetAuthToken(viper.GetString("token"))
}

func (c *Client) GetRepositories() ([]Repository, error) {
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

func (c *Client) GetContexts() ([]Context, error) {
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
