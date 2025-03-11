package pipeline

import (
	"github.com/kubero-dev/kubero-cli/internal/pipeline"
	"github.com/spf13/cobra"
)

func PipelineDownCmds() []*cobra.Command {
	plDownRootCmd := &cobra.Command{
		Use:   "down",
		Short: "Undeploy your pipelines and apps from the cluster",
		Long: `Use the pipeline or app subcommand to undeploy your pipelines and apps from the cluster
Subcommands:
  kubero down [pipeline|app]`,
	}
	plDownRootCmd.AddCommand(cmdDownPL())
	plDownRootCmd.AddCommand(cmdDownAppPL())
	plDownRootCmd.AddCommand(cmdDownPipelinePL())

	return []*cobra.Command{
		plDownRootCmd,
	}
}

func cmdDownPL() *cobra.Command {
	var pipelineName, stageName, appName string
	var force bool

	var downCmd = &cobra.Command{
		Use:     "down",
		Aliases: []string{"undeploy", "dn"},
		Short:   "Undeploy your pipelines and apps from the cluster",
		Long: `Use the pipeline or app subcommand to undeploy your pipelines and apps from the cluster
Subcommands:
  kubero down [pipeline|app]`,
		Run: func(cmd *cobra.Command, args []string) {
			pl := pipeline.NewPipelineManager(pipelineName, stageName, appName)
			if pipelineName != "" && appName == "" {
				pl.DownPipeline()
			} else if appName != "" {
				pl.DownApp()
			} else {
				pl.DownAllPipelines()
			}
		},
	}

	downCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
	downCmd.Flags().StringVarP(&stageName, "stage", "s", "", "Name of the stage [test|stage|production]")
	downCmd.Flags().StringVarP(&appName, "app", "a", "", "name of the app")
	downCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Skip asking for confirmation")

	return downCmd
}

func cmdDownAppPL() *cobra.Command {
	var pipelineName, stageName, appName string
	var force bool

	var downAppCmd = &cobra.Command{
		Use:   "app",
		Short: "Undeploy an apps from the cluster",
		Long:  `Use the app subcommand to undeploy your apps from the cluster`,
		Run: func(cmd *cobra.Command, args []string) {
			pl := pipeline.NewPipelineManager(pipelineName, stageName, appName)

			pl.DownApp()
		},
	}

	downAppCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
	downAppCmd.Flags().StringVarP(&stageName, "stage", "s", "", "Name of the stage [test|stage|production]")
	downAppCmd.Flags().StringVarP(&appName, "app", "a", "", "name of the app")
	downAppCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Skip asking for confirmation")

	return downAppCmd
}

func cmdDownPipelinePL() *cobra.Command {
	var pipelineName, stageName, appName string
	var force bool

	var downPipelineCmd = &cobra.Command{
		Use:     "pipeline",
		Aliases: []string{"pl"},
		Short:   "Undeploy a pipeline from the cluster",
		Long:    `Use the pipeline subcommand to undeploy your pipelines from the cluster`,
		Run: func(cmd *cobra.Command, args []string) {
			pl := pipeline.NewPipelineManager(pipelineName, stageName, appName)
			pl.DownPipeline()
		},
	}

	downPipelineCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
	downPipelineCmd.Flags().StringVarP(&stageName, "stage", "s", "", "Name of the stage [test|stage|production]")
	downPipelineCmd.Flags().StringVarP(&appName, "app", "a", "", "name of the app")
	downPipelineCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Skip asking for confirmation")

	return downPipelineCmd
}
