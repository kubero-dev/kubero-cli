package kuberoCli

import (
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Undeploy your pipelines and apps from the cluster",
	Long: `Use the pipeline or app subcommand to undeploy your pipelines and apps from the cluster
Subcommands:
  kubero down [pipeline|app]`,
	Run: func(cmd *cobra.Command, args []string) {
		cfmt.Println("{{Undeploying your pipeline}}::green")
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
	downCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
	downCmd.Flags().StringVarP(&appName, "app", "a", "", "name of the app")
}
