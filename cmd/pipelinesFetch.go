package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var pipelinesFetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch a existing pipeline configuration",
	Run: func(cmd *cobra.Command, args []string) {
		createPipeline := pipelinesFetchForm()

		client.SetBody(createPipeline.Spec)
		p, pipelineErr := client.Get("/api/cli/pipelines/" + createPipeline.Spec.Name)

		if pipelineErr != nil {
			fmt.Println(pipelineErr)
		} else {
			json.Unmarshal(p.Body(), &createPipeline.Spec)
			//json.Unmarshal(p.Body(), &createPipeline)
			//writePipelineYaml(createPipeline)
		}
	},
}

func init() {
	pipelinesFetchCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Skip asking for confirmation")
	pipelinesFetchCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "Name of the pipeline")
	pipelinesCmd.AddCommand(pipelinesFetchCmd)
}

func pipelinesFetchForm() PipelineCRD {

	var cp PipelineCRD

	cp.APIVersion = "application.kubero.dev/v1alpha1"
	cp.Kind = "KuberoPipeline"

	if pipelineName == "" {
		pipelineName = pipelineConfig.GetString("spec.name")
	}
	cp.Spec.Name = promptLine("Pipeline", "", pipelineName)

	return cp
}
