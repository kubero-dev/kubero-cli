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
	"github.com/spf13/viper"
)

func loadAllLocalPipelines() pipelinesConfigsList {
	pipelines := getAllLocalPipelines()

	//var pipelinesConfigsList pipelinesConfigsList
	pipelinesConfigsList := make(pipelinesConfigsList)

	for _, pipeline := range pipelines {
		pipelinesConfigsList[pipeline] = loadLocalPipeline(pipeline)
	}

	return pipelinesConfigsList

}

func printPipeline(r *resty.Response) {
	//fmt.Println(r)

	var pipeline Pipeline
	json.Unmarshal(r.Body(), &pipeline)

	cfmt.Printf("{{Name:}}::lightWhite %v \n", pipeline.Name)
	cfmt.Printf("{{Buildpack:}}::lightWhite %v\n", pipeline.Buildpack.Name)
	cfmt.Printf("{{Language:}}::lightWhite %v\n", pipeline.Buildpack.Language)
	if pipeline.Dockerimage != "" {
		fmt.Printf("{{Docker Image:}}::lightWhite %v \n", pipeline.Dockerimage)
	}
	cfmt.Printf("{{Deployment Strategy:}}::lightWhite %v \n", pipeline.Deploymentstrategy)
	cfmt.Printf("{{Git:}}::lightWhite %v:%v \n", pipeline.Git.Repository.SSHURL, pipeline.Git.Repository.DefaultBranch)
	/*
		cfmt.Printf("{{Review Apps:}}::lightWhite %v \n", pipeline.Reviewapps)
		cfmt.Printf("{{Phases:}}::lightWhite \n")
		for _, phase := range pipeline.Phases {
			if phase.Enabled {
				fmt.Printf(" - %v (%v) \n", phase.Name, phase.Context)
			}
		}
	*/
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
		"Buildpack",
		"reviewapps",
		"test",
		"staging",
		"production",
		//"Docker Image",
		//"Deployment Strategy",
		//"Review Apps"
	})

	var pipelinesList PipelinesList
	json.Unmarshal(r.Body(), &pipelinesList)

	for _, pipeline := range pipelinesList.Items {
		table.Append([]string{
			pipeline.Name,
			pipeline.Git.Repository.SSHURL,
			//pipeline.Git.Repository.DefaultBranch,
			pipeline.Buildpack.Name,
			//pipeline.Dockerimage,
			//pipeline.Deploymentstrategy,
			//fmt.Sprintf("%t", pipeline.Reviewapps)
			boolToEmoji(pipeline.Phases[0].Enabled),
			boolToEmoji(pipeline.Phases[1].Enabled),
			boolToEmoji(pipeline.Phases[2].Enabled),
			boolToEmoji(pipeline.Phases[3].Enabled),
		})
	}

	printCLI(table, r)
}

func getAllLocalPipelines() []string {

	basePath := "/.kubero/"
	gitdir := getGitdir()
	dir := gitdir + basePath + pipelineName

	pipelineNames := []string{}
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

func loadPipelineConfig(pipelineName string) *viper.Viper {

	basePath := "/.kubero/"
	gitdir := getGitdir()
	dir := gitdir + basePath + pipelineName
	//fmt.Println(dir)

	pipelineConfig := viper.New()
	pipelineConfig.SetConfigName("pipeline") // name of config file (without extension)
	pipelineConfig.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name
	pipelineConfig.AddConfigPath(dir)        // path to look for the config file in
	pipelineConfig.ReadInConfig()

	return pipelineConfig
}

func loadLocalPipeline(pipelineName string) kuberoApi.PipelineCRD {

	pipelineConfig := loadPipelineConfig(pipelineName)

	var pipelineCRD kuberoApi.PipelineCRD

	pipelineConfig.Unmarshal(&pipelineCRD)

	return pipelineCRD
}
