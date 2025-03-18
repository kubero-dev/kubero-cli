package config

import (
	l "github.com/kubero-dev/kubero-cli/internal/log"
	v "github.com/spf13/viper"
	"os"
)

type CredentialsManager struct {
	credentialsCfg *v.Viper
}

func NewCredentialsManager() *CredentialsManager {
	return &CredentialsManager{
		credentialsCfg: v.New(),
	}
}

func (c *CredentialsManager) GetCredentials() *v.Viper    { return c.credentialsCfg }
func (c *CredentialsManager) SetCredentials(cfg *v.Viper) { c.credentialsCfg = cfg }
func (c *CredentialsManager) LoadCredentials() error {
	if c == nil {
		nc := NewCredentialsManager()
		if nc.credentialsCfg == nil {
			nc.credentialsCfg = v.New()
		}
		c = nc
	}
	c.credentialsCfg.SetConfigName("credentials")
	c.credentialsCfg.SetConfigType("yaml")
	c.credentialsCfg.AddConfigPath("/etc/kubero/")
	c.credentialsCfg.AddConfigPath("$HOME/.kubero/")

	if cfgInfo, cfgStat := os.Stat(c.credentialsCfg.ConfigFileUsed()); cfgStat != nil || cfgInfo.Size() == 0 {
		l.Debug("No credentials found", map[string]interface{}{
			"context": "kubero-cli",
			"pkg":     "config",
			"method":  "loadCredentials",
		})
		return nil
	}

	err := c.credentialsCfg.ReadInConfig()
	if err != nil {
		l.Error("Failed to read credentials", map[string]interface{}{
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
		c.credentialsCfg = v.New()
	}
	c.credentialsCfg.SetConfigName("credentials")
	c.credentialsCfg.SetConfigType("yaml")
	c.credentialsCfg.AddConfigPath("/etc/kubero/")
	c.credentialsCfg.AddConfigPath("$HOME/.kubero/")
	if err := c.credentialsCfg.WriteConfig(); err != nil {
		l.Error("Failed to write credentials", map[string]interface{}{
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
