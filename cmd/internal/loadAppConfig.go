package internal

import (
	"log"

	"github.com/spf13/viper"
)

func loadAppConfig(pipelineName string, stageName string, appName string) *viper.Viper {

	dir := ".kubero/" + pipelineName + "/" + stageName

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
