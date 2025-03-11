package cli

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"github.com/spf13/cobra"
)

// appCmd represents the app command
var downPipelineCmd = &cobra.Command{
	Use:     "pipeline",
	Aliases: []string{"pl"},
	Short:   "Undeploy a pipeline from the cluster",
	Long:    `Use the pipeline subcommand to undeploy your pipelines from the cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		downPipeline()
	},
}

func init() {
	downCmd.AddCommand(downPipelineCmd)
	downPipelineCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
}
