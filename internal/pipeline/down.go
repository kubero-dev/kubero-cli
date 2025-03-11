package pipeline

func (m *ManagerPipeline) DownPipeline() {
	pipelinesList := getAllRemotePipelines()
	ensurePipelineIsSet(pipelinesList)
	downPipelineByName(pipelineName)
}

func (m *ManagerPipeline) downPipelineByName(pipelineName string) {
	confirmationLine("Are you sure you want to undeploy the pipeline '"+pipelineName+"'?", "y")

	_, err := api.UnDeployPipeline(pipelineName)
	if err != nil {
		panic("Unable to undeploy Pipeline")
	}
}

func (m *ManagerPipeline) DownApp() {

	pipelinesList := m.getAllRemotePipelines()
	ensurePipelineIsSet(pipelinesList)

	ensureStageNameIsSet()

	appsList := getAllRemoteApps()
	ensureAppNameIsSelected(appsList)

	confirmationLine("Are you sure you want to undeploy the app "+appName+" from "+stageName+" in "+pipelineName+"?", "y")

	_, err := api.UnDeployApp(pipelineName, stageName, appName)
	if err != nil {
		panic("Unable to undeploy App")
	}
}

func (m *ManagerPipeline) DownAllPipelines() {
	confirmationLine("Are you sure you want to undeploy all pipelines?", "y")
	pipelinesList := getAllLocalPipelines()
	for _, pipeline := range pipelinesList {
		downPipelineByName(pipeline)
	}
}
