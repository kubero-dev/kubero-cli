package cli

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"os"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/cobra"
)

// appCmd represents the app command
var fetchAppCmd = &cobra.Command{
	Use:   "app",
	Short: "Fetch an app",
	Long:  `Fetch an app`,
	Run: func(cmd *cobra.Command, args []string) {

		pipelinesList := getAllRemotePipelines()
		if len(pipelinesList) == 0 {
			_, _ = cfmt.Println("\n{{ERROR:}}::red No pipelines found")
			os.Exit(1)
		}
		ensurePipelineIsSet(pipelinesList)
		ensureStageNameIsSet()
		fetchPipeline(pipelineName)

		appsList := getAllRemoteApps()
		if len(appsList) == 0 {
			_, _ = cfmt.Println("\n{{ERROR:}}::red No apps found in pipeline '" + pipelineName + "'")
			os.Exit(1)
		}

		ensureAppNameIsSelected(appsList)
		ensureAppNameIsSet()
		fetchApp(appName, stageName, pipelineName)
	},
}

func init() {
	fetchCmd.AddCommand(fetchAppCmd)
	fetchAppCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "Name of the pipeline")
	fetchAppCmd.Flags().StringVarP(&stageName, "stage", "s", "", "Name of the stage [test|stage|production]")
	fetchAppCmd.Flags().StringVarP(&appName, "app", "a", "", "Name of the app")
}
