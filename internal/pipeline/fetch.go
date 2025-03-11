package pipeline

import (
	"encoding/json"
	"fmt"
	a "github.com/faelmori/kubero-cli/internal/api"
	"github.com/faelmori/kubero-cli/internal/log"
	"github.com/faelmori/kubero-cli/types"
	"github.com/i582/cfmt/cmd/cfmt"
	"gopkg.in/yaml.v3"
	"os"
)

func (m *PipelineManager) FetchPipeline(pipelineName string) error {
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
				return pipelineErr
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
		if writeYamlErr := m.WriteYamlPipeline(&pipeline); writeYamlErr != nil {
			log.Error("Unable to write pipeline to file")
			return writeYamlErr
		}
	}
	return nil
}

func (m *PipelineManager) FetchApp(appName string, stageName string, pipelineName string) error {

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
			return appErr
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

	if writeYamlErr := m.WriteYamlApp(&app); writeYamlErr != nil {
		log.Error("Unable to write app to file")
		return writeYamlErr
	}

	return nil
}

func (m *PipelineManager) FetchAllPipelines() {
	confirmation := promptLine("Are you sure you want to fetch all pipelines?", "[y,n]", "n")
	if confirmation == "y" {
		_, _ = cfmt.Println("{{Fetching all pipelines}}::yellow")
	} else {
		_, _ = cfmt.Println("{{Aborted}}::red")
		return
	}
}

func (m *PipelineManager) WriteYamlPipeline(pipeline *types.PipelineCRD) error {
	if pipeline == nil {
		return os.ErrInvalid
	}
	yamlData, yamlErr := yaml.Marshal(pipeline)
	if yamlErr != nil {
		return yamlErr
	}
	file, fileErr := os.Create(pipeline.Metadata.Name + ".yaml")
	if fileErr != nil {
		return fileErr
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	_, writeErr := file.Write(yamlData)
	if writeErr != nil {
		return writeErr
	}
	return nil
}

func (m *PipelineManager) WriteYamlApp(app *types.AppCRD) error {
	if app == nil {
		return os.ErrInvalid
	}
	yamlData, yamlErr := yaml.Marshal(app)
	if yamlErr != nil {
		return yamlErr
	}
	file, fileErr := os.Create(app.Metadata.Name + ".yaml")
	if fileErr != nil {
		return fileErr
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	_, writeErr := file.Write(yamlData)
	if writeErr != nil {
		return writeErr
	}
	return nil
}
