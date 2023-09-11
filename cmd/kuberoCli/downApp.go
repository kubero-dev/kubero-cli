/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package kuberoCli

import (
	"github.com/spf13/cobra"
)

// appCmd represents the app command
var downAppCmd = &cobra.Command{
	Use:   "app",
	Short: "Undeploy an apps from the cluster",
	Long:  `Use the app subcommand to undeploy your apps from the cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		downApp()
	},
}

func init() {
	downCmd.AddCommand(downAppCmd)
	downAppCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
	downAppCmd.Flags().StringVarP(&appName, "app", "a", "", "name of the app")
}
