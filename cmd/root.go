package main

import (
	_ "embed"
	"github.com/faelmori/kubero-cli/cmd/cli"
	"github.com/faelmori/kubero-cli/cmd/cli/config"
	"github.com/faelmori/kubero-cli/cmd/cli/debug"
	"github.com/faelmori/kubero-cli/cmd/cli/install"
	"github.com/faelmori/kubero-cli/cmd/cli/pipeline"
	a "github.com/faelmori/kubero-cli/internal/api"
	c "github.com/faelmori/kubero-cli/internal/config"
	"github.com/faelmori/kubero-cli/internal/log"
	t "github.com/faelmori/kubero-cli/types"
	"github.com/faelmori/logz"
	"github.com/spf13/cobra"
	_ "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"sync"
)

var (
	once sync.Once
)

type KuberoClient struct {
	credentialsConfig   *viper.Viper
	currentInstanceName string
	currentInstance     t.Instance
	configMgr           c.IConfigManager
	rootCmd             *cobra.Command
	log                 logz.Logger
	api                 *a.Client
	db                  *gorm.DB
}

func NewKuberoClient() *KuberoClient {
	kuberoCLI := KuberoClient{}
	kuberoCLI.init()
	return &kuberoCLI
}

func (k *KuberoClient) Command() *cobra.Command {
	if k.rootCmd != nil {
		return k.rootCmd
	}

	var rootCmd = &cobra.Command{
		Use:   "kubero",
		Short: "Kubero is a platform as a service (PaaS) that enables developers to build, run, and operate applications on Kubernetes.",
		Long: `
	,--. ,--.        ,--.
	|  .'   /,--.,--.|  |-.  ,---. ,--.--. ,---.
	|  .   ' |  ||  || .-. '| .-. :|  .--'| .-. |
	|  |\   \'  ''  '| '-' |\   --.|  |   ' '-' '
	'--' '--' '----'  '---'  '----''--'    '---'
Documentation:
  https://docs.kubero.dev
`,
		Example: `kubero install`,
		Aliases: []string{"kbr"},
	}

	rootCmd.AddCommand(install.InstallCmds()...)
	rootCmd.AddCommand(config.ConfigCmds()...)
	rootCmd.AddCommand(debug.DebugCmds()...)
	rootCmd.AddCommand(pipeline.FetchPipelineCmds()...)
	rootCmd.AddCommand(pipeline.PipelineListCmds()...)
	rootCmd.AddCommand(pipeline.PipelineDownCmds()...)
	rootCmd.AddCommand(cli.TunnelCmds()...)

	for _, cmd := range rootCmd.Commands() {
		SetUsageDefinition(cmd)
	}

	rootCmd.CompletionOptions.HiddenDefaultCmd = false

	k.rootCmd = rootCmd

	return rootCmd
}

func (k *KuberoClient) Execute() {
	if err := k.rootCmd.Execute(); err != nil {
		log.Error("Failed to execute the command.", map[string]interface{}{
			"context": "kubero-cli",
			"action":  "Execute",
			"error":   err.Error(),
		})
	}
}

func (k *KuberoClient) initDB() error {
	var err error
	k.db, err = gorm.Open(sqlite.Open("kubero.db"), &gorm.Config{})
	if err != nil {
		log.Error("Failed to connect to local database.", map[string]interface{}{
			"context": "kubero-cli",
			"action":  "initDB",
			"error":   err.Error(),
		})
		return err
	}
	autoMigrateErr := k.db.AutoMigrate(&t.Instance{})
	if autoMigrateErr != nil {
		log.Error("Failed to auto migrate the database.", map[string]interface{}{
			"context": "kubero-cli",
			"action":  "initDB",
			"error":   autoMigrateErr.Error(),
		})
		return autoMigrateErr
	}
	return nil
}

func (k *KuberoClient) initConfig() error {
	name := ""
	path := os.Getenv("KUBERO_CONFIG_PATH")
	if path != "" {
		if isDir := filepath.IsAbs(path); !isDir {
			log.Fatalln("KUBERO_CONFIG_PATH must be an absolute path")
		}
		name = filepath.Base(path)
		path = filepath.Dir(path)
	}

	k.configMgr = c.NewViperConfig(path, name)
	if loadErr := k.configMgr.LoadConfig(); loadErr != nil {
		log.Error("Failed to load configuration.", map[string]interface{}{
			"context": "kubero-cli",
			"action":  "initConfig",
			"error":   loadErr.Error(),
		})
		return loadErr
	}
	if loadErr := k.configMgr.GetCredentialsManager().LoadCredentials(); loadErr != nil {
		log.Error("Failed to load credentials.", map[string]interface{}{
			"context": "kubero-cli",
			"action":  "initConfig",
			"error":   loadErr.Error(),
		})
		return loadErr
	}
	plName := k.configMgr.GetProp("pipelineName").(string)
	if plName != "" {
		if _, loadErr := k.configMgr.LoadPLConfigs(plName); loadErr != nil {
			log.Error("Failed to load pipeline configuration.", map[string]interface{}{
				"context": "kubero-cli",
				"action":  "initConfig",
				"error":   loadErr.Error(),
			})
			return loadErr
		}
	}
	if k.configMgr.GetProp("instanceName") != nil {
		k.currentInstanceName = k.configMgr.GetProp("instanceName").(string)
		if k.currentInstanceName != "" {
			instance := k.configMgr.GetInstanceManager().GetCurrentInstance()
			k.currentInstance = *instance
		}
	}

	return nil
}

func (k *KuberoClient) init() {
	once.Do(func() {
		k.log = log.Logger()
		k.rootCmd = k.Command()
		_ = k.initDB()
		_ = k.initConfig()
		//k.initAPI()
	})
}
