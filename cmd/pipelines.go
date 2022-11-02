/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// pipelinesCmd represents the pipelines command
var pipelinesCmd = &cobra.Command{
	Use:   "pipelines",
	Short: "Manage your pipelines",
	Long: `List your pipelines
An App runs allways in a Pipeline. A Pipeline is a collection of Apps.`,
	Run: func(cmd *cobra.Command, args []string) {

		if pipeline != "" {
			// get a single pipeline
			pipeline, _ := client.Get("/api/cli/pipelines/" + pipeline)
			printPipeline(pipeline)
		} else {
			// get the pipelines
			resp, _ := client.Get("/api/cli/pipelines")
			printPipelinesList(resp)
		}
	},
}

type Pipeline struct {
	Buildpack struct {
		Build struct {
			Command    string `json:"command"`
			Repository string `json:"repository"`
			Tag        string `json:"tag"`
		} `json:"build"`
		Fetch struct {
			Repository string `json:"repository"`
			Tag        string `json:"tag"`
		} `json:"fetch"`
		Language string `json:"language"`
		Name     string `json:"name"`
		Run      struct {
			Command    string `json:"command"`
			Repository string `json:"repository"`
			Tag        string `json:"tag"`
		} `json:"run"`
	} `json:"buildpack"`
	Deploymentstrategy string `json:"deploymentstrategy"`
	Dockerimage        string `json:"dockerimage"`
	Git                struct {
		Keys struct {
			CreatedAt time.Time `json:"created_at"`
			ID        int       `json:"id"`
			Priv      string    `json:"priv"`
			Pub       string    `json:"pub"`
			ReadOnly  bool      `json:"read_only"`
			Title     string    `json:"title"`
			URL       string    `json:"url"`
			Verified  bool      `json:"verified"`
		} `json:"keys"`
		Repository struct {
			Admin         bool   `json:"admin"`
			CloneURL      string `json:"clone_url"`
			DefaultBranch string `json:"default_branch"`
			Description   string `json:"description"`
			Homepage      string `json:"homepage"`
			ID            int    `json:"id"`
			Language      string `json:"language"`
			Name          string `json:"name"`
			NodeID        string `json:"node_id"`
			Owner         string `json:"owner"`
			Private       bool   `json:"private"`
			Push          bool   `json:"push"`
			SSHURL        string `json:"ssh_url"`
			Visibility    string `json:"visibility"`
		} `json:"repository"`
		Webhook struct {
			Active    bool      `json:"active"`
			CreatedAt time.Time `json:"created_at"`
			Events    []string  `json:"events"`
			ID        int       `json:"id"`
			Insecure  string    `json:"insecure"`
			URL       string    `json:"url"`
		} `json:"webhook"`
		Webhooks struct {
		} `json:"webhooks"`
	} `json:"git"`
	Name   string `json:"name"`
	Phases []struct {
		Context string `json:"context"`
		Enabled bool   `json:"enabled"`
		Name    string `json:"name"`
	} `json:"phases"`
	Reviewapps bool `json:"reviewapps"`
}
type PipelinesList struct {
	Items []Pipeline `json:"items"`
}

var pipeline string

func init() {
	rootCmd.AddCommand(pipelinesCmd)
	pipelinesCmd.PersistentFlags().StringVarP(&pipeline, "pipeline", "p", "", "name of the pipeline")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pipelinesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pipelinesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// print the response as a table
func printPipelinesList(r *resty.Response) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Name",
		"Repository",
		"Branch",
		"Buildpack",
		"Docker Image",
		"Deployment Strategy",
		"Review Apps"})
	//table.SetBorder(false)

	var pipelinesList PipelinesList
	json.Unmarshal(r.Body(), &pipelinesList)

	for _, pipeline := range pipelinesList.Items {
		table.Append([]string{
			pipeline.Name,
			pipeline.Git.Repository.SSHURL,
			pipeline.Git.Repository.DefaultBranch,
			pipeline.Buildpack.Name,
			pipeline.Dockerimage,
			pipeline.Deploymentstrategy,
			fmt.Sprintf("%t", pipeline.Reviewapps)})
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
