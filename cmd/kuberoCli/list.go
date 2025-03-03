package kuberoCli

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"fmt"
	"os"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List pipelines and apps",
	Long:    `List pipelines and apps`,
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
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
	listCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "output format [table, json]")
}
