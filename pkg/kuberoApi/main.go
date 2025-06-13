package kuberoApi

import (
	_ "embed"
	"fmt"
	"net/url"
	"os"
	"regexp"

	"github.com/golang-jwt/jwt/v4"

	"github.com/go-resty/resty/v2"
)

type KuberoClient struct {
	baseURL     string
	bearerToken string
	host        string
	client      *resty.Request
}

type ClientNotInitializedError struct{}
type BaseURLNotSetError struct{}
type BearerTokenNotSetError struct{}
type NotAuthenticatedError struct{}

func (e *BaseURLNotSetError) Error() string {
	return "base URL not set"
}

func (e *ClientNotInitializedError) Error() string {
	return "client not initialized"
}

func (e *BearerTokenNotSetError) Error() string {
	return "bearer token not set"
}

func (e *NotAuthenticatedError) Error() string {
	return "not authenticated"
}

//go:embed VERSION
var version string

func (k *KuberoClient) Init(baseURL string, bearerToken string) *resty.Request {
	k.SetApiUrl(baseURL, bearerToken)

	return k.client
}

func (k *KuberoClient) validateClient() error {
	auth, _ := k.checkAuth()
	if !auth {
		fmt.Println("Error: Not authenticated. Run 'kubero login' to authenticate")
		os.Exit(0)
		return &NotAuthenticatedError{}
	}

	if k.client == nil {
		fmt.Println("Error: Client not initialized. Run 'kubero login' to authenticate")
		os.Exit(0)
		return &ClientNotInitializedError{}
	}

	if k.baseURL == "" {
		fmt.Println("Error: Base URL not set. Run 'kubero login' to authenticate")
		os.Exit(0)
		return &BaseURLNotSetError{}
	}

	if k.bearerToken == "" {
		return &BearerTokenNotSetError{}
	} else if !k.validateToken() {
		return &NotAuthenticatedError{}
	}

	return nil
}

func (k *KuberoClient) SetApiUrl(apiUrl string, bearerToken string) {

	parsedUrl, err := url.Parse(apiUrl)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	// resty needs to resolve the url. kubero.localhost will not be resolved
	// so we need to set the host header to the correct value
	matched, _ := regexp.MatchString(`localhost`, parsedUrl.Host)
	if matched {
		k.host = "kubero.localhost"
		k.baseURL = parsedUrl.Scheme + "://localhost:" + parsedUrl.Port()
	} else {
		k.baseURL = apiUrl
		k.host = parsedUrl.Host
	}

	k.client = resty.New().SetBaseURL(k.baseURL).R().
		EnableTrace().
		SetAuthScheme("Bearer").
		SetAuthToken(bearerToken).
		SetHeader("Host", k.host).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "kubero-cli/"+version)

	k.bearerToken = bearerToken

}

func (k *KuberoClient) validateToken() bool {
	if k.bearerToken == "" {
		return true
	}
	token, _, err := new(jwt.Parser).ParseUnverified(k.bearerToken, jwt.MapClaims{})
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			if exp < 0 {
				fmt.Println("Token expired")
				return false
			}
		}
	}

	return true
}

func (k *KuberoClient) checkAuth() (bool, error) {
	res, err := k.client.Get("/api/auth/session")
	if err != nil {
		panic(err)
	}
	if res.StatusCode() > 299 {
		return false, &NotAuthenticatedError{}
	}
	return true, nil
}

func (k *KuberoClient) DeployPipeline(pipeline PipelineCRD) (*resty.Response, error) {
	k.client.SetBody(pipeline.Spec)
	res, err := k.client.Post("/api/pipelines/")

	return res, err
}

func (k *KuberoClient) DeletePipeline(pipelineName string) (*resty.Response, error) {
	res, err := k.client.Delete("/api/pipelines/" + pipelineName)

	return res, err
}

func (k *KuberoClient) GetPipeline(pipelineName string) (*resty.Response, error) {
	res, err := k.client.Get("/api/pipelines/" + pipelineName)

	return res, err
}

func (k *KuberoClient) DeleteApp(pipelineName string, stageName string, appName string) (*resty.Response, error) {
	res, err := k.client.Delete("/api/pipelines/" + pipelineName + "/" + stageName + "/" + appName)

	return res, err
}

func (k *KuberoClient) GetApp(pipelineName string, stageName string, appName string) (*resty.Response, error) {
	res, err := k.client.Get("/api/pipelines/" + pipelineName + "/" + stageName + "/" + appName)

	return res, err
}

func (k *KuberoClient) GetApps() (*resty.Response, error) {
	res, err := k.client.Get("/api/apps")

	return res, err
}

func (k *KuberoClient) GetPipelines() (*resty.Response, error) {
	k.validateClient()
	res, err := k.client.Get("/api/pipelines")

	return res, err
}

func (k *KuberoClient) DeployApp(pipelineName string, phaseName string, appName string, app AppCRD) (*resty.Response, error) {
	k.client.SetBody(app.Spec)
	res, err := k.client.Post("/api/apps/" + pipelineName + "/" + phaseName + "/" + appName)

	return res, err
}

func (k *KuberoClient) GetPipelineApps(pipelineName string) (*resty.Response, error) {
	res, err := k.client.Get("/api/pipelines/" + pipelineName + "/apps")

	return res, err
}

func (k *KuberoClient) GetAddons() (*resty.Response, error) {
	k.validateClient()
	res, err := k.client.Get("/api/addons")

	return res, err
}

func (k *KuberoClient) GetRunpacks() (*resty.Response, error) {
	k.validateClient()
	res, err := k.client.Get("/api/config/runpacks")

	return res, err
}

func (k *KuberoClient) GetPodsize() (*resty.Response, error) {
	k.validateClient()
	res, err := k.client.Get("/api/config/podsizes")

	return res, err
}

func (k *KuberoClient) GetRepositories() (*resty.Response, error) {
	res, err := k.client.Get("/api/config/repositories")

	return res, err
}

func (k *KuberoClient) GetContexts() (*resty.Response, error) {
	res, err := k.client.Get("/api/kubernetes/context")

	return res, err
}

func (k *KuberoClient) GetEvents() (*resty.Response, error) {
	k.client.QueryParam.Add("namespace", "")
	res, err := k.client.Get("/api/kubernetes/namespace")

	return res, err
}

func (k *KuberoClient) GetLogs(pipelineName string, phaseName string, appName string, container string) (*resty.Response, error) {
	res, err := k.client.Get("/api/logs/" + pipelineName + "/" + phaseName + "/" + appName + "/" + container + "/history")

	return res, err
}

func (k *KuberoClient) Login(user string, pass string) (*resty.Response, error) {

	k.client.SetBody(map[string]string{"username": user, "password": pass})
	res, err := k.client.Post("/api/auth/login")

	return res, err
}
