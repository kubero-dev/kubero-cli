package config

import (
	"errors"
	"fmt"
	"github.com/faelmori/kubero-cli/internal/log"
	"github.com/faelmori/kubero-cli/types"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"path/filepath"
)

func (v *ConfigManager) LoadConfig() error {
	path, pathErr := v.GetConfigDir()
	if pathErr != nil {
		return pathErr
	}
	name := v.GetConfigName()
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType(filepath.Ext(name)[1:])
	viper.AutomaticEnv()
	if readInConfigErr := viper.ReadInConfig(); readInConfigErr != nil {
		return readInConfigErr
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Debug("Config file changed:", map[string]interface{}{
			"context": "kubero-cli",
			"pkg":     "config",
			"method":  "LoadConfig",
			"event":   e.Name,
		})
		rErr := viper.ReadInConfig()
		if rErr != nil {
			// Attempt to resolve the config file across multiple directories and methods.
			// If all attempts fail, terminate with a clear error message (critical for app functionality).
			log.Fatal("Error reloading configs!", map[string]interface{}{
				"context": "kubero-cli",
				"pkg":     "config",
				"method":  "LoadConfig",
				"event":   "OnConfigChange",
				"error":   rErr.Error(),
			})
		}
	})

	return nil
}
func (v *ConfigManager) LoadPLConfigs(pipelineName string) (*viper.Viper, error) {
	baseDir := v.GetIACBaseDir()
	dir := baseDir + "/" + pipelineName
	pipelineConfig := viper.New()
	pipelineConfig.SetConfigName("pipeline")
	pipelineConfig.SetConfigType("yaml")
	pipelineConfig.AddConfigPath(dir)
	readInConfigErr := pipelineConfig.ReadInConfig()
	if readInConfigErr != nil {
		log.Error("Failed to read pipeline config", map[string]interface{}{
			"context": "kubero-cli",
			"pkg":     "config",
			"method":  "LoadPLConfigs",
			"error":   readInConfigErr.Error(),
		})
		return nil, readInConfigErr
	}
	return pipelineConfig, nil
}
func (v *ConfigManager) loadCLIConfig() {
	dir := v.getGitDir()
	repoConfig := viper.New()
	repoConfig.SetConfigName("kubero")
	repoConfig.SetConfigType("yaml")
	repoConfig.AddConfigPath(dir)
	repoConfig.ConfigFileUsed()
	errCred := repoConfig.ReadInConfig()

	viper.SetDefault("api.url", "http://default:2000")
	viper.SetConfigName("kubero")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/kubero/")
	viper.AddConfigPath("$HOME/.kubero/")
	err := viper.ReadInConfig()

	if err != nil && errCred != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			fmt.Println("No config file found; using defaults")
		} else {
			fmt.Println("Error while loading config file:", err)
			return
		}
	}

	viperUnmarshalErr := viper.UnmarshalKey("instances", &v.instanceList)
	if viperUnmarshalErr != nil {
		fmt.Println("Error while unmarshalling instances:", viperUnmarshalErr)
		return
	}
	for instanceName, instance := range v.instanceList {
		instance.Name = instanceName
		instance.ConfigPath = viper.ConfigFileUsed()
		v.instanceList[instanceName] = instance
	}

	var repoInstancesList map[string]types.Instance
	unmarshalKeyErr := repoConfig.UnmarshalKey("instances", &repoInstancesList)
	if unmarshalKeyErr != nil {
		fmt.Println("Error while unmarshalling instances:", unmarshalKeyErr)
		return
	}
	for instanceName, repoInstance := range repoInstancesList {
		repoInstance.Name = instanceName
		repoInstance.ConfigPath = repoConfig.ConfigFileUsed()
		v.instanceList[instanceName] = repoInstance
	}

	var instanceNameList = make([]string, 0)
	currentInstanceName := viper.GetString("currentInstance")
	for instanceName, instance := range v.instanceList {
		instance.Name = instanceName
		instanceNameList = append(instanceNameList, instanceName)
		if instanceName == currentInstanceName {
			v.current = instance
		}
	}
}
func (v *ConfigManager) loadCredentials() error {
	if v.credentials == nil {
		v.credentials = viper.New()
	}
	v.credentials.SetConfigName("credentials")
	v.credentials.SetConfigType("yaml")
	v.credentials.AddConfigPath("/etc/kubero/")
	v.credentials.AddConfigPath("$HOME/.kubero/")
	err := v.credentials.ReadInConfig()
	if err != nil {
		log.Error("Failed to read credentials", map[string]interface{}{
			"context": "kubero-cli",
			"pkg":     "config",
			"method":  "loadCredentials",
			"error":   err.Error(),
		})
		return err
	}
	return nil
}
