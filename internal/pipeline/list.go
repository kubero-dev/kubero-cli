package pipeline

import (
	"fmt"
	"github.com/i582/cfmt/cmd/cfmt"
	a "github.com/kubero-dev/kubero-cli/internal/api"
	"os"
)

func (m *PipelineManager) ListPipelines(pipelineName, outputFormat string) error {
	ac := a.NewClient()
	client := ac.Init("/api/cli/pipelines", m.GetCredentialsManager().GetCredentials().GetString("token"))

	if pipelineName != "" {
		pipelineResp, err := client.Get("/api/cli/pipelines/" + pipelineName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if pipelineResp.StatusCode() == 404 {
			_, _ = cfmt.Println("{{  Pipeline not found}}::red")
			os.Exit(1)
		}

		m.printPipeline(pipelineResp)

		if appsListErr := m.AppsList(pipelineName, outputFormat); appsListErr != nil {
			fmt.Println(appsListErr)
			os.Exit(1)
		}
	} else {
		pipelineListResp, err := client.Get("/api/cli/pipelines")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if printPipelineList := m.printPipelinesList(pipelineListResp); printPipelineList != nil {
			fmt.Println(printPipelineList)
			os.Exit(1)
		}
	}

	return nil
}
