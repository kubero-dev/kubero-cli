package kuberoCli

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"github.com/spf13/cobra"
)

// appCmd represents the app command
var deleteAppCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"d"},
	Short:   "Deletes an app from the cluster",
	Long:    `Use the app subcommand to undeploy your apps from the cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		downApp()
	},
}

func init() {
	AppCmd.AddCommand(deleteAppCmd)
	deleteAppCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
	deleteAppCmd.Flags().StringVarP(&stageName, "stage", "s", "", "Name of the stage [test|stage|production]")
	deleteAppCmd.Flags().StringVarP(&appName, "app", "a", "", "name of the app")
}

func downApp() {

	pipelinesList := getAllRemotePipelines()
	ensurePipelineIsSet(pipelinesList)

	ensureStageNameIsSet()

	appsList := getAllRemoteApps()
	ensureAppNameIsSelected(appsList)

	confirmationLine("Are you sure you want to undeploy the app "+appName+" from "+stageName+" in "+pipelineName+"?", "y")

	_, err := api.DeleteApp(pipelineName, stageName, appName)
	if err != nil {
		panic("Unable to undeploy App")
	}
}
