package kuberoCli

import (
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Deploy your pipelines and apps to the cluster",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if pipelineName != "" && appName == "" {
			upPipeline()
		} else if appName != "" {
			upApp()
		} else {
			upAllPipelines()
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	upCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
	upCmd.Flags().StringVarP(&appName, "app", "a", "", "name of the app")
}

func upPipeline() {
	confirmation := promptLine("Are you sure you want to deploy the pipeline "+pipelineName+"?", "[y,n]", "y")
	if confirmation != "y" {
		cfmt.Println("{{Aborted}}::red")
	}

	cfmt.Println("{{Deploying pipeline}}::yellow " + pipelineName)

	pipeline := loadLocalPipeline(pipelineName)
	api.DeployPipeline(pipeline)
}

func upApp() {

	if pipelineName == "" {
		cfmt.Println("{{Please specify a pipeline}}::red")
		return
	}

	confirmation := promptLine("Are you sure you want to deploy the app "+appName+" to "+pipelineName+"?", "[y,n]", "y")
	if confirmation == "y" {
		cfmt.Println("{{Deploying app}}::yellow " + appName + " to " + pipelineName)
	} else {
		cfmt.Println("{{Aborted}}::red")
		return
	}
}

func upAllPipelines() {

	confirmation := promptLine("Are you sure you want to deploy all pipelines?", "[y,n]", "n")
	if confirmation != "y" {
		cfmt.Println("{{Aborted}}::red")
		return
	}
	pipelinesConfigs := loadAllLocalPipelines()

	cfmt.Println("{{Deploying all pipelines}}::yellow")
	//iterate over pipelinesConfigs
	for _, pipeline := range pipelinesConfigs {
		cfmt.Println("{{Deploying pipeline}}::yellow " + pipeline.Spec.Name + "")
		api.DeployPipeline(pipeline)
	}
}
