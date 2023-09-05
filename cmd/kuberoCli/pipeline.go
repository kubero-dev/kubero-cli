package kuberoCli

import (
	"encoding/json"
	"fmt"
	"kubero/pkg/kuberoApi"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
)

func listAllLocalPipelines() []string {
	gitdir := getGitdir()
	basePath := "/.kubero/" //TODO Make it dynamic
	dir := gitdir + basePath

	var pipelines []string
	files, _ := os.ReadDir(dir)
	for _, f := range files {
		if f.IsDir() {
			pipelines = append(pipelines, f.Name())
		}
	}
	return pipelines
}

func loadAllLocalPipelines() pipelinesConfigsList {
	pipelines := listAllLocalPipelines()

	//var pipelinesConfigsList pipelinesConfigsList
	pipelinesConfigsList := make(pipelinesConfigsList)

	for _, pipeline := range pipelines {
		pipelinesConfigsList[pipeline] = loadLocalPipeline(pipeline)
	}

	return pipelinesConfigsList

}

func loadLocalPipeline(pipelineName string) kuberoApi.PipelineCRD {

	gitdir := getGitdir()
	basePath := "/.kubero/" //TODO Make it dynamic
	dir := gitdir + basePath + pipelineName
	fmt.Println(dir)

	pipelineConfig := viper.New()
	pipelineConfig.SetConfigName("pipeline") // name of config file (without extension)
	pipelineConfig.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name
	pipelineConfig.AddConfigPath(dir)        // path to look for the config file in
	pipelineConfig.ReadInConfig()

	var pipelineCRD kuberoApi.PipelineCRD

	pipelineConfig.Unmarshal(&pipelineCRD)

	return pipelineCRD
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
