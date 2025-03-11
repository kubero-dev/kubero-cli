package config

// IManagerConfig defines the interface for loading configuration
type IManagerConfig interface {
	LoadConfigs() error
	GetConfigDir() (string, error)
	GetConfigName() string

	GetProp(key string) interface{}
	SetProp(key string, value interface{})

	saveConfig() error
}
