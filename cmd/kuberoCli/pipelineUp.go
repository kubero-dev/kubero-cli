package kuberoCli

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/cobra"
)

// pipelineCmd represents the pipeline command
var upPipelineCmd = &cobra.Command{
	Use:     "iac:up",
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
	pipelineCmd.AddCommand(upPipelineCmd)
	upPipelineCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
}

func upPipeline() {
	confirmationLine("Are you sure you want to deploy the pipeline '"+pipelineName+"'?", "y")

	pipeline := loadPipelineConfig(pipelineName, true)
	_, deployPipelineErr := api.DeployPipeline(pipeline)
	if deployPipelineErr != nil {
		_, _ = cfmt.Println("{{Error deploying pipeline}}::red", deployPipelineErr)
		return
	}
}

/*
// does not make sense to deploy all pipelines at once
func upAllPipelines() {

	confirmationLine("Are you sure you want to deploy all pipelines?", "y")
	pipelinesConfigs := loadAllLocalPipelines()

	_, _ = cfmt.Println("{{Deploying all pipelines}}::yellow")

	for _, pipelineCRD := range pipelinesConfigs {
		_, _ = cfmt.Println("{{Deploying pipeline}}::yellow " + pipelineCRD.Spec.Name + "")
		_, deployPipelineErr := api.DeployPipeline(pipelineCRD)
		if deployPipelineErr != nil {
			_, _ = cfmt.Println("{{Error deploying pipeline}}::red", deployPipelineErr)
			return
		}
	}
}
*/
