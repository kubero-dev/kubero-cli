package pipeline

import (
	"encoding/json"
	a "github.com/faelmori/kubero-cli/internal/api"
	c "github.com/faelmori/kubero-cli/internal/config"
	"github.com/faelmori/kubero-cli/internal/log"
	"github.com/faelmori/kubero-cli/types"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func (m *PipelineManager) AppsList(pipelineName, outputFormat string) error {
	api := a.NewClient()
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

	return nil
}

func (m *PipelineManager) GetAllRemoteApps() []string {
	api := a.NewClient()
	apps, _ := api.GetApps()
	var appShortList []types.AppShort
	jsonUnmarshalErr := json.Unmarshal(apps.Body(), &appShortList)
	if jsonUnmarshalErr != nil {
		log.Fatal(jsonUnmarshalErr)
		return nil
	}

	var appsList []string
	//var pl types.Pipeline
	for _, app := range appShortList {
		if m.pipelineName != "" && app.Name != m.pipelineName {
			continue
		}
		if m.stageName != "" && app.Phase != m.stageName {
			continue
		}
		if m.appName != "" && app.Name != m.appName {
			continue
		}
		appsList = append(appsList, app.Name)
	}

	return appsList
}

func (m *PipelineManager) GetAllLocalApps() []string {
	cfg := c.NewViperConfig("", "")
	baseDir := cfg.GetIACBaseDir()
	dir := baseDir + "/" + m.pipelineName + "/" + m.stageName

	var appsList []string
	appFiles, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, appFileName := range appFiles {

		// remove the .yaml extension
		appName := strings.TrimSuffix(appFileName.Name(), ".yaml")

		ap := m.LoadLocalApp(m.pipelineName, m.stageName, appName)

		if ap.Kind == "KuberoApp" && ap.Metadata.Name != "" {
			appsList = append(appsList, ap.Metadata.Name)
		}
	}
	return appsList
}

func (m *PipelineManager) LoadLocalApp(pipelineName string, stageName string, appName string) types.AppCRD {

	appConfig := m.LoadAppConfig(pipelineName, stageName, appName)

	var appCRD types.AppCRD

	appConfigUnmarshalErr := appConfig.Unmarshal(&appCRD)
	if appConfigUnmarshalErr != nil {
		log.Fatal(appConfigUnmarshalErr)
		return types.AppCRD{}
	}

	return appCRD
}

func (m *PipelineManager) LoadAppConfig(pipelineName string, stageName string, appName string) *viper.Viper {
	cfg := c.NewViperConfig("", "")
	baseDir := cfg.GetIACBaseDir()

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
