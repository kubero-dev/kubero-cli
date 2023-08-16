package kuberoCli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Deploy your pipelines and apps to the cluster",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("up called")
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	upCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
	upCmd.Flags().StringVarP(&appName, "app", "a", "", "name of the app")
}
