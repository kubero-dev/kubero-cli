package kuberoCli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// pipelineCmd represents the pipeline command
var listPipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "List available pipelines",
	Long:  `List available pipelines`,
	Run: func(cmd *cobra.Command, args []string) {

		if pipelineName != "" {
			// get a single pipeline
			pipelineResp, _ := client.Get("/api/cli/pipelines/" + pipelineName)
			printPipeline(pipelineResp)
		} else {
			// get the pipelines
			pipelineListResp, _ := client.Get("/api/cli/pipelines")
			printPipelinesList(pipelineListResp)
		}
	},
}

func init() {
	listCmd.AddCommand(listPipelineCmd)
	listPipelineCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
}

// print the response as a table
func printPipelinesList(r *resty.Response) {

	table := tablewriter.NewWriter(os.Stdout)
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
	//table.SetBorder(false)

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
			fmt.Sprintf("%t", pipeline.Phases[0].Enabled),
			fmt.Sprintf("%t", pipeline.Phases[1].Enabled),
			fmt.Sprintf("%t", pipeline.Phases[2].Enabled),
			fmt.Sprintf("%t", pipeline.Phases[3].Enabled),
		})
	}

	printCLI(table, r)
}

func printPipeline(r *resty.Response) {
	//fmt.Println(r)

	var pipeline Pipeline
	json.Unmarshal(r.Body(), &pipeline)

	cfmt.Printf("{{Name:}}::lightWhite %v \n", pipeline.Name)
	cfmt.Printf("{{Buildpack:}}::lightWhite %v, %v \n", pipeline.Buildpack.Name, pipeline.Buildpack.Language)
	if pipeline.Dockerimage != "" {
		fmt.Printf("{{Docker Image:}}::lightWhite %v \n", pipeline.Dockerimage)
	}
	cfmt.Printf("{{Deployment Strategy:}}::lightWhite %v \n", pipeline.Deploymentstrategy)
	cfmt.Printf("{{Git:}}::lightWhite %v:%v \n", pipeline.Git.Repository.SSHURL, pipeline.Git.Repository.DefaultBranch)
	cfmt.Printf("{{Review Apps:}}::lightWhite %v \n", pipeline.Reviewapps)
	cfmt.Printf("{{Phases:}}::lightWhite \n")
	for _, phase := range pipeline.Phases {
		if phase.Enabled {
			fmt.Printf(" - %v (%v) \n", phase.Name, phase.Context)
		}
	}
}
