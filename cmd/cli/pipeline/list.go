package pipeline

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"github.com/faelmori/kubero-cli/internal/pipeline"
	"github.com/spf13/cobra"
)

func PipelineListCmds() []*cobra.Command {
	return []*cobra.Command{
		cmdPLList(),
	}
}

func cmdPLList() *cobra.Command {
	var pipelineName, outputFormat string
	var listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List pipelines and apps",
		Long:    `List pipelines and apps`,
		RunE: func(cmd *cobra.Command, args []string) error {
			pl := pipeline.NewPipelineManager(pipelineName, "", "")
			return pl.ListPipelines(pipelineName, outputFormat)
		},
	}

	listCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
	listCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "output format [table, json]")

	return listCmd
}
