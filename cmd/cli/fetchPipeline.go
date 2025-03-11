package cli

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"github.com/spf13/cobra"
)

// pipelineCmd represents the pipeline command
var fetchPipelineCmd = &cobra.Command{
	Use:     "pipeline",
	Aliases: []string{"pl"},
	Short:   "Fetch a pipeline",
	Long:    `Fetch a pipeline`,
	Run: func(cmd *cobra.Command, args []string) {
		pipelinesList := getAllRemotePipelines()
		ensurePipelineIsSet(pipelinesList)
		fetchPipeline(pipelineName)
	},
}

func init() {
	fetchCmd.AddCommand(fetchPipelineCmd)
	fetchPipelineCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "Name of the pipeline")
}
