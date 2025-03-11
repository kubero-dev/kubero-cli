package config

import (
	"fmt"
	"github.com/faelmori/kubero-cli/internal/log"
	"github.com/faelmori/kubero-cli/internal/utils"
	"github.com/faelmori/kubero-cli/types"
	logz "github.com/faelmori/logz/logger"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var promptLine = utils.NewConsolePrompt().PromptLine

type ConfigManager struct {
	path, name            string
	logz                  *logz.LogzCore
	Logger                *logz.LogzLogger
	globals, pipelineConf *viper.Viper
	instanceManager       *InstanceManager
	credentialsManager    *CredentialsManager
}

func NewViperConfig(path, name string) IConfigManager {
	var vpCfg ConfigManager
	logger := log.Logger()
	vpCfg = ConfigManager{
		path:    path,
		name:    name,
		globals: viper.New(),
		Logger:  &logger,
	}
	if loadConfigsErr := vpCfg.LoadConfig(); loadConfigsErr != nil {
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

func (v *ConfigManager) WriteCLIConfig(argDomain, argPort, argApiToken string) error {
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

func (v *ConfigManager) GetConfig() types.Config {
	var config types.Config
	if err := v.globals.Unmarshal(&config); err != nil {
		log.Error("Failed to unmarshal config", map[string]interface{}{
			"context": "kubero-cli",
			"pkg":     "config",
			"method":  "GetConfig",
			"error":   err.Error(),
		})
	}
	return config
}
func (v *ConfigManager) GetConfigManager() *ConfigManager { return v }
func (v *ConfigManager) GetLogger() *types.Logger         { return v.Logger }
func (v *ConfigManager) GetViper() *viper.Viper           { return v.globals }
func (v *ConfigManager) GetPath() string                  { return v.path }
func (v *ConfigManager) GetName() string                  { return v.name }

func (v *ConfigManager) GetConfigDir() (string, error) {
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

func (v *ConfigManager) GetConfigName() string {
	if v.name != "" {
		return v.name
	}
	return "config.yaml"
}

func (v *ConfigManager) GetProp(key string) interface{} {
	if v.globals == nil {
		return nil
	}
	return v.globals.Get(key)
}
func (v *ConfigManager) SetProp(key string, value interface{}) {
	if v.globals == nil {
		v.globals = viper.New()
	}
	v.globals.Set(key, value)
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
func (v *ConfigManager) SetPath(path string) error {
	if filepath.IsAbs(path) {
		if filepath.Base(path) == "" {
			v.path = filepath.Dir(path)
			v.name = v.GetConfigName()
		} else {
			v.path = path
			v.name = filepath.Base(path)
		}
	}
	if err := v.LoadConfig(); err != nil {
		return err
	}

	return nil
}
func (v *ConfigManager) SetName(name string) error {
	if name == "" {
		v.name = v.GetConfigName()
	} else {
		v.name = name
	}
	if err := v.LoadConfig(); err != nil {
		return err
	}

	return nil
}

func (v *ConfigManager) GetInstanceManager() *InstanceManager {
	if v.instanceManager == nil {
		if v.credentialsManager == nil {
			v.credentialsManager = NewCredentialsManager()
			if loadCredentialsErr := v.credentialsManager.LoadCredentials(); loadCredentialsErr != nil {
				// Attempt to resolve the config file across multiple directories and methods.
				// If all attempts fail, terminate with a clear error message (critical for app functionality).
				log.Error("Error loading credentials!", map[string]interface{}{
					"context": "kubero-cli",
					"pkg":     "config",
					"method":  "GetInstanceManager",
					"error":   loadCredentialsErr.Error(),
				})
				return nil
			}
		}
		v.instanceManager = NewInstanceManager(v.credentialsManager.GetCredentials())
	}
	return v.instanceManager
}
func (v *ConfigManager) GetCredentialsManager() *CredentialsManager {
	if v.credentialsManager == nil {
		v.credentialsManager = NewCredentialsManager()
		if loadCredentialsErr := v.credentialsManager.LoadCredentials(); loadCredentialsErr != nil {
			// Attempt to resolve the config file across multiple directories and methods.
			// If all attempts fail, terminate with a clear error message (critical for app functionality).
			log.Error("Error loading credentials!", map[string]interface{}{
				"context": "kubero-cli",
				"pkg":     "config",
				"method":  "GetCredentialsManager",
				"error":   loadCredentialsErr.Error(),
			})
			return nil
		}
	}
	return v.credentialsManager
}
func (v *ConfigManager) GetIACBaseDir() string { return v.GetProp("iac.baseDir").(string) }
