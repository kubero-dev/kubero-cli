/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package kuberoCli

import (
	"github.com/spf13/cobra"
)

// appCmd represents the app command
var fetchAppCmd = &cobra.Command{
	Use:   "app",
	Short: "Fetch an app",
	Long:  `Fetch an app`,
	Run: func(cmd *cobra.Command, args []string) {
		ensureStageNameIsSet()
		fetchPipeline(pipelineName)
		fetchApp(appName, stageName, pipelineName)
	},
}

func init() {
	fetchCmd.AddCommand(fetchAppCmd)
	fetchAppCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "Name of the pipeline")
	fetchAppCmd.Flags().StringVarP(&stageName, "stage", "s", "", "Name of the stage [test|stage|production]")
	fetchAppCmd.Flags().StringVarP(&appName, "app", "a", "", "Name of the app")
}
