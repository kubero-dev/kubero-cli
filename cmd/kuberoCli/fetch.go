package kuberoCli

import (
	"encoding/json"
	"fmt"
	"kubero/pkg/kuberoApi"
	"os"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:     "fetch",
	Aliases: []string{"pull", "fe"},
	Short:   "Fetch your remote pipelines and apps to your local repository",
	Long:    `Use the pipeline or app subcommand to fetch your pipelines and apps to your local repository`,
	Run: func(cmd *cobra.Command, args []string) {
		if pipelineName != "" && appName == "" {
			fetchPipeline(pipelineName)
		} else if pipelineName != "" && appName != "" {
			ensureStageNameIsSet()
			fetchPipeline(pipelineName)
			fetchApp(appName, stageName, pipelineName)
		} else {
			fetchAllPipelines()
		}
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "Name of the pipeline")
	fetchCmd.Flags().StringVarP(&stageName, "stage", "s", "", "Name of the stage [test|stage|production]")
	fetchCmd.Flags().StringVarP(&appName, "app", "a", "", "Name of the app")
}

func fetchPipeline(pipelineName string) {
	confirmation := promptLine("Do you want to fetch the pipeline '"+pipelineName+"'?", "[y,n]", "y")
	if confirmation == "y" {
		_, _ = cfmt.Println("{{Fetching pipeline}}::yellow " + pipelineName)

		var pipeline kuberoApi.PipelineCRD

		pipeline.APIVersion = "application.kubero.dev/v1alpha1"
		pipeline.Kind = "KuberoPipeline"
		pipeline.Spec.Name = pipelineName
		pipeline.Metadata.Name = appName

		p, pipelineErr := api.GetPipeline(pipelineName)

		if pipelineErr != nil {
			if p == nil {
				_, _ = cfmt.Println("{{ERROR:}}::red Pipeline '" + pipelineName + "' not found ")
				os.Exit(1)
			}
			if p.StatusCode() == 404 {
				_, _ = cfmt.Println("{{ERROR:}}::red Pipeline '" + pipelineName + "' not found ")
				os.Exit(1)
			}
			fmt.Println(pipelineErr)
			os.Exit(1)
		}

		jsonUnmarshalErr := json.Unmarshal(p.Body(), &pipeline.Spec)
		if jsonUnmarshalErr != nil {
			fmt.Println(jsonUnmarshalErr)
			return
		}
		writePipelineYaml(pipeline)

	} else {
		return
	}
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

func fetchAllPipelines() {
	confirmation := promptLine("Are you sure you want to fetch all pipelines?", "[y,n]", "n")
	if confirmation == "y" {
		_, _ = cfmt.Println("{{Fetching all pipelines}}::yellow")
	} else {
		_, _ = cfmt.Println("{{Aborted}}::red")
		return
	}
}
