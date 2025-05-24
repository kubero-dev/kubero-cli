package kuberoCli

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"github.com/spf13/cobra"
)

// appCmd represents the app command
var deletePipelineCmd = &cobra.Command{
	Use:     "down",
	Aliases: []string{"d"},
	Short:   "Undeploy a pipeline from the cluster",
	Long:    `Use the pipeline subcommand to undeploy your pipelines from the cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		deletePipeline()
	},
}

func init() {
	pipelineCmd.AddCommand(deletePipelineCmd)
	deletePipelineCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
}

func deletePipeline() {
	pipelinesList := getAllRemotePipelines()
	ensurePipelineIsSet(pipelinesList)
	deletePipelineByName(pipelineName)
}

func deletePipelineByName(pipelineName string) {
	confirmationLine("Are you sure you want to undeploy the pipeline '"+pipelineName+"'?", "y")

	_, err := api.DeletePipeline(pipelineName)
	if err != nil {
		panic("Unable to undeploy Pipeline")
	}
}
