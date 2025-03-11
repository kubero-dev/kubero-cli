package config

import (
	"github.com/kubero-dev/kubero-cli/internal/log"
	"github.com/spf13/viper"
	"os"
)

type CredentialsManager struct {
	credentialsCfg *viper.Viper
}

func NewCredentialsManager() *CredentialsManager {
	return &CredentialsManager{
		credentialsCfg: viper.New(),
	}
}

func (c *CredentialsManager) GetCredentials() *viper.Viper    { return c.credentialsCfg }
func (c *CredentialsManager) SetCredentials(cfg *viper.Viper) { c.credentialsCfg = cfg }
func (c *CredentialsManager) LoadCredentials() error {
	if c == nil {
		nc := NewCredentialsManager()
		if nc.credentialsCfg == nil {
			nc.credentialsCfg = viper.New()
		}
		c = nc
	}
	c.credentialsCfg.SetConfigName("credentials")
	c.credentialsCfg.SetConfigType("yaml")
	c.credentialsCfg.AddConfigPath("/etc/kubero/")
	c.credentialsCfg.AddConfigPath("$HOME/.kubero/")

	if cfgInfo, cfgStat := os.Stat(c.credentialsCfg.ConfigFileUsed()); cfgStat != nil || cfgInfo.Size() == 0 {
		log.Debug("No credentials found", map[string]interface{}{
			"context": "kubero-cli",
			"pkg":     "config",
			"method":  "loadCredentials",
		})
		return nil
	}

	err := c.credentialsCfg.ReadInConfig()
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
func (c *CredentialsManager) WriteCredentials() error {
	if c.credentialsCfg == nil {
		c.credentialsCfg = viper.New()
	}
	c.credentialsCfg.SetConfigName("credentials")
	c.credentialsCfg.SetConfigType("yaml")
	c.credentialsCfg.AddConfigPath("/etc/kubero/")
	c.credentialsCfg.AddConfigPath("$HOME/.kubero/")
	if err := c.credentialsCfg.WriteConfig(); err != nil {
		log.Error("Failed to write credentials", map[string]interface{}{
			"context": "kubero-cli",
			"pkg":     "config",
			"method":  "WriteCredentials",
			"error":   err.Error(),
		})
		return err
	}
	return nil
}
func (c *CredentialsManager) GetCredentialsDir() (string, error) {
	if c.credentialsCfg.ConfigFileUsed() == "" {
		return "", nil
	}
	return c.credentialsCfg.ConfigFileUsed(), nil
}
