package pipeline

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
)

func (m *ManagerPipeline) ensurePipelineIsSet(pipelinesList []string) {
	if m.pipelineName == "" {
		fmt.Println("")
		prompt := &survey.Select{
			Message: "Select a pipeline",
			Options: pipelinesList,
		}
		askOneErr := survey.AskOne(prompt, &m.pipelineName)
		if askOneErr != nil {
			fmt.Println("Error while selecting pipeline:", askOneErr)
			return
		}
	}
}

func (m *ManagerPipeline) ensureAppNameIsSet() {
	if m.appName == "" {
		m.appName = promptLine("Define a app name", "", m.appName)
	}
}

func (m *ManagerPipeline) ensureStageNameIsSet() {
	if m.stageName == "" {
		fmt.Println("")
		pipelineConfig := m.loadPipelineConfig(m.pipelineName)
		availablePhases := m.getPipelinePhases(pipelineConfig)
		prompt := &survey.Select{
			Message: "Select a stage",
			Options: availablePhases,
		}
		askOneErr := survey.AskOne(prompt, &m.stageName)
		if askOneErr != nil {
			fmt.Println("Error while selecting stage:", askOneErr)
			return
		}
	}
}

func (m *ManagerPipeline) ensureAppNameIsSelected(availableApps []string) {
	if m.appName == "" {
		fmt.Println("")
		prompt := &survey.Select{
			Message: "Select an app",
			Options: availableApps,
		}
		askOneErr := survey.AskOne(prompt, &m.appName)
		if askOneErr != nil {
			fmt.Println("Error while selecting app:", askOneErr)
			return
		}
	}
}
