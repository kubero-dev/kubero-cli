package kuberoApi

import (
	_ "embed"

	"github.com/go-resty/resty/v2"
)

//go:embed VERSION
var version string

func InitClient(baseURL string, bearerToken string) *resty.Request {
	client := resty.New().SetBaseURL(baseURL).R().
		EnableTrace().
		SetAuthScheme("Bearer").
		SetAuthToken(bearerToken).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "kubero-cli/"+version)

	return client
}
