/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var pipelinesDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a existing pipeline",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete called")

		pipeline, _ := client.Delete("/api/cli/pipelines/" + pipeline)
		fmt.Println(pipeline)
	},
}

func init() {
	pipelinesCmd.AddCommand(pipelinesDeleteCmd)
	pipelinesDeleteCmd.Flags().StringVarP(&pipeline, "pipeline", "p", "", "Name of the pipeline")

	pipelinesDeleteCmd.MarkFlagRequired("pipeline")
}
