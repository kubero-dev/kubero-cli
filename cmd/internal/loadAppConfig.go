package internal

import (
	"log"

	"github.com/spf13/viper"
)

func loadAppConfig(pipelineName string, stageName string, appName string) *viper.Viper {

	baseDir := getIACBaseDir()
	dir := baseDir + "/" + pipelineName + "/" + stageName

	appConfig := viper.New()
	appConfig.SetConfigName(appName)
	appConfig.SetConfigType("yaml")
	appConfig.AddConfigPath(dir)
	readInConfigErr := appConfig.ReadInConfig()
	if readInConfigErr != nil {
		log.Fatal(readInConfigErr)
		return nil
	}

	return appConfig
}
