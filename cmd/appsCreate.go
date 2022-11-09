/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new app",
	Long:  `Create a new app in a Pipeline`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")

		appPipeline := appsForm()
		fmt.Println(appPipeline)
	},
}

func init() {
	appsCmd.AddCommand(createCmd)
}

type CreateApp struct {
	PipelineName  string `json:"pipelineName"`
	RepoProvider  string `json:"repoprovider"`
	RepositoryURL string `json:"repositoryURL"`
	Phases        []struct {
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`
		Context string `json:"context"`
	} `json:"phases"`
	Reviewapps         bool   `json:"reviewapps"`
	Dockerimage        string `json:"dockerimage"`
	Deploymentstrategy string `json:"deploymentstrategy"`
	Buildpack          string `json:"buildpack"`
}

func appsForm() CreateApp {

	var ca CreateApp

	return ca
}
