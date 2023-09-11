/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package kuberoCli

import (
	"github.com/spf13/cobra"
)

// appCmd represents the app command
var upAppCmd = &cobra.Command{
	Use:   "app",
	Short: "Deploy an apps to the cluster",
	Long:  `Use the app subcommand to deploy your apps to the cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		upApp()
	},
}

func init() {
	upCmd.AddCommand(upAppCmd)
	upAppCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
	upAppCmd.Flags().StringVarP(&appName, "app", "a", "", "name of the app")
}
