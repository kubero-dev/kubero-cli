package cli

import (
	"github.com/faelmori/kubero-cli/internal/pipeline"
	"github.com/spf13/cobra"
)

func FetchPipelineCmds() []*cobra.Command {
	return []*cobra.Command{
		cmdFetchPL(),
	}
}

func cmdFetchPL() *cobra.Command {
	var fetchCmd = &cobra.Command{
		Use:     "fetch",
		Aliases: []string{"pull", "fe"},
		Short:   "Fetch your remote pipelines and apps to your local repository",
		Long:    `Use the pipeline or app subcommand to fetch your pipelines and apps to your local repository`,
		Run: func(cmd *cobra.Command, args []string) {
			//pipelineManager := pipeline.NewPipelineManager()
			if pipelineName != "" && appName == "" {
				pipeline.FetchPipeline(pipelineName)
			} else if pipelineName != "" && appName != "" {
				ensureStageNameIsSet()
				pipeline.FetchPipeline(pipelineName)
				pipeline.FetchApp(appName, stageName, pipelineName)
			} else {
				pipeline.FetchAllPipelines()
			}
		},
	}

	fetchCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "Name of the pipeline")
	fetchCmd.Flags().StringVarP(&stageName, "stage", "s", "", "Name of the stage [test|stage|production]")
	fetchCmd.Flags().StringVarP(&appName, "app", "a", "", "Name of the app")

	return fetchCmd
}
