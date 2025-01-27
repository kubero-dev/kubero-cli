package kuberoCli

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"fmt"
	"kubero/pkg/kuberoApi"

	"os"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
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
	}

	pipelinesList := getAllLocalPipelines()
	ensurePipelineIsSet(pipelinesList)
	ensureStageNameIsSet()
	ensureAppNameIsSet()
	createApp()
}

func appForm() kuberoApi.AppCRD {

	var appCRD kuberoApi.AppCRD

	appconfig := loadAppConfig(pipelineName, stageName, appName)
	pipelineConfig := loadPipelineConfig(pipelineName)

	appCRD.APIVersion = "application.kubero.dev/v1alpha1"
	appCRD.Kind = "KuberoApp"

	appCRD.Spec.Name = appName
	appCRD.Spec.Pipeline = pipelineName
	appCRD.Spec.Phase = stageName

	appCRD.Spec.Domain = promptLine("Domain", "", appconfig.GetString("spec.domain"))

	unmarshalKeyErr := pipelineConfig.UnmarshalKey("spec.git.repository", &appCRD.Spec.Gitrepo)
	if unmarshalKeyErr != nil {
		fmt.Println(unmarshalKeyErr)
		return kuberoApi.AppCRD{}
	}

	gitURL := pipelineConfig.GetString("spec.git.repository.sshurl")
	appCRD.Spec.Branch = promptLine("Branch", gitURL+":", appconfig.GetString("spec.branch"))

	appCRD.Spec.Buildpack = pipelineConfig.GetString("spec.buildpack.name")

	autodeployDefault := "n"
	if !appconfig.GetBool("spec.autodeploy") {
		autodeployDefault = "y"
	}
	autodeploy := promptLine("Autodeploy", "[y,n]", autodeployDefault)
	if autodeploy == "Y" {
		appCRD.Spec.Autodeploy = true
	} else {
		appCRD.Spec.Autodeploy = false
	}

	envCount, _ := strconv.Atoi(promptLine("Number of Env Vars", "", "0"))
	appCRD.Spec.EnvVars = []string{}
	for i := 0; i < envCount; i++ {
		appCRD.Spec.EnvVars = append(appCRD.Spec.EnvVars, promptLine("Env Var", "", ""))
	}

	appCRD.Spec.Image.ContainerPort, _ = strconv.Atoi(promptLine("Container Port", "8080", appconfig.GetString("spec.image.containerport")))

	appCRD.Spec.Web = kuberoApi.Web{}
	appCRD.Spec.Web.ReplicaCount, _ = strconv.Atoi(promptLine("Web Pods", "1", appconfig.GetString("spec.web.replicacount")))

	appCRD.Spec.Worker = kuberoApi.Worker{}
	appCRD.Spec.Worker.ReplicaCount, _ = strconv.Atoi(promptLine("Worker Pods", "0", appconfig.GetString("spec.worker.replicacount")))

	return appCRD
}

func createApp() {

	appCRD := appForm()

	writeAppYaml(appCRD)

	_, _ = cfmt.Println("\n\n{{Created appCRD.yaml}}::green")
}

func writeAppYaml(appCRD kuberoApi.AppCRD) {
	// write pipeline.yaml
	yamlData, err := yaml.Marshal(&appCRD)

	if err != nil || appCRD.Spec.Name == "" {
		panic("Unable to write data into the file")
	}

	fileName := ".kubero/" + appCRD.Spec.Pipeline + "/" + appCRD.Spec.Phase + "/" + appCRD.Spec.Name + ".yaml"

	err = os.WriteFile(fileName, yamlData, 0644)
	if err != nil {
		panic("Unable to write data into the file")
	}
}

func createPipeline() kuberoApi.PipelineCRD {

	loadConfigs(pipelineName)

	loadRepositories()
	loadContexts()
	loadBuildpacks()
	pipelineCRD := pipelinesForm()

	writePipelineYaml(pipelineCRD)

	_, _ = cfmt.Println("\n\n{{Created pipeline.yaml}}::green")
	_, _ = cfmt.Println(pipelineName)

	return pipelineCRD
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

	var pipelineCRD kuberoApi.PipelineCRD

	if pipelineName == "" {
		pipelineName = promptLine("Define a PipelineName name", "", "")
	}
	pipelineCRD.Spec.Name = pipelineName

	pipelineCRD.APIVersion = "application.kubero.dev/v1alpha1"
	pipelineCRD.Kind = "KuberoPipeline"

	fmt.Println("")
	prompt := &survey.Select{
		Message: "Select a buildpack",
		Options: buildPacksSimpleList,
	}
	askOneErr := survey.AskOne(prompt, &pipelineCRD.Spec.Buildpack.Name)
	if askOneErr != nil {
		fmt.Println(askOneErr.Error())
		return kuberoApi.PipelineCRD{}
	}

	domain := pipelineConfig.GetString("spec.domain")
	pipelineCRD.Spec.Domain = promptLine("FQDN Domain ", "", domain)

	// those fields are deprecated and may be removed in the future
	pipelineCRD.Spec.DockerImage = ""
	pipelineCRD.Spec.DeploymentStrategy = "git"

	gitconnection := promptLine("Connect pipeline to a Git repository (GitOps)", "[y,n]", "n")

	contextDefault := contextSimpleList[0]
	if gitconnection == "y" {
		gitProvider := pipelineConfig.GetString("spec.git.repository.provider")
		pipelineCRD.Spec.Git.Repository.Provider = promptLine("Repository Provider", fmt.Sprint(repoSimpleList), gitProvider)

		gitURL := pipelineConfig.GetString("spec.git.repository.sshurl")
		pipelineCRD.Spec.Git.Repository.SshUrl = promptLine("Repository URL", "["+getGitRemote()+"]", gitURL)

		phaseReview := promptLine("enable reviewapps", "[y,n]", "n")
		if phaseReview == "y" {
			pipelineCRD.Spec.ReviewApps = true
			pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, kuberoApi.Phase{
				Name:    "review",
				Enabled: true,
				Context: promptLine("Context for reviewapps", fmt.Sprint(contextSimpleList), contextDefault),
			})
		} else {
			pipelineCRD.Spec.ReviewApps = false
			pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, kuberoApi.Phase{
				Name:    "review",
				Enabled: false,
				Context: "",
			})
		}
	}

	phaseTest := promptLine("enable test", "[y,n]", "n")
	if phaseTest == "y" {
		pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, kuberoApi.Phase{
			Name:    "test",
			Enabled: true,
			Context: promptLine("Context for test", fmt.Sprint(contextSimpleList), contextDefault),
		})
	} else {
		pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, kuberoApi.Phase{
			Name:    "test",
			Enabled: false,
			Context: "",
		})
	}

	phaseStage := promptLine("enable stage", "[y,n]", "n")
	if phaseStage == "y" {
		pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, kuberoApi.Phase{
			Name:    "stage",
			Enabled: true,
			Context: promptLine("Context for stage", fmt.Sprint(contextSimpleList), contextDefault),
		})
	} else {
		pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, kuberoApi.Phase{
			Name:    "stage",
			Enabled: false,
			Context: "",
		})
	}

	phaseProduction := promptLine("enable production", "[y,n]", "y")
	if phaseProduction != "n" {
		pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, kuberoApi.Phase{
			Name:    "production",
			Enabled: true,
			Context: promptLine("Context for production ", fmt.Sprint(contextSimpleList), contextDefault),
		})
	} else {
		pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, kuberoApi.Phase{
			Name:    "production",
			Enabled: false,
			Context: "",
		})
	}

	return pipelineCRD
}
