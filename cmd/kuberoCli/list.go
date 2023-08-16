/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package kuberoCli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List pipelines and apps",
	Long:  `List pipelines and apps`,
	Run: func(cmd *cobra.Command, args []string) {

		if pipelineName != "" {
			// get a single pipeline
			pipelineResp, err := client.Get("/api/cli/pipelines/" + pipelineName)
			if pipelineResp.StatusCode() == 404 {
				cfmt.Println("{{  Pipeline not found}}::red")
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
			pipelineListResp, err := client.Get("/api/cli/pipelines")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			printPipelinesList(pipelineListResp)

		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
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
