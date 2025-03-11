package main

import (
	_ "embed"
	"github.com/faelmori/kubero-cli/cmd/cli/config"
	"github.com/faelmori/kubero-cli/cmd/cli/install"
	"github.com/faelmori/kubero-cli/cmd/cli/pipeline"
	a "github.com/faelmori/kubero-cli/internal/api"
	"github.com/faelmori/kubero-cli/internal/log"
	"github.com/faelmori/logz"
	"github.com/kubero-dev/kubero-cli/cmd/cli"
	"github.com/kubero-dev/kubero-cli/cmd/cli/debug"
	"github.com/kubero-dev/kubero-cli/cmd/cli/instance"
	"github.com/kubero-dev/kubero-cli/cmd/common"
	c "github.com/kubero-dev/kubero-cli/internal/config"
	t "github.com/kubero-dev/kubero-cli/types"
	"github.com/spf13/cobra"
	_ "github.com/spf13/pflag"
	"github.com/spf13/viper"
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
		Annotations: common.GetDescriptions(
			[]string{
				"Kubero is a platform as a service (PaaS) that enables developers to build, run, and operate applications on Kubernetes.",
				"Kubero is a platform as a service (PaaS) that enables developers to build, run, and operate applications on Kubernetes.",
			}, false,
		),
		Example: common.ConcatenateExamples("kubero install", "kubero pipeline create", "kubero net dashboard"),
	}

	rootCmd.AddCommand(config.ConfigCmds()...)
	rootCmd.AddCommand(debug.DebugCmds()...)

	rootCmd.AddCommand(install.InstallCmds()...)

	plRootCmd := &cobra.Command{
		Use:   "pipeline",
		Short: "Pipeline commands",
		Long: `Pipeline commands. Use the pipeline subcommand to manage your pipelines.
Subcommands:
  kubero pipeline [create|fetch|list|down]`,
	}

	plRootCmd.AddCommand(instance.InstanceCmds()...)
	plRootCmd.AddCommand(pipeline.CreateCmds()...)
	plRootCmd.AddCommand(pipeline.FetchPipelineCmds()...)
	plRootCmd.AddCommand(pipeline.PipelineListCmds()...)
	plRootCmd.AddCommand(pipeline.PipelineDownCmds()...)

	rootCmd.AddCommand(plRootCmd)

	netRootCmd := &cobra.Command{
		Use:   "net",
		Short: "Network commands",
		Long: `Network commands. Use the net subcommand to manage your network.
Subcommands:
  kubero net [dashboard|login|tunnel]`,
	}

	netRootCmd.AddCommand(cli.DashboardCmds()...)
	netRootCmd.AddCommand(cli.LoginCmds()...)
	netRootCmd.AddCommand(cli.TunnelCmds()...)

	rootCmd.AddCommand(netRootCmd)

	for _, cmd := range rootCmd.Commands() {
		SetUsageDefinition(cmd, true)
	}

	SetUsageDefinition(rootCmd, false)

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
		log.Debug("Failed to load configuration.", map[string]interface{}{
			"context": "kubero-cli",
			"action":  "initConfig",
			"stage":   "loadConfig",
			"error":   loadErr.Error(),
		})
		//return loadErr
	}
	if loadErr := k.configMgr.GetCredentialsManager().LoadCredentials(); loadErr != nil {
		log.Debug("Failed to load credentials.", map[string]interface{}{
			"context": "kubero-cli",
			"action":  "initConfig",
			"stage":   "loadCredentials",
			"error":   loadErr.Error(),
		})
		//return loadErr
	}
	plName := k.configMgr.GetProp("pipelineName")
	if plName != nil {
		if stPlName, ok := plName.(string); ok && stPlName != "" {
			if _, loadErr := k.configMgr.LoadPLConfigs(stPlName); loadErr != nil {
				log.Debug("Failed to load pipeline configuration.", map[string]interface{}{
					"context": "kubero-cli",
					"action":  "initConfig",
					"stage":   "loadPLConfigs",
					"error":   loadErr.Error(),
				})
				//return loadErr
			}
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
		_ = k.initConfig()
		//k.initAPI()
	})
}
