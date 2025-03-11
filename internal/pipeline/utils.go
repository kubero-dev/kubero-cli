package pipeline

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	a "github.com/faelmori/kubero-cli/internal/api"
)

func (m *PipelineManager) ensurePipelineIsSet(pipelinesList []string) error {
	if m.pipelineName == "" {
		fmt.Println("")
		prompt := &survey.Select{
			Message: "Select a pipeline",
			Options: pipelinesList,
		}
		askOneErr := survey.AskOne(prompt, &m.pipelineName)
		if askOneErr != nil {
			fmt.Println("Error while selecting pipeline:", askOneErr)
			return askOneErr
		}
	}

	return nil
}

func (m *PipelineManager) ensureAppNameIsSet() {
	if m.appName == "" {
		m.appName = promptLine("Define a app name", "", m.appName)
	}
}

func (m *PipelineManager) ensureStageNameIsSet() error {
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
			return askOneErr
		}
	}

	return nil
}

func (m *PipelineManager) ensureAppNameIsSelected(availableApps []string) {
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

func (m *PipelineManager) LoadRepositories() {
	if m.repositories == nil {
		repo := a.NewRepository("", "")
		repoReq, err := repo.GetRepositories()
		if err != nil {
			fmt.Println(err)
			return
		}
		m.repo = &repoReq
	}
}

func (m *PipelineManager) LoadContexts() {
	//TODO
}

func (m *PipelineManager) LoadBuildpacks() {
	//TODO
}
