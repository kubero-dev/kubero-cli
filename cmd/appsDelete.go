/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var appsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a existing app in a pipeline",
	Long:  `Delete a existing app in a pipeline`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete called")

		appInstance, _ := client.Delete("/api/cli/pipelines/" + pipeline + "/" + stage + "/" + app)

		fmt.Println(appInstance)

		fmt.Println(appInstance.StatusCode())
	},
}

var stage string
var app string

func init() {
	appsDeleteCmd.Flags().StringVarP(&pipeline, "pipeline", "p", "", "* Name of the pipeline")
	appsDeleteCmd.MarkFlagRequired("pipeline")

	appsDeleteCmd.Flags().StringVarP(&stage, "stage", "s", "", "* Name of the stage")
	appsDeleteCmd.MarkFlagRequired("stage")

	appsDeleteCmd.Flags().StringVarP(&app, "app", "a", "", "* Name of the app")
	appsDeleteCmd.MarkFlagRequired("app")

	appsCmd.AddCommand(appsDeleteCmd)
}
