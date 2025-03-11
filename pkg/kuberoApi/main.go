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
var Version string

func NewKuberoClient() *KuberoClient {
	return &KuberoClient{}
}

func (k *KuberoClient) Init(baseURL string, bearerToken string) *resty.Request {

	client := resty.New().SetBaseURL(baseURL).R().
		EnableTrace().
		SetAuthScheme("Bearer").
		SetAuthToken(bearerToken).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "kubero-cli/"+Version)

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

func (k *KuberoClient) GetPipeline(pipelineName string) (*resty.Response, error) {
	res, err := k.client.Get("/api/cli/pipelines/" + pipelineName)

	return res, err
}

func (k *KuberoClient) UnDeployApp(pipelineName string, stageName string, appName string) (*resty.Response, error) {
	res, err := k.client.Delete("/api/cli/pipelines/" + pipelineName + "/" + stageName + "/" + appName)

	return res, err
}

func (k *KuberoClient) GetApp(pipelineName string, stageName string, appName string) (*resty.Response, error) {
	res, err := k.client.Get("/api/cli/pipelines/" + pipelineName + "/" + stageName + "/" + appName)

	return res, err
}

func (k *KuberoClient) GetApps() (*resty.Response, error) {
	res, err := k.client.Get("/api/cli/apps")

	return res, err
}

func (k *KuberoClient) GetPipelines() (*resty.Response, error) {
	res, err := k.client.Get("/api/cli/pipelines")

	return res, err
}

func (k *KuberoClient) DeployApp(app AppCRD) (*resty.Response, error) {
	k.client.SetBody(app.Spec)
	res, err := k.client.Post("/api/cli/apps")

	return res, err
}

func (k *KuberoClient) GetPipelineApps(pipelineName string) (*resty.Response, error) {
	res, err := k.client.Get("/api/cli/pipelines/" + pipelineName + "/apps")

	return res, err
}

func (k *KuberoClient) GetAddons() (*resty.Response, error) {
	res, err := k.client.Get("/api/cli/addons")

	return res, err
}

func (k *KuberoClient) GetBuildpacks() (*resty.Response, error) {
	res, err := k.client.Get("/api/cli/config/buildpacks")

	return res, err
}

func (k *KuberoClient) GetPodsize() (*resty.Response, error) {
	res, err := k.client.Get("/api/cli/config/podsize")

	return res, err
}

func (k *KuberoClient) GetRepositories() (*resty.Response, error) {
	res, err := k.client.Get("/api/cli/config/repositories")

	return res, err
}

func (k *KuberoClient) GetContexts() (*resty.Response, error) {
	res, err := k.client.Get("/api/cli/config/k8s/context")

	return res, err
}
