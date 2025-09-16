package kuberoCli

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"encoding/json"
	"fmt"
	"kubero/pkg/kuberoApi"
	"os"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/cobra"
)

// pipelineCmd represents the pipeline command
var fetchPipelineCmd = &cobra.Command{
	Use:   "iac:fetch",
	Short: "Fetch a pipeline",
	Long:  `Fetch a pipeline`,
	Run: func(cmd *cobra.Command, args []string) {
		pipelinesList := getAllRemotePipelines()
		ensurePipelineIsSet(pipelinesList)
		fetchPipeline(pipelineName)
	},
}

func init() {
	pipelineCmd.AddCommand(fetchPipelineCmd)
	fetchPipelineCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "Name of the pipeline")
}

func fetchPipeline(pipelineName string) {
	confirmation := promptLine("Do you want to fetch the pipeline '"+pipelineName+"'?", "[y,n]", "y")
	if confirmation == "y" {
		_, _ = cfmt.Println("{{Fetching pipeline}}::yellow " + pipelineName)

		var pipeline kuberoApi.PipelineCRD

		pipeline.APIVersion = "application.kubero.dev/v1alpha1"
		pipeline.Kind = "KuberoPipeline"
		pipeline.Spec.Name = pipelineName
		pipeline.Metadata.Name = appName

		p, pipelineErr := api.GetPipeline(pipelineName)

		if pipelineErr != nil {
			if p == nil {
				_, _ = cfmt.Println("{{ERROR:}}::red Pipeline '" + pipelineName + "' not found ")
				os.Exit(1)
			}
			if p.StatusCode() == 404 {
				_, _ = cfmt.Println("{{ERROR:}}::red Pipeline '" + pipelineName + "' not found ")
				os.Exit(1)
			}
			fmt.Println(pipelineErr)
			os.Exit(1)
		}

		jsonUnmarshalErr := json.Unmarshal(p.Body(), &pipeline.Spec)
		if jsonUnmarshalErr != nil {
			fmt.Println(jsonUnmarshalErr)
			return
		}
		writePipelineYaml(pipeline)

	} else {
		return
	}
}

func fetchAllPipelines() {

	confirmation := promptLine("Fetch App or Pipeline?", "[app,pipeline]", "app")
	if confirmation == "app" {
		_, _ = cfmt.Println("{{Fetching app}}::yellow")
		fetchAppCmd.Run(fetchAppCmd, []string{})
	} else {
		_, _ = cfmt.Println("{{Fetching pipelines}}::yellow")
		pipelinesList := getAllRemotePipelines()
		if len(pipelinesList) == 0 {
			_, _ = cfmt.Println("{{ERROR:}}::red No pipelines found")
			os.Exit(1)
		}
		for _, pipeline := range pipelinesList {
			fetchPipeline(pipeline)
		}
	}
}
