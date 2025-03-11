package pipeline

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"fmt"
	"github.com/faelmori/kubero-cli/internal/pipeline"
	"github.com/spf13/cobra"
)

func CreateCmds() []*cobra.Command {
	return []*cobra.Command{
		cmdCreate(),
		cmdCreatePipeline(),
		cmdCreateApp(),
	}
}

func cmdCreate() *cobra.Command {
	var createCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"cr", "add", "new"},
		Short:   "Create a new pipeline and/or app",
		Long:    `Initiate a new pipeline and app in your current repository.`,
		Run: func(cmd *cobra.Command, args []string) {
			createPipelineAndApp()
		},
	}

	return createCmd
}

func cmdCreatePipeline() *cobra.Command {
	var createPipelineCmd = &cobra.Command{
		Use:     "pipeline",
		Aliases: []string{"pl"},
		Short:   "Create a new pipeline",
		Long:    `Create a new Pipeline`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("create a new pipeline")

			_ = createPipeline()
		},
	}

	return createPipelineCmd
}

func cmdCreateApp() *cobra.Command {
	var pipelineName string

	var createAppCmd = &cobra.Command{
		Use:   "app",
		Short: "Create a new app in a Pipeline",
		Long: `Create a new app in a Pipeline.

If called without arguments, it will ask for all the required information`,
		Run: func(cmd *cobra.Command, args []string) {
			c := pipeline.NewPipelineManager(pipelineName, stageName, appName)

			pipelinesList := c.getAllLocalPipelines()
			ensurePipelineIsSet(pipelinesList)
			ensureStageNameIsSet()
			ensureAppNameIsSet()
			createApp()

		},
	}

	createAppCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "Pipeline name")

	return createAppCmd
}
