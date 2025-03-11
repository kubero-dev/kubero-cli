package pipeline

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"fmt"
	"github.com/faelmori/kubero-cli/cmd/common"
	"github.com/faelmori/kubero-cli/internal/pipeline"
	"github.com/spf13/cobra"
)

func CreateCmds() []*cobra.Command {
	createRootCmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"cr", "add", "new"},
		Short:   "Create a new pipeline and/or app",
		Long:    `Initiate a new pipeline and app in your current repository.`,
		Annotations: common.GetDescriptions([]string{
			"Create a new pipeline and/or app",
			`Initiate a new pipeline and app in your current repository.`,
		}, false),
	}
	createRootCmd.AddCommand(cmdCreatePipeline())
	createRootCmd.AddCommand(cmdCreateApp())
	createRootCmd.AddCommand(cmdCreate())

	return []*cobra.Command{
		createRootCmd,
	}
}

func cmdCreate() *cobra.Command {
	var pipelineName, stageName, appName string

	var createCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"cr", "add", "new"},
		Short:   "Create a new pipeline and/or app",
		Long:    `Initiate a new pipeline and app in your current repository.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			pl := pipeline.NewPipelineManager(pipelineName, stageName, appName)
			return pl.CreatePipelineAndApp()
		},
	}

	createCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "Pipeline name")
	createCmd.Flags().StringVarP(&stageName, "stage", "s", "", "Stage name")
	createCmd.Flags().StringVarP(&appName, "app", "a", "", "App name")

	return createCmd
}

func cmdCreatePipeline() *cobra.Command {
	var pipelineName, stageName, appName string

	var createPipelineCmd = &cobra.Command{
		Use:     "pipeline",
		Aliases: []string{"pl"},
		Short:   "Create a new pipeline",
		Long:    `Create a new Pipeline`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("create a new pipeline")
			pl := pipeline.NewPipelineManager(pipelineName, stageName, appName)

			_, err := pl.CreatePipeline()
			if err != nil {
				return err
			}
			return nil
		},
	}

	return createPipelineCmd
}

func cmdCreateApp() *cobra.Command {
	var pipelineName, stageName, appName string

	var createAppCmd = &cobra.Command{
		Use:   "app",
		Short: "Create a new app in a Pipeline",
		Long: `Create a new app in a Pipeline.

If called without arguments, it will ask for all the required information`,
		RunE: func(cmd *cobra.Command, args []string) error {
			pl := pipeline.NewPipelineManager(pipelineName, stageName, appName)
			var err error
			pipelinesList := pl.GetAllLocalPipelines()
			if err = pl.EnsurePipelineIsSet(pipelinesList); err != nil {
				fmt.Println(err)
				return err
			}
			if err = pl.EnsureStageNameIsSet(); err != nil {
				fmt.Println(err)
				return err
			}
			if err = pl.EnsureAppNameIsSet(); err != nil {
				fmt.Println(err)
				return err
			}
			if err = pl.CreateApp(); err != nil {
				fmt.Println(err)
				return err
			}
			return nil
		},
	}

	createAppCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "Pipeline name")
	createAppCmd.Flags().StringVarP(&stageName, "stage", "s", "", "Stage name")
	createAppCmd.Flags().StringVarP(&appName, "app", "a", "", "App name")

	return createAppCmd
}
