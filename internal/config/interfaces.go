package config

import (
	"github.com/faelmori/kubero-cli/types"
	"github.com/spf13/viper"
)

// IConfigManager defines the interface for loading configuration
type IConfigManager interface {
	LoadConfig() error
	LoadPLConfigs(pipelineName string) (*viper.Viper, error)
	GetConfigDir() (string, error)
	GetConfigName() string

	GetConfig() types.Config
	GetLogger() *types.Logger
	GetViper() *types.Viper
	getLogz() *types.LogzCore
	GetName() string
	GetPath() string

	SetPath(path string) error
	SetName(name string) error
	GetProp(key string) interface{}
	SetProp(key string, value interface{})
	saveConfig() error

	GetGitDir() string
	GetGitRemote() string
	GetIACBaseDir() string

	WriteCLIConfig(argDomain, argPort, argToken string) error
	loadCLIConfig()

	GetInstanceManager() *InstanceManager
	GetCredentialsManager() *CredentialsManager
}
