package kuberoCli

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"github.com/spf13/cobra"
)

// appCmd represents the app command
var logsAppCmd = &cobra.Command{
	Use:     "logs",
	Aliases: []string{"d"},
	Short:   "Load logs for an app",
	Long:    `Load logs for an app in a specific stage of a pipeline`,
	Example: `kubero app logs -p my-pipeline -s stage -a my-app`,
	Run: func(cmd *cobra.Command, args []string) {
		getAppLogHistory()
	},
}

func init() {
	AppCmd.AddCommand(logsAppCmd)
	logsAppCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
	logsAppCmd.Flags().StringVarP(&stageName, "stage", "s", "", "Name of the stage [test|stage|production]")
	logsAppCmd.Flags().StringVarP(&appName, "app", "a", "", "name of the app")
}

func getAppLogHistory() {

	pipelinesList := getAllRemotePipelines()
	ensurePipelineIsSet(pipelinesList)
	ensureStageNameIsSet()

	appsList := getAllRemoteApps()
	ensureAppNameIsSelected(appsList)

	_, err := api.GetLogs(pipelineName, stageName, appName, "kubero-web")
	if err != nil {
		panic("Unable to fetch App logs")
	}
}
