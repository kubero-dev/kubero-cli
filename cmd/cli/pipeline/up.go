package pipeline

import (
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/cobra"
	"os"
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

var upAppCmd = &cobra.Command{
	Use:   "app",
	Short: "Deploy an apps to the cluster",
	Long:  `Use the app subcommand to deploy your apps to the cluster`,
	Run: func(cmd *cobra.Command, args []string) {

		pipelinesList := getAllLocalPipelines()
		ensurePipelineIsSet(pipelinesList)

		ensureStageNameIsSet()

		appsList := getAllLocalApps()
		if len(appsList) == 0 {
			_, _ = cfmt.Println("\n{{ERROR:}}::red No apps found in pipeline '" + pipelineName + "'")
			os.Exit(1)
		}
		ensureAppNameIsSelected(appsList)
		upApp()
	},
}

var upPipelineCmd = &cobra.Command{
	Use:     "pipeline",
	Aliases: []string{"pl"},
	Short:   "Deploy a pipeline to the cluster",
	Long:    `Use the pipeline subcommand to deploy your pipelines to the cluster`,
	Run: func(cmd *cobra.Command, args []string) {

		pipelinesList := getAllLocalPipelines()
		ensurePipelineIsSet(pipelinesList)
		upPipeline()
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
	_, deployPipelineErr := api.DeployPipeline(pipeline)
	if deployPipelineErr != nil {
		_, _ = cfmt.Println("{{Error deploying pipeline}}::red", deployPipelineErr)
		return
	}
}

func upApp() {
	confirmationLine("Are you sure you want to deploy the app "+appName+" to "+pipelineName+"?", "y")
	app := loadLocalApp(pipelineName, stageName, appName)
	app.Spec.Pipeline = pipelineName             // ensure pipeline is set
	app.Spec.Phase = stageName                   // ensure stage is set
	app.Spec.Security.VulnerabilityScans = false // TODO: ask for this
	_, DeployAppErr := api.DeployApp(app)
	if DeployAppErr != nil {
		_, _ = cfmt.Println("{{Error deploying app}}::red", DeployAppErr)
		return
	}
}

func upAllPipelines() {

	confirmationLine("Are you sure you want to deploy all pipelines?", "y")
	pipelinesConfigs := loadAllLocalPipelines()

	_, _ = cfmt.Println("{{Deploying all pipelines}}::yellow")

	for _, pipelineCRD := range pipelinesConfigs {
		_, _ = cfmt.Println("{{Deploying pipeline}}::yellow " + pipelineCRD.Spec.Name + "")
		_, deployPipelineErr := api.DeployPipeline(pipelineCRD)
		if deployPipelineErr != nil {
			_, _ = cfmt.Println("{{Error deploying pipeline}}::red", deployPipelineErr)
			return
		}
	}
}
