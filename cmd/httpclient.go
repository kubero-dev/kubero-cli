package cmd

import (
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

var client *resty.Request

func InitClient() {
	client = resty.New().SetBaseURL(viper.GetString("api.url")).R().
		EnableTrace().
		SetAuthScheme("Bearer").
		SetAuthToken(viper.GetString("api.token")).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "kubero-cli/0.0.1")
}
