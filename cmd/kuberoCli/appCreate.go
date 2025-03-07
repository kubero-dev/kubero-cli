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

// appCmd represents the app command
var createAppCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new app in a Pipeline",
	Long: `Create a new app in a Pipeline.

If called without arguments, it will ask for all the required information`,
	Run: func(cmd *cobra.Command, args []string) {

		pipelinesList := getAllRemotePipelines()
		ensurePipelineIsSet(pipelinesList)
		ensureStageNameIsSet()
		ensureAppNameIsSet()
		createRemoteApp()

	},
}

func init() {
	AppCmd.AddCommand(createAppCmd)
	createAppCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "Name of the pipeline")
	createAppCmd.Flags().StringVarP(&stageName, "stage", "s", "", "Name of the stage")
	createAppCmd.Flags().StringVarP(&appName, "app", "a", "", "Name of the app")
}

func appForm() kuberoApi.AppCRD {

	var appCRD kuberoApi.AppCRD

	pipelineConfig := loadPipelineConfig(pipelineName, false)

	appCRD.APIVersion = "application.kubero.dev/v1alpha1"
	appCRD.Kind = "KuberoApp"

	appCRD.Spec.Name = appName
	appCRD.Spec.Pipeline = pipelineName
	appCRD.Spec.Phase = stageName

	appCRD.Spec.Domain = promptLine("Domain", "", "")
	gitURL := pipelineConfig.Spec.Git.Repository.SshUrl
	if gitURL != "" {
		appCRD.Spec.Branch = promptLine("Branch", gitURL+":", "main")
	} else {
		appCRD.Spec.Branch = promptLine("Branch", "", "main")
	}

	appCRD.Spec.Buildpack = pipelineConfig.Spec.Buildpack.Name

	autodeploy := promptLine("Autodeploy", "[y,n]", "n")
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

	appCRD.Spec.Image.ContainerPort, _ = strconv.Atoi(promptLine("Container Port", "8080", "8080"))

	appCRD.Spec.Web = kuberoApi.Web{}
	appCRD.Spec.Web.ReplicaCount, _ = strconv.Atoi(promptLine("Web Pods", "1", "1"))

	appCRD.Spec.Worker = kuberoApi.Worker{}
	appCRD.Spec.Worker.ReplicaCount, _ = strconv.Atoi(promptLine("Worker Pods", "0", "0"))

	return appCRD
}

func createRemoteApp() {
	appCRD := appForm()

	_, err := api.DeployApp(pipelineName, stageName, appName, appCRD)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		_, _ = cfmt.Print("\n{{Created appCRD}}::green\n\n")
	}

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
