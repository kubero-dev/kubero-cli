package pipeline

import (
	"fmt"
	"github.com/faelmori/kubero-cli/internal/api"
	"github.com/i582/cfmt/cmd/cfmt"
	"os"
)

func (m *ManagerPipeline) ListPipelines(pipelineName, outputFormat string) error {
	apiClient := api.NewClient("/api/cli/pipelines/"+pipelineName, "")
	client := apiClient.RestyClient.GetClient()

	if pipelineName != "" {
		pipelineResp, err := client.Get("/api/cli/pipelines/" + pipelineName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if pipelineResp.StatusCode == 404 {
			_, _ = cfmt.Println("{{  Pipeline not found}}::red")
			os.Exit(1)
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		printPipeline(pipelineResp.Request)

		appsList()
	} else {
		// get the pipelines
		pipelineListResp, err := client.Get("/api/cli/pipelines")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		printPipelinesList(pipelineListResp)

	}
}
