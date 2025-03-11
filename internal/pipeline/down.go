package pipeline

import a "github.com/faelmori/kubero-cli/internal/api"

func (m *PipelineManager) DownPipeline() {
	pipelinesList := m.GetAllRemotePipelines()
	m.EnsurePipelineIsSet(pipelinesList)
	m.downPipelineByName(m.pipelineName)
}

func (m *PipelineManager) downPipelineByName(pipelineName string) {
	confirmationLine("Are you sure you want to undeploy the pipeline '"+pipelineName+"'?", "y")
	client := a.NewClient()
	_, err := client.UnDeployPipeline(pipelineName)
	if err != nil {
		panic("Unable to undeploy Pipeline")
	}
}

func (m *PipelineManager) DownApp() {

	pipelinesList := m.GetAllRemotePipelines()
	m.EnsurePipelineIsSet(pipelinesList)

	m.EnsureStageNameIsSet()

	appsList := m.GetAllRemoteApps()
	m.ensureAppNameIsSelected(appsList)

	confirmationLine("Are you sure you want to undeploy the app "+m.appName+" from "+m.stageName+" in "+m.pipelineName+"?", "y")
	client := a.NewClient()
	_, err := client.UnDeployApp(m.pipelineName, m.stageName, m.appName)
	if err != nil {
		panic("Unable to undeploy App")
	}
}

func (m *PipelineManager) DownAllPipelines() {
	confirmationLine("Are you sure you want to undeploy all pipelines?", "y")
	pipelinesList := m.GetAllLocalPipelines()
	for _, pipeline := range pipelinesList {
		m.downPipelineByName(pipeline)
	}
}
