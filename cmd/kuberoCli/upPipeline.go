/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package kuberoCli

import (
	"github.com/spf13/cobra"
)

// pipelineCmd represents the pipeline command
var upPipelineCmd = &cobra.Command{
	Use:     "pipeline",
	Aliases: []string{"pl"},
	Short:   "Deploy a pipeline to the cluster",
	Long:    `Use the pipeline subcommand to deploy your pipelines to the cluster`,
	Run: func(cmd *cobra.Command, args []string) {

		pipelinesList := getAllLocalPipelines()
		ensurePipelineIsSet(pipelinesList)
		upPipeline()
	},
}

func init() {
	upCmd.AddCommand(upPipelineCmd)
	upPipelineCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
}
