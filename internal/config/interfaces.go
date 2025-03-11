package config

import "github.com/faelmori/kubero-cli/types"

// IManagerConfig defines the interface for loading configuration
type IManagerConfig interface {
	LoadConfigs() error
	GetConfigDir() (string, error)
	GetConfigName() string
	WriteCLIConfig(argDomain, argPort, argToken string) error
	GetConfig() types.Config
	GetCurrentInstance() types.Instance
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
}
