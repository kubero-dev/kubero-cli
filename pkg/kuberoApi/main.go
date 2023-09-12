package kuberoApi

import (
	_ "embed"

	"github.com/go-resty/resty/v2"
)

type KuberoClient struct {
	baseURL     string
	bearerToken string
	client      *resty.Request
}

//go:embed VERSION
var version string

func (k *KuberoClient) Init(baseURL string, bearerToken string) *resty.Request {

	client := resty.New().SetBaseURL(baseURL).R().
		EnableTrace().
		SetAuthScheme("Bearer").
		SetAuthToken(bearerToken).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "kubero-cli/"+version)

	k.baseURL = baseURL
	k.bearerToken = bearerToken
	k.client = client

	return client
}

func (k *KuberoClient) DeployPipeline(pipeline PipelineCRD) (*resty.Response, error) {
	k.client.SetBody(pipeline.Spec)
	res, err := k.client.Post("/api/cli/pipelines/")

	return res, err
}

func (k *KuberoClient) UnDeployPipeline(pipelineName string) (*resty.Response, error) {
	res, err := k.client.Delete("/api/cli/pipelines/" + pipelineName)

	return res, err
}

func (k *KuberoClient) UnDeployApp(pipelineName string, stageName string, appName string) (*resty.Response, error) {
	res, err := k.client.Delete("/api/cli/pipelines/" + pipelineName + "/" + stageName + "/" + appName)

	return res, err
}
