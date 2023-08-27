/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package kuberoCli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"cr", "add", "new"},
	Short:   "Create a new pipeline and/or app",
	Long:    `Initiate a new pipeline and app in your current repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		createPipelineAndApp()
		fmt.Println("create called")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.PersistentFlags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
}

func createPipelineAndApp() {
	createPipelineAndApp := promptLine("Create a new pipeline", "[y,n]", "y")
	if createPipelineAndApp == "y" {
		createPipeline()
	} else {
		fmt.Println("TODO : Load existing pipelines to select one") //TODO
	}

}
