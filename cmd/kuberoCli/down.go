package kuberoCli

import (
	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Undeploy your pipelines and apps from the cluster",
	Long: `Use the pipeline or app subcommand to undeploy your pipelines and apps from the cluster
Subcommands:
  kubero down [pipeline|app]`,
	Run: func(cmd *cobra.Command, args []string) {
		if pipelineName != "" && appName == "" {
			downPipeline()
		} else if appName != "" {
			downApp()
		} else {
			downAllPipelines()
		}
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
	downCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
	downCmd.Flags().StringVarP(&stageName, "stage", "s", "", "Name of the stage [test|stage|production]")
	downCmd.Flags().StringVarP(&appName, "app", "a", "", "name of the app")
	downCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Skip asking for confirmation")
}

func downPipeline() {
	ensurePipelineIsSet()
	downPipelineByName(pipelineName)
}

func downPipelineByName(pipelineName string) {
	confirmationLine("Are you sure you want to undeploy the pipeline '"+pipelineName+"'?", "y")

	_, err := api.UnDeployPipeline(pipelineName)
	if err != nil {
		panic("Unable to undeploy Pipeline")
	}
}

func downApp() {

	ensurePipelineIsSet()
	ensurePipelineIsSet()
	ensureStageNameIsSet()

	confirmationLine("Are you sure you want to undeploy the app "+appName+" from "+stageName+" in "+pipelineName+"?", "y")

	_, err := api.UnDeployApp(pipelineName, stageName, appName)
	if err != nil {
		panic("Unable to undeploy App")
	}
}

func downAllPipelines() {
	confirmationLine("Are you sure you want to undeploy all pipelines?", "y")
	pipelinesList := getAllLocalPipelines()
	for _, pipeline := range pipelinesList {
		downPipelineByName(pipeline)
	}
}
