package kuberoCli

import (
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:     "up",
	Aliases: []string{"deploy", "dp"},
	Short:   "Deploy your pipelines and apps to the cluster",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		if pipelineName != "" && appName == "" {

			pipelinesList := getAllLocalPipelines()
			ensurePipelineIsSet(pipelinesList)
			upPipeline()
		} else if appName != "" {

			pipelinesList := getAllLocalPipelines()
			ensurePipelineIsSet(pipelinesList)
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
	confirmationLine("Are you sure you want to deploy the pipeline '"+pipelineName+"'?", "y")

	pipeline := loadLocalPipeline(pipelineName)
	api.DeployPipeline(pipeline)
}

func upApp() {
	confirmationLine("Are you sure you want to deploy the app "+appName+" to "+pipelineName+"?", "y")
	app := loadLocalApp(pipelineName, stageName, appName)
	api.DeployApp(app)
}

func upAllPipelines() {

	confirmationLine("Are you sure you want to deploy all pipelines?", "y")
	pipelinesConfigs := loadAllLocalPipelines()

	cfmt.Println("{{Deploying all pipelines}}::yellow")
	//iterate over pipelinesConfigs
	for _, pipelineCRD := range pipelinesConfigs {
		cfmt.Println("{{Deploying pipeline}}::yellow " + pipelineCRD.Spec.Name + "")
		api.DeployPipeline(pipelineCRD)
	}
}
