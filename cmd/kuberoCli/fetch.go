package kuberoCli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Sync youre remote pipelines and apps to your local repository",
	Long:  `Use the pipeline or app subcommand to sync your pipelines and apps to your local repository`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fetch called")
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
	fetchCmd.Flags().StringVarP(&appName, "app", "a", "", "name of the app")
}
