package kuberoCli

import (
	"github.com/spf13/cobra"
)

// pipelinesCmd represents the pipelines command
var pipelinesCmd = &cobra.Command{
	Use:   "pipelines",
	Short: "**DEPRECATED** Manage your pipelines",
	Long: `List your pipelines
An App runs allways in a Pipeline. A Pipeline is a collection of Apps.`,
}

func init() {
	rootCmd.AddCommand(pipelinesCmd)
	pipelinesCmd.PersistentFlags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
}
