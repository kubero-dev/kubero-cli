package api

import (
	"fmt"
	"github.com/faelmori/kubero-cli/types"
	"github.com/faelmori/kubero-cli/version"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	baseURL     string
	bearerToken string
	client      *resty.Request
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Init(baseURL string, bearerToken string) *resty.Request {
	if baseURL == "" || bearerToken == "" {
		panic("baseURL and bearerToken are required to initialize the API client")
	}

	client := resty.New().SetBaseURL(baseURL).R().
		EnableTrace().
		SetAuthScheme("Bearer").
		SetAuthToken(bearerToken).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "kubero-cli/"+version.Version())

	c.baseURL = baseURL
	c.bearerToken = bearerToken
	c.client = client

	return client
}
func (c *Client) DeployPipeline(pipeline types.PipelineCRD) (*resty.Response, error) {
	c.client.SetBody(pipeline.Spec)
	res, err := c.client.Post("/api/cli/pipelines/")

	return res, err
}
func (c *Client) UnDeployPipeline(pipelineName string) (*resty.Response, error) {
	res, err := c.client.Delete("/api/cli/pipelines/" + pipelineName)

	return res, err
}
func (c *Client) GetPipeline(pipelineName string) (*resty.Response, error) {
	res, err := c.client.Get("/api/cli/pipelines/" + pipelineName)

	return res, err
}
func (c *Client) UnDeployApp(pipelineName string, stageName string, appName string) (*resty.Response, error) {
	res, err := c.client.Delete("/api/cli/pipelines/" + pipelineName + "/" + stageName + "/" + appName)

	return res, err
}
func (c *Client) GetApp(pipelineName string, stageName string, appName string) (*resty.Response, error) {
	res, err := c.client.Get("/api/cli/pipelines/" + pipelineName + "/" + stageName + "/" + appName)

	return res, err
}
func (c *Client) GetApps() (*resty.Response, error) {
	res, err := c.client.Get("/api/cli/apps")

	return res, err
}
func (c *Client) GetPipelines() (*resty.Response, error) {
	res, err := c.client.Get("/api/cli/pipelines")
	return res, handleError(res, err)
}
func (c *Client) DeployApp(app types.AppCRD) (*resty.Response, error) {
	c.client.SetBody(app.Spec)
	res, err := c.client.Post("/api/cli/apps")

	return res, err
}
func (c *Client) GetPipelineApps(pipelineName string) (*resty.Response, error) {
	res, err := c.client.Get("/api/cli/pipelines/" + pipelineName + "/apps")

	return res, err
}
func (c *Client) GetAddons() (*resty.Response, error) {
	res, err := c.client.Get("/api/cli/addons")

	return res, err
}
func (c *Client) GetBuildpacks() (*resty.Response, error) {
	res, err := c.client.Get("/api/cli/config/buildpacks")

	return res, err
}
func (c *Client) GetPodsize() (*resty.Response, error) {
	res, err := c.client.Get("/api/cli/config/podsize")

	return res, err
}
func (c *Client) GetRepositories() (*resty.Response, error) {
	res, err := c.client.Get("/api/cli/config/repositories")

	return res, err
}
func (c *Client) GetContexts() (*resty.Response, error) {
	res, err := c.client.Get("/api/cli/config/k8s/context")

	return res, err
}
func (c *Client) WithBody(body interface{}) *Client {
	c.client.SetBody(body)
	return c
}

func handleError(response *resty.Response, err error) error {
	if err != nil {
		return err
	}

	if response.IsError() {
		return fmt.Errorf("API error: %s", response.String())
	}

	return nil
}
