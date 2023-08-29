package kuberoCli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// appCmd represents the app command
var createAppCmd = &cobra.Command{
	Use:   "app",
	Short: "Create a new app in a Pipeline",
	Long: `Create a new app in a Pipeline.

If called without arguments, it will ask for all the required information`,
	Run: func(cmd *cobra.Command, args []string) {

		createApp := appForm(appName, pipelineName)
		writeAppYaml(createApp)

	},
}

func init() {
	createCmd.AddCommand(createAppCmd)
	createAppCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "Name of the pipeline")
	createAppCmd.Flags().StringVarP(&stageName, "stage", "s", "", "Name of the stage")
	createAppCmd.Flags().StringVarP(&appName, "app", "a", "", "Name of the app")
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
