package internal

import (
	"log"

	"github.com/spf13/viper"
	"kubero/pkg/kuberoApi"
)

func loadLocalApp(pipelineName string, stageName string, appName string) kuberoApi.AppCRD {

	appConfig := loadAppConfig(pipelineName, stageName, appName)

	var appCRD kuberoApi.AppCRD

	appConfigUnmarshalErr := appConfig.Unmarshal(&appCRD)
	if appConfigUnmarshalErr != nil {
		log.Fatal(appConfigUnmarshalErr)
		return kuberoApi.AppCRD{}
	}

	return appCRD
}
