package pipeline

import (
	"encoding/json"
	"fmt"
	a "github.com/faelmori/kubero-cli/internal/api"
	"github.com/faelmori/kubero-cli/internal/log"
	"github.com/faelmori/kubero-cli/types"
	"github.com/i582/cfmt/cmd/cfmt"
)

func (m *ManagerPipeline) FetchPipeline(pipelineName string) error {
	confirmation := promptLine("Do you want to fetch the pipeline '"+pipelineName+"'?", "[y,n]", "y")
	if confirmation == "y" {
		log.Info("Fetching pipeline " + pipelineName)

		var pipeline types.PipelineCRD

		pipeline.APIVersion = "application.kubero.dev/v1alpha1"
		pipeline.Kind = "KuberoPipeline"
		pipeline.Spec.Name = pipelineName
		pipeline.Metadata.Name = m.appName
		mdLbList := pipeline.Metadata.Labels.(map[string]interface{})
		token := ""
		if t, ok := mdLbList["token"]; ok {
			token = t.(string)
		}

		api := a.NewClient()
		api.Init(pipeline.Spec.Domain, token)
		p, pipelineErr := api.GetPipeline(pipelineName)

		if pipelineErr != nil {
			if p == nil {
				log.Error("Pipeline '" + pipelineName + "' not found ")
			}
			if p.StatusCode() == 404 {
				log.Error("Pipeline '" + pipelineName + "' not found ")
			}
			fmt.Println(pipelineErr)
			return pipelineErr
		}

		jsonUnmarshalErr := json.Unmarshal(p.Body(), &pipeline.Spec)
		if jsonUnmarshalErr != nil {
			log.Error("Unable to decode response")
			return jsonUnmarshalErr
		}
		m.WritePipelineYaml(pipeline)
	}
	return nil
}

func (m *ManagerPipeline) FetchApp(appName string, stageName string, pipelineName string) error {

	confirmation := promptLine("Do you want to fetch the app '"+appName+"' from '"+pipelineName+"'?", "[y,n]", "y")
	if confirmation == "y" {
		log.Info("Fetching app " + appName)
	} else {
		log.Warn("Aborted")
		return nil
	}

	var app types.AppCRD
	app.APIVersion = "application.kubero.dev/v1alpha1"
	app.Kind = "KuberoApp"

	app.Spec.Pipeline = pipelineName
	app.Spec.Phase = stageName
	app.Spec.Name = appName
	app.Metadata.Name = appName

	api := a.NewClient()
	api.Init(app.Spec.Domain, "")

	apiClientApp, appErr := api.GetApp(pipelineName, stageName, appName)

	if appErr != nil {
		if apiClientApp == nil {
			log.Error("App '" + appName + "' not found ")
		}
		if apiClientApp.StatusCode() == 404 {
			log.Error("App '" + appName + "' not found ")
		}
		return appErr
	}

	jsonUnmarshalErr := json.Unmarshal(apiClientApp.Body(), &app)
	if jsonUnmarshalErr != nil {
		log.Error("Unable to decode response")
		return jsonUnmarshalErr
	}

	m.WriteAppYaml(app)

}

func (m *ManagerPipeline) FetchAllPipelines() {
	confirmation := promptLine("Are you sure you want to fetch all pipelines?", "[y,n]", "n")
	if confirmation == "y" {
		_, _ = cfmt.Println("{{Fetching all pipelines}}::yellow")
	} else {
		_, _ = cfmt.Println("{{Aborted}}::red")
		return
	}
}
