package pipeline

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	a "github.com/faelmori/kubero-cli/internal/api"
)

func (m *PipelineManager) EnsurePipelineIsSet(pipelinesList []string) error {
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

func (m *PipelineManager) EnsureAppNameIsSet() error {
	if m.appName == "" {
		m.appName = promptLine("Define a app name", "", m.appName)
	}
	return nil
}

func (m *PipelineManager) EnsureStageNameIsSet() error {
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

func (m *PipelineManager) LoadRepositories() error {
	if m.repositories == nil {
		repo := a.NewRepository("", "")
		_, err := repo.GetRepositories()
		if err != nil {
			fmt.Println(err)
			return err
		}
		m.repo = &repo
	}
	return nil
}

func (m *PipelineManager) LoadContexts() {
	//TODO
}

func (m *PipelineManager) LoadBuildpacks() {
	//TODO
}
