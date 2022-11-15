/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
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
	deleteCmd.Flags().StringVarP(&pipeline, "pipeline", "p", "", "* Name of the pipeline")
	deleteCmd.MarkFlagRequired("pipeline")

	deleteCmd.Flags().StringVarP(&stage, "stage", "s", "", "* Name of the stage")
	deleteCmd.MarkFlagRequired("stage")

	deleteCmd.Flags().StringVarP(&app, "app", "a", "", "* Name of the app")
	deleteCmd.MarkFlagRequired("app")

	appsCmd.AddCommand(deleteCmd)
}
