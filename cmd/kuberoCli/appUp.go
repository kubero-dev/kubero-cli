package kuberoCli

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"os"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/cobra"
)

// appCmd represents the app command
var upAppCmd = &cobra.Command{
	Use:   "iac:up",
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

func init() {
	AppCmd.AddCommand(upAppCmd)
	upAppCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
	upAppCmd.Flags().StringVarP(&stageName, "stage", "s", "", "Name of the stage [test|stage|production]")
	upAppCmd.Flags().StringVarP(&appName, "app", "a", "", "name of the app")
}

func upApp() {
	confirmationLine("Are you sure you want to deploy the app "+appName+" to "+pipelineName+"?", "y")
	app := loadLocalApp(pipelineName, stageName, appName)
	app.Spec.Pipeline = pipelineName             // ensure pipeline is set
	app.Spec.Phase = stageName                   // ensure stage is set
	app.Spec.Security.VulnerabilityScans = false // TODO: ask for this
	_, DeployAppErr := api.DeployApp(pipelineName, stageName, appName, app)
	if DeployAppErr != nil {
		_, _ = cfmt.Println("{{Error deploying app}}::red", DeployAppErr)
		return
	}
}
