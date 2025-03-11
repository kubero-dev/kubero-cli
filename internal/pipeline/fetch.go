package pipeline

import (
	"fmt"
	"github.com/faelmori/kubero-cli/internal/api"
	"github.com/faelmori/kubero-cli/internal/utils"
	"github.com/faelmori/kubero-cli/pkg/kuberoApi"
	"github.com/i582/cfmt/cmd/cfmt"
	"os"
)

var promptLine = utils.NewConsolePrompt().PromptLine

func FetchPipeline(pipelineName string) {
	confirmation := promptLine("Do you want to fetch the pipeline '"+pipelineName+"'?", "[y,n]", "y")
	if confirmation == "y" {
		_, _ = cfmt.Println("{{Fetching pipeline}}::yellow " + pipelineName)

		var pipeline kuberoApi.PipelineCRD

		pipeline.APIVersion = "application.kubero.dev/v1alpha1"
		pipeline.Kind = "KuberoPipeline"
		pipeline.Spec.Name = pipelineName
		pipeline.Metadata.Name = appName
		mdLbList := pipeline.Metadata.Labels.(map[string]interface{})
		token := ""
		if t, ok := mdLbList["token"]; ok {
			token = t.(string)
		}
		apiClient := api.NewClient(pipeline.Spec.Domain, token)
		p, pipelineErr := apiClient.GetPipeline(pipelineName)

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

func FetchApp(appName string, stageName string, pipelineName string) {

	confirmation := promptLine("Do you want to fetch the app '"+appName+"' from '"+pipelineName+"'?", "[y,n]", "y")
	if confirmation == "y" {
		_, _ = cfmt.Println("{{Fetching app}}::yellow " + appName + "")
	} else {
		_, _ = cfmt.Println("{{Aborted}}::red")
		return
	}

	var app kuberoApi.AppCRD
	app.APIVersion = "application.kubero.dev/v1alpha1"
	app.Kind = "KuberoApp"

	app.Spec.Pipeline = pipelineName
	app.Spec.Phase = stageName
	app.Spec.Name = appName
	app.Metadata.Name = appName

	a, appErr := api.GetApp(pipelineName, stageName, appName)

	if appErr != nil {
		if a == nil {
			_, _ = cfmt.Println("{{ERROR:}}::red App '" + appName + "' not found ")
			os.Exit(1)
		}
		if a.StatusCode() == 404 {
			_, _ = cfmt.Println("{{ERROR:}}::red App '" + appName + "' not found ")
			os.Exit(1)
		}
		fmt.Println(appErr)
		os.Exit(1)
	}

	jsonUnmarshalErr := json.Unmarshal(a.Body(), &app)
	if jsonUnmarshalErr != nil {
		fmt.Println(jsonUnmarshalErr)
		return
	}

	writeAppYaml(app)

}

func FetchAllPipelines() {
	confirmation := promptLine("Are you sure you want to fetch all pipelines?", "[y,n]", "n")
	if confirmation == "y" {
		_, _ = cfmt.Println("{{Fetching all pipelines}}::yellow")
	} else {
		_, _ = cfmt.Println("{{Aborted}}::red")
		return
	}
}
