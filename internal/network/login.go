package network

import (
	"fmt"
	c "github.com/faelmori/kubero-cli/internal/config"
	"github.com/spf13/viper"
)

type Login struct {
	credentialsCfg  *viper.Viper
	instanceManager *c.InstanceManager
}

func NewLogin(cfg c.IConfigManager) *Login {
	if cfg == nil {
		cfg = c.NewViperConfig("", "")
	}
	if loadCredErr := cfg.GetCredentialsManager().LoadCredentials(); loadCredErr != nil {
		fmt.Println("Error loading credentials: ", loadCredErr)
		return nil
	}
	instanceManager := c.NewInstanceManager(cfg.GetCredentialsManager().GetCredentials())
	credentialsCfg := cfg.GetProp("credentials")
	if credentialsCfg == nil {
		credentialsCfg = viper.New()
		cfg.SetProp("credentials", credentialsCfg)
	}

	return &Login{
		credentialsCfg:  credentialsCfg.(*viper.Viper),
		instanceManager: instanceManager,
	}
}

func (l *Login) EnsureInstanceOrCreate() error {
	cfg := c.NewViperConfig("", "")
	if loadCredErr := cfg.GetCredentialsManager().LoadCredentials(); loadCredErr != nil {
		fmt.Println("Error loading credentials: ", loadCredErr)
		return loadCredErr
	}

	instanceNameList := l.instanceManager.GetInstanceNameList()
	instanceName := selectFromList("Select an instance", instanceNameList, l.instanceManager.GetCurrentInstance().Name)
	instance := l.instanceManager.GetInstance(instanceName)
	if instance.ApiUrl == "" {
		if createInstanceErr := l.instanceManager.CreateInstanceForm(); createInstanceErr != nil {
			return createInstanceErr
		}
	} else {
		if setCurInstanceErr := l.instanceManager.SetCurrentInstance(instanceName); setCurInstanceErr != nil {
			fmt.Println("Error setting current instance: ", setCurInstanceErr)
			return setCurInstanceErr
		}
	}

	return nil
}

func (l *Login) SetKuberoCredentials(token string) error {
	if token == "" {
		token = promptLine("Kubero Token", "", "")
	}

	l.credentialsCfg.Set(l.instanceManager.GetCurrentInstance().Name, token)
	writeConfigErr := l.credentialsCfg.WriteConfig()
	if writeConfigErr != nil {
		fmt.Println("Error writing config file: ", writeConfigErr)
		return writeConfigErr
	}

	return nil
}
