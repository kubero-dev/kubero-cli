package kuberoCli

import (
	"github.com/i582/cfmt/cmd/cfmt"
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
	downCmd.Flags().StringVarP(&appName, "app", "a", "", "name of the app")
}

func downPipeline() {
	confirmation := promptLine("Are you sure you want to undeploy the pipeline "+pipelineName+"?", "[y,n]", "y")
	if confirmation == "y" {
		cfmt.Println("{{Undeploying pipeline}} " + pipelineName + "::yellow")
	} else {
		cfmt.Println("{{Aborted}}::red")
		return
	}
}

func downApp() {

	if pipelineName == "" {
		cfmt.Println("{{Please specify a pipeline}}::red")
		return
	}

	confirmation := promptLine("Are you sure you want to undeploy the app "+appName+"?", "[y,n]", "y")
	if confirmation == "y" {
		cfmt.Println("{{Undeploying app}} " + appName + "::yellow")
	} else {
		cfmt.Println("{{Aborted}}::red")
		return
	}
}

func downAllPipelines() {
	confirmation := promptLine("Are you sure you want to undeploy all pipelines?", "[y,n]", "n")
	if confirmation == "y" {
		cfmt.Println("{{Undeploying all pipelines}}::yellow")
	} else {
		cfmt.Println("{{Aborted}}::red")
		return
	}
}
