package pipeline

import (
	"encoding/json"
	"github.com/faelmori/kubero-cli/internal/log"
	"github.com/faelmori/kubero-cli/internal/utils"
	"github.com/faelmori/kubero-cli/pkg/kuberoApi"
	"github.com/faelmori/kubero-cli/types"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type ManagerPipeline struct {
	pipelineName string
	stageName    string
	appName      string
}

func NewPipelineManager(pipelineName, stageName, appName string) *ManagerPipeline {
	return &ManagerPipeline{
		pipelineName: pipelineName,
		stageName:    stageName,
		appName:      appName,
	}
}

func (m *ManagerPipeline) AppsList(pipelineName, outputFormat string) error {
	api := kuberoApi.NewKuberoClient()
	pipelineResp, _ := api.GetPipelineApps(pipelineName)

	var pl types.PipelineSpec
	jsonUnmarshalErr := json.Unmarshal(pipelineResp.Body(), &pl)
	if jsonUnmarshalErr != nil {
		log.Error("Unable to decode response")
		return jsonUnmarshalErr
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

		utils.PrintCLI(table, pipelineResp, outputFormat)
	}
}

func (m *ManagerPipeline) GetAllRemoteApps() []string {
	api := kuberoApi.NewKuberoClient()
	apps, _ := api.GetApps()
	var appShortList []types.AppShort
	jsonUnmarshalErr := json.Unmarshal(apps.Body(), &appShortList)
	if jsonUnmarshalErr != nil {
		log.Fatal(jsonUnmarshalErr)
		return nil
	}

	var appsList []string
	for _, app := range appShortList {
		if m.pipelineName != "" && api.Pipeline != m.pipelineName {
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

func (m *ManagerPipeline) GetAllLocalApps() []string {

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

func (m *ManagerPipeline) LoadLocalApp(pipelineName string, stageName string, appName string) kuberoApi.AppCRD {

	appConfig := loadAppConfig(pipelineName, stageName, appName)

	var appCRD kuberoApi.AppCRD

	appConfigUnmarshalErr := appConfig.Unmarshal(&appCRD)
	if appConfigUnmarshalErr != nil {
		log.Fatal(appConfigUnmarshalErr)
		return kuberoApi.AppCRD{}
	}

	return appCRD
}

func (m *ManagerPipeline) LoadAppConfig(pipelineName string, stageName string, appName string) *viper.Viper {

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
