package pipeline

import (
	a "github.com/kubero-dev/kubero-cli/internal/api"
	c "github.com/kubero-dev/kubero-cli/internal/config"
	l "github.com/kubero-dev/kubero-cli/internal/log"
	u "github.com/kubero-dev/kubero-cli/internal/utils"
	t "github.com/kubero-dev/kubero-cli/types"

	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
)

var (
	utilsPrompt      = u.NewConsolePrompt()
	promptLine       = utilsPrompt.PromptLine
	confirmationLine = utilsPrompt.ConfirmationLine
	utils            = u.NewUtils()
)

type PipelineManager struct {
	c.ConfigManager
	repo         *a.Repository
	repositories []*t.GitRepository
	contexts     []*a.Context
	pipelines    *t.PipelinesList
	buildPacks   *t.BuildPacks
	pipelineName string
	stageName    string
	appName      string
}

func NewPipelineManager(pipelineName, stageName, appName string) *PipelineManager {
	return &PipelineManager{
		pipelineName: pipelineName,
		stageName:    stageName,
		appName:      appName,
	}
}

func (m *PipelineManager) GetPipelineConfig(pipelineName string) *viper.Viper {
	return m.loadPipelineConfig(pipelineName)
}

func (m *PipelineManager) LoadAllLocalPipelines() t.PipelinesConfigsList {
	pipelines := m.GetAllLocalPipelines()

	pipelinesConfigsList := make(t.PipelinesConfigsList)

	for _, pipeline := range pipelines {
		pipelinesConfigsList[pipeline] = m.loadLocalPipeline(pipeline)
	}

	return pipelinesConfigsList

}

func (m *PipelineManager) printPipeline(r *resty.Response) {
	//fmt.Println(r)

	var pipeline t.Pipeline
	unmarshalErr := json.Unmarshal(r.Body(), &pipeline)
	if unmarshalErr != nil {
		fmt.Println("Error: ", "Unable to decode response")
		return
	}

	_, _ = cfmt.Printf("{{Name:}}::lightWhite %v \n", pipeline.Name)
	_, _ = cfmt.Printf("{{BuildPack:}}::lightWhite %v\n", pipeline.BuildPack.Name)
	_, _ = cfmt.Printf("{{Language:}}::lightWhite %v\n", pipeline.BuildPack.Language)
	if pipeline.DockerImage != "" {
		fmt.Printf("{{Docker Image:}}::lightWhite %v \n", pipeline.DockerImage)
	}
	_, _ = cfmt.Printf("{{Deployment Strategy:}}::lightWhite %v \n", pipeline.DeploymentStrategy)
	_, _ = cfmt.Printf("{{Git:}}::lightWhite %v:%v \n", pipeline.Git.Repository.SshUrl, pipeline.Git.Repository.DefaultBranch)
}

func (m *PipelineManager) printPipelinesList(r *resty.Response) error {

	table := tablewriter.NewWriter(os.Stdout)
	//table.SetAutoFormatHeaders(true)
	//table.SetBorder(false)
	//table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetHeader([]string{
		"Name",
		"Repository",
		//"Branch",
		"BuildPack",
		"reviewapps",
		"test",
		"staging",
		"production",
		//"Docker Image",
		//"Deployment Strategy",
		//"Review Apps"
	})

	var pipelinesList t.PipelinesList
	unmarshalErr := json.Unmarshal(r.Body(), &pipelinesList)
	if unmarshalErr != nil {
		l.Error("Unable to decode response")
		return unmarshalErr
	}

	for _, pipeline := range pipelinesList.Items {
		table.Append([]string{
			pipeline.Name,
			pipeline.Git.Repository.SshUrl,
			//pipeline.Git.Repository.DefaultBranch,
			pipeline.BuildPack.Name,
			//pipeline.DockerImage,
			//pipeline.DeploymentStrategy,
			//fmt.Sprintf("%t", pipeline.ReviewApps)
			utils.BoolToEmoji(pipeline.Phases[0].Enabled),
			utils.BoolToEmoji(pipeline.Phases[1].Enabled),
			utils.BoolToEmoji(pipeline.Phases[2].Enabled),
			utils.BoolToEmoji(pipeline.Phases[3].Enabled),
		})
	}

	utils.PrintCLI(table, r, "table")

	return nil
}

func (m *PipelineManager) GetAllRemotePipelines() []string {
	var pipelinesList t.PipelinesList

	api := a.NewClient()
	res, err := api.GetPipelines()
	if err != nil {
		fmt.Println("Error: ", "Unable to load pipelines")
		fmt.Println(err)
		os.Exit(1)
	}

	unmarshalErr := json.Unmarshal(res.Body(), &pipelinesList)
	if unmarshalErr != nil {
		fmt.Println("Error: ", "Unable to decode response")
		return nil
	}

	var pipelines []string
	pipelines = make([]string, 0)

	for _, pipeline := range pipelinesList.Items {
		pipelines = append(pipelines, pipeline.Name)
	}

	return pipelines
}

func (m *PipelineManager) GetAllLocalPipelines() []string {
	baseDir := m.GetIACBaseDir()
	dir := baseDir + "/" + m.pipelineName

	var pipelineNames []string
	pipelineNames = make([]string, 0)

	files, err := os.ReadDir(dir)
	if err != nil {
		l.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			if _, err := os.Stat(dir + "/" + f.Name() + "/pipeline.yaml"); err == nil {
				pipelineNames = append(pipelineNames, f.Name())
			}
		}
	}

	return pipelineNames
}

func (m *PipelineManager) getPipelinePhases(pipelineConfig *viper.Viper) []string {
	var phases []string

	//pipelineConfig := getPipelineConfig(pipelineName)

	phasesList := pipelineConfig.GetStringSlice("spec.phases")

	for p := range phasesList {
		enabled := pipelineConfig.GetBool("spec.phases." + strconv.Itoa(p) + ".enabled")
		if enabled {
			phases = append(phases, pipelineConfig.GetString("spec.phases."+strconv.Itoa(p)+".name"))
		}
	}
	return phases
}

func (m *PipelineManager) loadPipelineConfig(pipelineName string) *viper.Viper {

	baseDir := m.GetIACBaseDir()
	dir := baseDir + "/" + pipelineName

	pipelineConfig := viper.New()
	pipelineConfig.SetConfigName("pipeline") // name of config file (without extension)
	pipelineConfig.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name
	pipelineConfig.AddConfigPath(dir)        // path to look for the config file in
	readInConfigErr := pipelineConfig.ReadInConfig()
	if readInConfigErr != nil {
		fmt.Println("Error: ", "Unable to read config file")
		return nil
	}

	return pipelineConfig
}

func (m *PipelineManager) loadLocalPipeline(pipelineName string) t.PipelineCRD {

	pipelineConfig := m.loadPipelineConfig(pipelineName)

	var pipelineCRD t.PipelineCRD

	pipelineConfigUnmarshalErr := pipelineConfig.Unmarshal(&pipelineCRD)
	if pipelineConfigUnmarshalErr != nil {
		fmt.Println("Error: ", "Unable to unmarshal config file")
		return t.PipelineCRD{}
	}

	return pipelineCRD
}
