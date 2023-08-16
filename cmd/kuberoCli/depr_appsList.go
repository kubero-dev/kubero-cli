package kuberoCli

import (
	"os"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var appsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List apps in a pipeline",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		if pipelineName == "" {
			pipelineName = pipelineConfig.GetString("spec.name")
			if pipelineName == "" {
				cfmt.Println("{{  Pipeline not found in config file}}::red")
				os.Exit(1)
			}
		}

		appsList()
	},
}

func init() {
	appsListCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "Name of the Pipeline")
	//appsListCmd.MarkFlagRequired("pipeline")
	appsCmd.AddCommand(appsListCmd)
}
