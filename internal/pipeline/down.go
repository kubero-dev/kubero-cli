package pipeline

func DownPipeline() {
	pipelinesList := getAllRemotePipelines()
	ensurePipelineIsSet(pipelinesList)
	downPipelineByName(pipelineName)
}

func downPipelineByName(pipelineName string) {
	confirmationLine("Are you sure you want to undeploy the pipeline '"+pipelineName+"'?", "y")

	_, err := api.UnDeployPipeline(pipelineName)
	if err != nil {
		panic("Unable to undeploy Pipeline")
	}
}

func DownApp() {

	pipelinesList := getAllRemotePipelines()
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

func DownAllPipelines() {
	confirmationLine("Are you sure you want to undeploy all pipelines?", "y")
	pipelinesList := getAllLocalPipelines()
	for _, pipeline := range pipelinesList {
		downPipelineByName(pipeline)
	}
}
