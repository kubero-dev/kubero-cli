package config

import (
	"fmt"
	"github.com/faelmori/kubero-cli/internal/log"
	"github.com/faelmori/kubero-cli/internal/utils"
	"github.com/faelmori/kubero-cli/types"
	logz "github.com/faelmori/logz/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var promptLine = utils.NewConsolePrompt().PromptLine

type ManagerConfig struct {
	path, name string
	Viper      *viper.Viper
	Logz       *logz.LogzCore
	Logger     *logz.LogzLogger
}

func NewViperConfig(path, name string) *ManagerConfig {
	var vpCfg ManagerConfig
	logger := log.Logger()
	vpCfg = ManagerConfig{
		path:   path,
		name:   name,
		Viper:  viper.New(),
		Logger: &logger,
	}
	if loadConfigsErr := vpCfg.LoadConfigs(); loadConfigsErr != nil {
		// Attempt to resolve the config file across multiple directories and methods.
		// If all attempts fail, terminate with a clear error message (critical for app functionality).
		log.Fatal("Error reloading configs!", map[string]interface{}{
			"context": "kubero-cli",
			"pkg":     "config",
			"method":  "NewViperConfig",
			"error":   loadConfigsErr.Error(),
		})
	}
	return &vpCfg
}

func (v *ManagerConfig) WriteCLIConfig(argDomain, argPort, argApiToken string) error {
	ingressInstall := promptLine("10) Write the Kubero CLI config", "[y,n]", "n")
	if ingressInstall != "y" {
		log.Info("Skipping Kubero CLI config")
		return nil
	}

	//TODO consider using SSL here.
	url := promptLine("Kubero Host address", "", "http://"+argDomain+":"+argPort)
	viper.Set("api.url", url)

	token := promptLine("Kubero Token", "", argApiToken)
	viper.Set("api.token", token)

	var config types.Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Error("Failed to unmarshal config", map[string]interface{}{
			"context": "kubero-cli",
			"pkg":     "config",
			"method":  "WriteCLIConfig",
			"error":   err.Error(),
		})
		return err
	}
	fmt.Printf("%+v\n", config)

	if viperWrtErr := viper.WriteConfig(); viperWrtErr != nil {
		log.Error("Failed to write config", map[string]interface{}{
			"context": "kubero-cli",
			"pkg":     "config",
			"method":  "WriteCLIConfig",
			"error":   viperWrtErr.Error(),
		})
		return viperWrtErr
	}

	return nil
}
func (v *ManagerConfig) saveConfig() error {
	if v.Viper == nil {
		return nil
	}
	if v.path == "" {
		v.path, _ = v.GetConfigDir()
	}
	if v.name == "" {
		v.name = v.GetConfigName()
	}
	return v.Viper.WriteConfigAs(filepath.Join(v.path, v.name))
}
func (v *ManagerConfig) GetConfigDir() (string, error) {
	if v.path != "" {
		return v.path, nil
	}
	var (
		homeDirEnv                  = os.Getenv("HOME")
		homeDirUser, homeDirUserErr = os.UserHomeDir()
		cacheDir, cacheDirErr       = os.UserCacheDir()
		configDir, configDirErr     = os.UserConfigDir()
	)
	if homeDirUserErr != nil {
		if homeDirEnv == "" {
			if configDirErr != nil {
				if cacheDirErr != nil {
					return "", homeDirUserErr
				} else {
					homeDirUser = cacheDir
				}
			} else {
				homeDirUser = configDir
			}
		} else {
			homeDirUser = homeDirEnv
		}
	} else if homeDirEnv != homeDirUser {
		homeDirUser = homeDirEnv
	}
	if homeDirUser == "" {
		return "", homeDirUserErr
	} else {
		homeDirEnv = homeDirUser
	}
	homeDir := homeDirUser
	var configPath = filepath.Join(homeDir, "/.config/kubero-cli")
	mkdirAllErr := os.MkdirAll(configPath, 0755)
	if mkdirAllErr != nil || !os.IsExist(mkdirAllErr) {
		return "", mkdirAllErr
	}
	v.path = configPath
	return configPath, nil
}
func (v *ManagerConfig) GetConfigName() string {
	if v.name != "" {
		return v.name
	}
	return "config.yaml"
}
func (v *ManagerConfig) LoadConfigs() error {
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
			"method":  "LoadConfigs",
			"event":   e.Name,
		})
		rErr := viper.ReadInConfig()
		if rErr != nil {
			// Attempt to resolve the config file across multiple directories and methods.
			// If all attempts fail, terminate with a clear error message (critical for app functionality).
			log.Fatal("Error reloading configs!", map[string]interface{}{
				"context": "kubero-cli",
				"pkg":     "config",
				"method":  "LoadConfigs",
				"event":   "OnConfigChange",
				"error":   rErr.Error(),
			})
		}
	})
	return nil
}
func (v *ManagerConfig) GetProp(key string) interface{} {
	if v.Viper == nil {
		return nil
	}
	return v.Viper.Get(key)
}
func (v *ManagerConfig) SetProp(key string, value interface{}) {
	if v.Viper == nil {
		v.Viper = viper.New()
	}
	v.Viper.Set(key, value)
	go func() {
		saveConfigErr := v.saveConfig()
		if saveConfigErr != nil {
			// Attempt to resolve the config file across multiple directories and methods.
			// If all attempts fail, terminate with a clear error message (critical for app functionality).
			log.Fatal("Error saving config!", map[string]interface{}{
				"context": "kubero-cli",
				"pkg":     "config",
				"method":  "SetProp",
				"error":   saveConfigErr.Error(),
			})
		}
	}()
}
func (v *ManagerConfig) GetCurrentInstance() types.Instance {
	var instance types.Instance
	if err := v.Viper.UnmarshalKey("instance", &instance); err != nil {
		log.Error("Failed to unmarshal instance", map[string]interface{}{
			"context": "kubero-cli",
			"pkg":     "config",
			"method":  "GetCurrentInstance",
			"error":   err.Error(),
		})
	}
	return instance
}
func (v *ManagerConfig) GetConfig() types.Config {
	var config types.Config
	if err := v.Viper.Unmarshal(&config); err != nil {
		log.Error("Failed to unmarshal config", map[string]interface{}{
			"context": "kubero-cli",
			"pkg":     "config",
			"method":  "GetConfig",
			"error":   err.Error(),
		})
	}
	return config
}
func (v *ManagerConfig) GetLogger() *logz.LogzLogger { return v.Logger }
func (v *ManagerConfig) getLogz() *logz.LogzCore     { return v.Logz }
func (v *ManagerConfig) GetViper() *viper.Viper      { return v.Viper }
func (v *ManagerConfig) GetPath() string             { return v.path }
func (v *ManagerConfig) GetName() string             { return v.name }
func (v *ManagerConfig) SetPath(path string) error {
	if filepath.IsAbs(path) {
		if filepath.Base(path) == "" {
			v.path = filepath.Dir(path)
			v.name = v.GetConfigName()
		} else {
			v.path = path
			v.name = filepath.Base(path)
		}
	}
	if err := v.LoadConfigs(); err != nil {
		return err
	}

	return nil
}
func (v *ManagerConfig) SetName(name string) error {
	if name == "" {
		v.name = v.GetConfigName()
	} else {
		v.name = name
	}
	if err := v.LoadConfigs(); err != nil {
		return err
	}

	return nil
}
func (v *ManagerConfig) setViper(vp *viper.Viper) {
	if vp == nil {
		if v.Viper == nil {
			v.Viper = viper.New()
		} else {
			viper.Reset()
			v.Viper = viper.GetViper()
		}
	} else {
		v.Viper = vp
	}
}
