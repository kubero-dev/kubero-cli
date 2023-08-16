package kuberoCli

import (
	"encoding/json"
	"fmt"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var appsFetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch a existing app configuration in a pipeline",
	Run: func(cmd *cobra.Command, args []string) {

		app := appsFetchForm()
		a, appErr := client.Get("/api/cli/pipelines/" + app.Spec.Pipeline + "/" + app.Spec.Phase + "/" + app.Spec.Name)

		if appErr != nil {
			fmt.Println(appErr)
		} else {
			fmt.Println(a)
			cfmt.Println("{{App fetched successfully}}::green")
			json.Unmarshal(a.Body(), &app)
			fmt.Println(app)
			writeAppYaml(app)
		}
	},
}

func init() {
	appsFetchCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Skip asking for confirmation")
	appsFetchCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "Name of the pipeline")
	appsFetchCmd.Flags().StringVarP(&stage, "stage", "s", "", "Name of the stage")
	appsFetchCmd.Flags().StringVarP(&app, "app", "a", "", "Name of the app")
	appsCmd.AddCommand(appsFetchCmd)
}

func appsFetchForm() CreateApp {

	var ca CreateApp
	ca.APIVersion = "application.kubero.dev/v1alpha1"
	ca.Kind = "KuberoApp"

	if pipelineName == "" {
		pipelineName = pipelineConfig.GetString("spec.name")
	}
	ca.Spec.Pipeline = promptLine("Pipeline", "", pipelineName)

	/* TODO need remote Phases
	availablePhases := getPipelinePhases()
	ca.Spec.Phase = promptLine("Phase", fmt.Sprint(availablePhases), stage)
	*/
	ca.Spec.Phase = promptLine("Phase", "[review,test,stage,production]", stage)

	if app == "" {
		appconfig := loadAppConfig(ca.Spec.Phase)
		app = appconfig.GetString("spec.name")
	}
	ca.Spec.Name = promptLine("Name", "", app)

	return ca
}
