package kuberoCli

import (
	"github.com/AlecAivazis/survey/v2"
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
	upCmd.Flags().StringVarP(&stageName, "stage", "s", "", "Name of the stage [test|stage|production]")
	upCmd.Flags().StringVarP(&appName, "app", "a", "", "name of the app")
	upCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Skip asking for confirmation")
}

func upPipeline() {

	pipelinesList := listAllLocalPipelines()
	if pipelineName == "" {
		prompt := &survey.Select{
			Message: "Select a pipeline",
			Options: pipelinesList,
		}
		survey.AskOne(prompt, &pipelineName)
	}

	confirmationLine("Are you sure you want to deploy the pipeline '"+pipelineName+"'?", "y")

	pipeline := loadLocalPipeline(pipelineName)
	api.DeployPipeline(pipeline)
}

func upApp() {

	if pipelineName == "" {
		cfmt.Println("{{Please specify a pipeline}}::red")
		return
	}

	confirmationLine("Are you sure you want to deploy the app "+appName+" to "+pipelineName+"?", "y")
	// TODO: implement app deployment
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
