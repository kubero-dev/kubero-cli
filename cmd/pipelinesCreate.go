package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var PipelineCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "**DEPRECATED** Create a new pipeline",
	Long:  `Create a new Pipeline`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create a new pipeline")

		loadRepositories()
		loadContexts()
		loadBuildpacks()
		createPipeline := pipelinesForm()

		client.SetBody(createPipeline.Spec)
		pipeline, pipelineErr := client.Post("/api/cli/pipelines/")

		if pipelineErr != nil {
			fmt.Println(pipelineErr)
		} else {
			cfmt.Println("{{Pipeline created successfully}}::green")
			json.Unmarshal(pipeline.Body(), &createPipeline.Spec)
			//writePipelineYaml(createPipeline)
		}

	},
}

func init() {
	pipelinesCmd.AddCommand(PipelineCreateCmd)
}
