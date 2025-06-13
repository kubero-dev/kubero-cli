package kuberoCli

import (
	"encoding/json"
	"fmt"
	"kubero/pkg/kuberoApi"
	"log"
	"os"
	"strings"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var AppCmd = &cobra.Command{
	Use:     "app",
	Aliases: []string{"apps", "application", "applications", "a"},
	Short:   "List apps in a Pipeline",
	Long: `Create a new app in a Pipeline.

If called without arguments, it will ask for all the required information`,
	Run: func(cmd *cobra.Command, args []string) {

		pipelinesList := getAllRemotePipelines()
		fmt.Println(pipelinesList)
		ensurePipelineIsSet(pipelinesList)
		//ensureStageNameIsSet()
		//ensureAppNameIsSet()
		appsList()

	},
}

func init() {
	rootCmd.AddCommand(AppCmd)
	AppCmd.PersistentFlags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
}

func appsList() {

	pipelineResp, _ := api.GetPipelineApps(pipelineName)

	var pl Pipeline
	jsonUnmarshalErr := json.Unmarshal(pipelineResp.Body(), &pl)
	if jsonUnmarshalErr != nil {
		log.Fatal("appsList ", jsonUnmarshalErr)
		return
	}

	for _, phase := range pl.Phases {
		if !phase.Enabled {
			continue
		}
		_, _ = cfmt.Print("\n")

		_, _ = cfmt.Println("{{  " + strings.ToUpper(phase.Name) + "}}::bold|white" + " (" + phase.Context + ")")

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{
			"Name",
			"Phase",
			"Pipeline",
			"Repository",
			"Domain",
		})
		table.SetBorder(false)

		for _, app := range phase.Apps {
			table.Append([]string{
				app.Name,
				app.Phase,
				app.Pipeline,
				app.Gitrepo.CloneURL + ":" +
					app.Gitrepo.DefaultBranch,
				app.Domain,
			})
		}

		printCLI(table, pipelineResp)
	}
}

func getAllRemoteApps() []string {
	apps, _ := api.GetApps()
	var appShortList []appShort
	jsonUnmarshalErr := json.Unmarshal(apps.Body(), &appShortList)
	if jsonUnmarshalErr != nil {
		log.Fatal(jsonUnmarshalErr)
		return nil
	}

	var appsList []string
	for _, app := range appShortList {
		if pipelineName != "" && app.Pipeline != pipelineName {
			continue
		}
		if stageName != "" && app.Phase != stageName {
			continue
		}
		if appName != "" && app.Name != appName {
			continue
		}
		appsList = append(appsList, app.Name)
	}

	return appsList
}

func getAllLocalApps() []string {

	baseDir := getIACBaseDir()
	dir := baseDir + "/" + pipelineName + "/" + stageName

	var appsList []string
	appFiles, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, appFileName := range appFiles {

		// remove the .yaml extension
		appName := strings.TrimSuffix(appFileName.Name(), ".yaml")

		a := loadLocalApp(pipelineName, stageName, appName)
		if a.Kind == "KuberoApp" && a.Metadata.Name != "" {
			appsList = append(appsList, a.Metadata.Name)
		}
	}
	return appsList
}

func loadLocalApp(pipelineName string, stageName string, appName string) kuberoApi.AppCRD {

	appConfig := loadAppConfig(pipelineName, stageName, appName)

	var appCRD kuberoApi.AppCRD

	appConfigUnmarshalErr := appConfig.Unmarshal(&appCRD)
	if appConfigUnmarshalErr != nil {
		log.Fatal(appConfigUnmarshalErr)
		return kuberoApi.AppCRD{}
	}

	return appCRD
}

func loadAppConfig(pipelineName string, stageName string, appName string) *viper.Viper {

	baseDir := getIACBaseDir()
	dir := baseDir + "/" + pipelineName + "/" + stageName

	appConfig := viper.New()
	appConfig.SetConfigName(appName)
	appConfig.SetConfigType("yaml")
	appConfig.AddConfigPath(dir)
	readInConfigErr := appConfig.ReadInConfig()
	if readInConfigErr != nil {
		log.Fatal(readInConfigErr)
		return nil
	}

	return appConfig
}
