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

func (k *KuberoClient) CreatePipeline(pipeline PipelineCRD) (*resty.Response, error) {
	k.client.SetBody(pipeline.Spec)
	pipelineResp, pipelineErr := k.client.Post("/api/cli/pipelines/")

	return pipelineResp, pipelineErr
}
