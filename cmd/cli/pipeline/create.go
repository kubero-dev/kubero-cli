package pipeline

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"fmt"
	"github.com/faelmori/kubero-cli/internal/pipeline"
	"github.com/faelmori/kubero-cli/types"

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

var createPipelineCmd = &cobra.Command{
	Use:     "pipeline",
	Aliases: []string{"pl"},
	Short:   "Create a new pipeline",
	Long:    `Create a new Pipeline`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create a new pipeline")

		_ = createPipeline()
	},
}

var createAppCmd = &cobra.Command{
	Use:   "app",
	Short: "Create a new app in a Pipeline",
	Long: `Create a new app in a Pipeline.

If called without arguments, it will ask for all the required information`,
	Run: func(cmd *cobra.Command, args []string) {
		c := pipeline.NewPipelineManager(pipelineName, stageName, appName)

		pipelinesList := c.getAllLocalPipelines()
		ensurePipelineIsSet(pipelinesList)
		ensureStageNameIsSet()
		ensureAppNameIsSet()
		createApp()

	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.PersistentFlags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
}

func createPipelineAndApp() {
	createPipelineAndAppArg := promptLine("Create a new pipeline", "[y,n]", "y")
	if createPipelineAndAppArg == "y" {
		createPipeline()
	}

	appConfig := pipeline.NewPipelineManager(pipelineName, stageName, appName)

	pipelinesList := appConfig.GetAllLocalApps()
	ensurePipelineIsSet(pipelinesList)
	ensureStageNameIsSet()
	ensureAppNameIsSet()
	createApp()
}

func appForm() types.AppCRD {

	var appCRD types.AppCRD

	appConfig := pipeline.NewPipelineManager(pipelineName, stageName, appName)
	plConfig := appConfig.LoadAppConfig(pipelineName, stageName, appName)

	appCRD.APIVersion = "application.kubero.dev/v1alpha1"
	appCRD.Kind = "KuberoApp"

	appCRD.Spec.Name = appName
	appCRD.Spec.Pipeline = pipelineName
	appCRD.Spec.Phase = stageName

	appCRD.Spec.Domain = promptLine("Domain", "", plConfig.GetString("spec.domain"))

	unmarshalKeyErr := plConfig.UnmarshalKey("spec.git.repository", &appCRD.Spec.Gitrepo)
	if unmarshalKeyErr != nil {
		fmt.Println(unmarshalKeyErr)
		return types.AppCRD{}
	}

	gitURL := pipelineConfig.GetString("spec.git.repository.sshurl")
	appCRD.Spec.Branch = promptLine("Branch", gitURL+":", plConfig.GetString("spec.branch"))

	appCRD.Spec.Buildpack = pipelineConfig.GetString("spec.buildpack.name")

	autodeployDefault := "n"
	if !plConfig.GetBool("spec.autodeploy") {
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

	appCRD.Spec.Image.ContainerPort, _ = strconv.Atoi(promptLine("Container Port", "8080", plConfig.GetString("spec.image.containerport")))

	appCRD.Spec.Web = types.Web{}
	appCRD.Spec.Web.ReplicaCount, _ = strconv.Atoi(promptLine("Web Pods", "1", plConfig.GetString("spec.web.replicacount")))

	appCRD.Spec.Worker = types.Worker{}
	appCRD.Spec.Worker.ReplicaCount, _ = strconv.Atoi(promptLine("Worker Pods", "0", plConfig.GetString("spec.worker.replicacount")))

	return appCRD
}

func createApp() {

	appCRD := appForm()

	writeAppYaml(appCRD)

	_, _ = cfmt.Println("\n\n{{Created appCRD.yaml}}::green")
}

func writeAppYaml(appCRD types.AppCRD) {
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

func createPipeline() types.PipelineCRD {

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

func writePipelineYaml(pipeline types.PipelineCRD) {
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

func pipelinesForm() types.PipelineCRD {

	var pipelineCRD types.PipelineCRD

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
		return types.PipelineCRD{}
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
			pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, types.Phase{
				Name:    "review",
				Enabled: true,
				Context: promptLine("Context for reviewapps", fmt.Sprint(contextSimpleList), contextDefault),
			})
		} else {
			pipelineCRD.Spec.ReviewApps = false
			pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, types.Phase{
				Name:    "review",
				Enabled: false,
				Context: "",
			})
		}
	}

	phaseTest := promptLine("enable test", "[y,n]", "n")
	if phaseTest == "y" {
		pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, types.Phase{
			Name:    "test",
			Enabled: true,
			Context: promptLine("Context for test", fmt.Sprint(contextSimpleList), contextDefault),
		})
	} else {
		pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, types.Phase{
			Name:    "test",
			Enabled: false,
			Context: "",
		})
	}

	phaseStage := promptLine("enable stage", "[y,n]", "n")
	if phaseStage == "y" {
		pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, types.Phase{
			Name:    "stage",
			Enabled: true,
			Context: promptLine("Context for stage", fmt.Sprint(contextSimpleList), contextDefault),
		})
	} else {
		pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, types.Phase{
			Name:    "stage",
			Enabled: false,
			Context: "",
		})
	}

	phaseProduction := promptLine("enable production", "[y,n]", "y")
	if phaseProduction != "n" {
		pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, types.Phase{
			Name:    "production",
			Enabled: true,
			Context: promptLine("Context for production ", fmt.Sprint(contextSimpleList), contextDefault),
		})
	} else {
		pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, types.Phase{
			Name:    "production",
			Enabled: false,
			Context: "",
		})
	}

	return pipelineCRD
}
