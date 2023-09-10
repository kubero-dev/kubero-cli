/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package kuberoCli

import (
	"fmt"
	"kubero/pkg/kuberoApi"
	"os"
	"strconv"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"cr", "add", "new"},
	Short:   "Create a new pipeline and/or app",
	Long:    `Initiate a new pipeline and app in your current repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		createPipelineAndApp()
		fmt.Println("create called")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.PersistentFlags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
}

func createPipelineAndApp() {
	createPipelineAndApp := promptLine("Create a new pipeline", "[y,n]", "y")
	if createPipelineAndApp == "y" {
		createPipeline()
	} else {
		fmt.Println("TODO : Load existing pipelines to select one") //TODO
	}

}

func appForm(AppName string, pipelineName string) AppCRD {

	var app AppCRD

	app.APIVersion = "application.kubero.dev/v1alpha1"
	app.Kind = "KuberoApp"

	if appName == "" {
		app.Spec.Name = promptLine("App Name", "", appName)
	} else {
		app.Spec.Name = appName
	}

	if pipelineName == "" {
		app.Spec.Pipeline = promptLine("Pipeline Name", "", pipelineName)
	} else {
		app.Spec.Pipeline = pipelineName
	}

	pipelineConfig := getPipelineConfig(pipelineName)
	availablePhases := getPipelinePhases(pipelineConfig)
	if stageName == "" {
		app.Spec.Phase = promptLine("Phase", fmt.Sprint(availablePhases), stageName)
	} else {
		app.Spec.Phase = stageName
	}

	app.Spec.Domain = promptLine("Domain", "", "")

	gitURL := pipelineConfig.GetString("spec.git.repository.sshurl")
	//ca.Spec.Gitrepo.SSHURL = promptLine("Git SSH URL", "["+getGitRemote()+"]", gitURL)

	//ca.Spec.Gitrepo.SSHURL = pipelineConfig.GetString("spec.git.repository")
	pipelineConfig.UnmarshalKey("spec.git.repository", &app.Spec.Gitrepo)
	app.Spec.Branch = promptLine("Branch", gitURL+":", "")

	app.Spec.Buildpack = pipelineConfig.GetString("spec.buildpack.name")

	autodeploy := promptLine("Autodeploy", "[y,n]", "n")
	if autodeploy == "Y" {
		app.Spec.Autodeploy = true
	} else {
		app.Spec.Autodeploy = false
	}

	envCount, _ := strconv.Atoi(promptLine("Number of Env Vars", "", "0"))
	for i := 0; i < envCount; i++ {
		app.Spec.EnvVars = append(app.Spec.EnvVars, promptLine("Env Var", "", ""))
	}

	app.Spec.Image.ContainerPort, _ = strconv.Atoi(promptLine("Container Port", "8080", ""))

	app.Spec.Web.ReplicaCount, _ = strconv.Atoi(promptLine("Web Pods", "1", ""))

	app.Spec.Worker.ReplicaCount, _ = strconv.Atoi(promptLine("Worker Pods", "0", ""))

	return app
}

func createPipeline() kuberoApi.PipelineCRD {

	if pipelineName == "" {
		pipelineName = promptLine("Pipeline Name", "", "")
	}

	loadConfigs("/.kubero/", pipelineName)

	loadRepositories()
	loadContexts()
	loadBuildpacks()
	pipelineYaml := pipelinesForm()

	writePipelineYaml(pipelineYaml)

	cfmt.Println("\n\n{{Created pipeline.yaml}}::green")
	cfmt.Println(pipelineName)

	return pipelineYaml
}

func writePipelineYaml(pipeline kuberoApi.PipelineCRD) {
	basePath := "/.kubero/" //TODO Make it dynamic

	gitdir := getGitdir()
	dir := gitdir + basePath + pipelineName
	err := os.MkdirAll(dir, 0755)

	if err != nil {
		fmt.Println(err)
		panic("Unable to create directory")
	}

	yamlData, err := yaml.Marshal(&pipeline)

	// iterate over phases to create the directory
	for _, phase := range pipeline.Spec.Phases {
		if phase.Enabled {
			err := os.MkdirAll(dir+"/"+phase.Name, 0755)
			if err != nil {
				fmt.Println(err)
				panic("Unable to create directory")
			}
		}
	}

	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}
	//fmt.Println(string(yamlData))

	fileName := dir + "/pipeline.yaml"
	err = os.WriteFile(fileName, yamlData, 0644)
	if err != nil {
		panic("Unable to write data into the file")
	}

}

func pipelinesForm() kuberoApi.PipelineCRD {

	var cp kuberoApi.PipelineCRD

	cp.APIVersion = "application.kubero.dev/v1alpha1"
	cp.Kind = "KuberoPipeline"

	selectedBuildpack := pipelineConfig.GetString("spec.buildpack.name")
	cp.Spec.Buildpack.Name = promptLine("Buildpack ", fmt.Sprint(buildPacksSimpleList), selectedBuildpack)

	domain := pipelineConfig.GetString("spec.domain")
	cp.Spec.Domain = promptLine("FQDN Domain ", "", domain)

	// those fields are deprecated and may be removed in the future
	cp.Spec.Dockerimage = ""
	cp.Spec.Deploymentstrategy = "git"

	gitconnection := promptLine("Connect pipeline to a Git repository (GitOps)", "[y,n]", "n")

	if gitconnection == "y" {
		gitPrivider := pipelineConfig.GetString("spec.git.repository.provider")
		cp.Spec.Git.Repository.Provider = promptLine("Repository Provider", fmt.Sprint(repoSimpleList), gitPrivider)

		gitURL := pipelineConfig.GetString("spec.git.repository.sshurl")
		cp.Spec.Git.Repository.SSHURL = promptLine("Repository URL", "["+getGitRemote()+"]", gitURL)
	}

	phaseReview := promptLine("enable reviewapps", "[y,n]", "n")
	if phaseReview == "y" {
		cp.Spec.Reviewapps = true
		contextDefault := pipelineConfig.GetString("spec.phases.0.context")
		cp.Spec.Phases = append(cp.Spec.Phases, kuberoApi.Phase{
			Name:    "review",
			Enabled: true,
			Context: promptLine("Context for reviewapps", fmt.Sprint(contextSimpleList), contextDefault),
		})
	} else {
		cp.Spec.Reviewapps = false
		cp.Spec.Phases = append(cp.Spec.Phases, kuberoApi.Phase{
			Name:    "review",
			Enabled: false,
			Context: "",
		})
	}

	phaseTest := promptLine("enable test", "[y,n]", "n")
	if phaseTest == "y" {
		contextDefault := pipelineConfig.GetString("spec.phases.1.context")
		cp.Spec.Phases = append(cp.Spec.Phases, kuberoApi.Phase{
			Name:    "test",
			Enabled: true,
			Context: promptLine("Context for test", fmt.Sprint(contextSimpleList), contextDefault),
		})
	} else {
		cp.Spec.Phases = append(cp.Spec.Phases, kuberoApi.Phase{
			Name:    "test",
			Enabled: false,
			Context: "",
		})
	}

	phaseStage := promptLine("enable stage", "[y,n]", "n")
	if phaseStage == "y" {
		contextDefault := pipelineConfig.GetString("spec.phases.2.context")
		cp.Spec.Phases = append(cp.Spec.Phases, kuberoApi.Phase{
			Name:    "stage",
			Enabled: true,
			Context: promptLine("Context for stage", fmt.Sprint(contextSimpleList), contextDefault),
		})
	} else {
		cp.Spec.Phases = append(cp.Spec.Phases, kuberoApi.Phase{
			Name:    "stage",
			Enabled: false,
			Context: "",
		})
	}

	phaseProduction := promptLine("enable production", "[y,n]", "y")
	//var phaseProductionContext string = ""
	if phaseProduction != "n" {
		contextDefault := pipelineConfig.GetString("spec.phases.3.context")
		cp.Spec.Phases = append(cp.Spec.Phases, kuberoApi.Phase{
			Name:    "production",
			Enabled: true,
			Context: promptLine("Context for production ", fmt.Sprint(contextSimpleList), contextDefault),
		})
	} else {
		cp.Spec.Phases = append(cp.Spec.Phases, kuberoApi.Phase{
			Name:    "production",
			Enabled: false,
			Context: "",
		})
	}

	return cp
}
