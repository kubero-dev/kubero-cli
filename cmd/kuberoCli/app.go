package kuberoCli

import (
	"encoding/json"
	"kubero/pkg/kuberoApi"
	"log"
	"os"
	"strings"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
)

func appsList() {

	pipelineResp, _ := client.Get("/api/cli/pipelines/" + pipelineName + "/apps")

	var pl Pipeline
	json.Unmarshal(pipelineResp.Body(), &pl)

	for _, phase := range pl.Phases {
		if !phase.Enabled {
			continue
		}
		cfmt.Print("\n")

		cfmt.Println("{{  " + strings.ToUpper(phase.Name) + "}}::bold|white" + " (" + phase.Context + ")")

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
	json.Unmarshal(apps.Body(), &appShortList)

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
	basePath := "/.kubero/"
	gitdir := getGitdir()
	dir := gitdir + basePath + pipelineName + "/" + stageName

	var appsList []string
	apps, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, app := range apps {
		a := loadLocalApp(pipelineName, stageName, app.Name())
		if a.Kind == "KuberoApp" && a.Spec.Name != "" {
			appsList = append(appsList, a.Spec.Name)
		}
	}
	return appsList
}

func loadLocalApp(pipelineName string, stageName string, appName string) kuberoApi.AppCRD {

	appConfig := loadAppConfig(pipelineName, stageName, appName)

	var appCRD kuberoApi.AppCRD

	appConfig.Unmarshal(&appCRD)

	return appCRD
}

func loadAppConfig(pipelineName string, stageName string, appName string) *viper.Viper {

	basePath := "/.kubero/"
	gitdir := getGitdir()
	dir := gitdir + basePath + pipelineName + "/" + stageName
	//fmt.Println(dir)

	appConfig := viper.New()
	appConfig.SetConfigName(appName) // name of config file (without extension)
	appConfig.SetConfigType("yaml")  // REQUIRED if the config file does not have the extension in the name
	appConfig.AddConfigPath(dir)     // path to look for the config file in
	appConfig.ReadInConfig()

	return appConfig
}
