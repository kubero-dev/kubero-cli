package cli

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pipelineName string
var stageName string
var appName string

// pipelineCmd represents the pipeline command
var createPipelineCmd = &cobra.Command{
	Use:     "pipeline",
	Aliases: []string{"pl"},
	Short:   "Create a new pipeline",
	Long:    `Create a new Pipeline`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create a new pipeline")

		_ = createPipeline()
	},
}

func init() {
	createCmd.AddCommand(createPipelineCmd)
}
