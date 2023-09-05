package kuberoCli

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var pipelinesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the Pipelines",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		if pipelineName != "" {
			// get a single pipeline
			pipelineResp, _ := client.Get("/api/cli/pipelines/" + pipelineName)
			printPipeline(pipelineResp)
		} else {
			// get the pipelines
			pipelineListResp, _ := client.Get("/api/cli/pipelines")
			printPipelinesList(pipelineListResp)
		}
	},
}

func init() {
	pipelinesCmd.AddCommand(pipelinesListCmd)
}
