package kuberoCli

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"encoding/json"
	"fmt"
	"kubero/pkg/kuberoApi"
	"os"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/cobra"
)

// appCmd represents the app command
var fetchAppCmd = &cobra.Command{
	Use:   "iac:fetch",
	Short: "Fetch an app configuration",
	Long:  `Fetch an app`,
	Run: func(cmd *cobra.Command, args []string) {

		pipelinesList := getAllRemotePipelines()
		if len(pipelinesList) == 0 {
			_, _ = cfmt.Println("\n{{ERROR:}}::red No pipelines found")
			os.Exit(1)
		}
		ensurePipelineIsSet(pipelinesList)
		ensureStageNameIsSet()
		fetchPipeline(pipelineName)

		appsList := getAllRemoteApps()
		if len(appsList) == 0 {
			_, _ = cfmt.Println("\n{{ERROR:}}::red No apps found in pipeline '" + pipelineName + "'")
			os.Exit(1)
		}

		ensureAppNameIsSelected(appsList)
		ensureAppNameIsSet()
		fetchApp(appName, stageName, pipelineName)
	},
}

func init() {
	AppCmd.AddCommand(fetchAppCmd)
	fetchAppCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "Name of the pipeline")
	fetchAppCmd.Flags().StringVarP(&stageName, "stage", "s", "", "Name of the stage [test|stage|production]")
	fetchAppCmd.Flags().StringVarP(&appName, "app", "a", "", "Name of the app")
}

func fetchApp(appName string, stageName string, pipelineName string) {

	confirmation := promptLine("Do you want to fetch the app '"+appName+"' from '"+pipelineName+"'?", "[y,n]", "y")
	if confirmation == "y" {
		_, _ = cfmt.Println("{{Fetching app}}::yellow " + appName + "")
	} else {
		_, _ = cfmt.Println("{{Aborted}}::red")
		return
	}

	var app kuberoApi.AppCRD
	app.APIVersion = "application.kubero.dev/v1alpha1"
	app.Kind = "KuberoApp"

	app.Spec.Pipeline = pipelineName
	app.Spec.Phase = stageName
	app.Spec.Name = appName
	app.Metadata.Name = appName

	a, appErr := api.GetApp(pipelineName, stageName, appName)

	if appErr != nil {
		if a == nil {
			_, _ = cfmt.Println("{{ERROR:}}::red App '" + appName + "' not found ")
			os.Exit(1)
		}
		if a.StatusCode() == 404 {
			_, _ = cfmt.Println("{{ERROR:}}::red App '" + appName + "' not found ")
			os.Exit(1)
		}
		fmt.Println(appErr)
		os.Exit(1)
	}

	jsonUnmarshalErr := json.Unmarshal(a.Body(), &app)
	if jsonUnmarshalErr != nil {
		fmt.Println(jsonUnmarshalErr)
		return
	}

	writeAppYaml(app)

}
