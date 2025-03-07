package kuberoCli

import (
	"encoding/json"
	"fmt"
	"kubero/pkg/kuberoApi"
	"log"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pipelineCmd = &cobra.Command{
	Use:     "pipeline",
	Aliases: []string{"p", "pipelines"},
	Short:   "List pipelines",
	Long:    `List pipelines`,
	Run: func(cmd *cobra.Command, args []string) {

		if pipelineName != "" {
			// get a single pipeline

			pipelineResp, err := api.GetPipeline(pipelineName)
			//pipelineResp, err := client.Get("/api/cli/pipelines/" + pipelineName)
			if pipelineResp.StatusCode() == 404 {
				_, _ = cfmt.Println("{{  Pipeline not found}}::red")
				os.Exit(1)
			}

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			printPipeline(pipelineResp)
			appsList()
		} else {
			// get the pipelines
			pipelineListResp, err := api.GetPipelines()
			//pipelineListResp, err := client.Get("/api/cli/pipelines")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			printPipelinesList(pipelineListResp)

		}
	},
}

func init() {
	rootCmd.AddCommand(pipelineCmd)
	pipelineCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "output format [table, json]")
}

func loadAllLocalPipelines() pipelinesConfigsList {
	pipelines := getAllLocalPipelines()

	//var pipelinesConfigsList pipelinesConfigsList
	pipelinesConfigsList := make(pipelinesConfigsList)

	for _, pipeline := range pipelines {
		pipelinesConfigsList[pipeline] = loadPipelineConfig(pipeline, true)
	}

	return pipelinesConfigsList

}

func printPipeline(r *resty.Response) {
	//fmt.Println(r)

	var pipeline Pipeline
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

// print the response as a table
func printPipelinesList(r *resty.Response) {

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

	var pipelinesList PipelinesList
	unmarshalErr := json.Unmarshal(r.Body(), &pipelinesList)
	if unmarshalErr != nil {
		fmt.Println("Error: ", "Unable to decode response")
		return
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
			boolToEmoji(pipeline.Phases[0].Enabled),
			boolToEmoji(pipeline.Phases[1].Enabled),
			boolToEmoji(pipeline.Phases[2].Enabled),
			boolToEmoji(pipeline.Phases[3].Enabled),
		})
	}

	printCLI(table, r)
}

func getAllRemotePipelines() []string {
	var pipelinesList PipelinesList

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

func getAllLocalPipelines() []string {

	baseDir := getIACBaseDir()
	dir := baseDir + "/" + pipelineName

	var pipelineNames []string
	pipelineNames = make([]string, 0)

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
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

func getPipelinePhases(pipelineConfig *viper.Viper) []string {
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

func getPipelinePhasesFromCRD(pipelineCRD kuberoApi.PipelineCRD) []string {
	var phases []string

	phasesList := pipelineCRD.Spec.Phases

	for p := range phasesList {
		if phasesList[p].Enabled {
			phases = append(phases, phasesList[p].Name)
		}
	}
	return phases
}

func loadLocalPipelineConfig(pipelineName string) *viper.Viper {

	baseDir := getIACBaseDir()
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

func loadPipelineConfig(pipelineName string, local bool) kuberoApi.PipelineCRD {
	var pipelineConfig *viper.Viper
	if local {
		pipelineConfig = loadLocalPipelineConfig(pipelineName)
	} else {
		pipelineConfig = loadRemotePipelineConfig(pipelineName)
	}

	var pipelineCRD kuberoApi.PipelineCRD

	pipelineConfigUnmarshalErr := pipelineConfig.Unmarshal(&pipelineCRD)
	if pipelineConfigUnmarshalErr != nil {
		fmt.Println("Error: ", "Unable to unmarshal config file")
		return kuberoApi.PipelineCRD{}
	}

	return pipelineCRD
}

func loadRemotePipelineConfig(pipelineName string) *viper.Viper {
	res, err := api.GetPipeline(pipelineName)
	if err != nil {
		fmt.Println("Error: ", "Unable to load pipeline")
		fmt.Println(err)
		os.Exit(1)
	}

	var pipeline kuberoApi.PipelineSpec
	unmarshalErr := json.Unmarshal(res.Body(), &pipeline)
	if unmarshalErr != nil {
		fmt.Println("Error: ", "Unable to decode response")
		return nil
	}

	pipelineConfig := viper.New()
	pipelineConfig.SetConfigType("yaml")
	pipelineConfig.SetConfigName("pipeline")
	pipelineConfig.Set("spec", pipeline)

	return pipelineConfig
}
