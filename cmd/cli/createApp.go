package cli

import (
	"github.com/spf13/cobra"
)

// appCmd represents the app command
var createAppCmd = &cobra.Command{
	Use:   "app",
	Short: "Create a new app in a Pipeline",
	Long: `Create a new app in a Pipeline.

If called without arguments, it will ask for all the required information`,
	Run: func(cmd *cobra.Command, args []string) {

		pipelinesList := getAllLocalPipelines()
		ensurePipelineIsSet(pipelinesList)
		ensureStageNameIsSet()
		ensureAppNameIsSet()
		createApp()

	},
}

func init() {
	createCmd.AddCommand(createAppCmd)
	createAppCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "Name of the pipeline")
	createAppCmd.Flags().StringVarP(&stageName, "stage", "s", "", "Name of the stage")
	createAppCmd.Flags().StringVarP(&appName, "app", "a", "", "Name of the app")
}
