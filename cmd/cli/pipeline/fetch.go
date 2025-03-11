package pipeline

import (
	p "github.com/faelmori/kubero-cli/internal/pipeline"
	"github.com/spf13/cobra"
)

func FetchPipelineCmds() []*cobra.Command {
	return []*cobra.Command{
		cmdFetchPL(),
		cmdFetchApp(),
		cmdFetchPipelinePL(),
	}
}

func cmdFetchPL() *cobra.Command {
	var pipelineName, appName, stageName string

	var fetchCmd = &cobra.Command{
		Use:     "fetch",
		Aliases: []string{"pull", "fe"},
		Short:   "Fetch your remote pipelines and apps to your local repository",
		Long:    `Use the pipeline or app subcommand to fetch your pipelines and apps to your local repository`,
		RunE: func(cmd *cobra.Command, args []string) error {
			pl := p.NewPipelineManager(pipelineName, appName, stageName)
			if pipelineName != "" && appName == "" {
				err := pl.FetchPipeline(pipelineName)
				if err != nil {
					return err
				}
			} else if pipelineName != "" && appName != "" {
				if err := pl.EnsureStageNameIsSet(); err != nil {
					return err
				}
				fPlErr := pl.FetchPipeline(pipelineName)
				if fPlErr != nil {
					return fPlErr
				}
				fAppErr := pl.FetchApp(appName, stageName, pipelineName)
				if fAppErr != nil {
					return fAppErr
				}
			} else {
				pl.FetchAllPipelines()
			}
			return nil
		},
	}

	fetchCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "Name of the pipeline")
	fetchCmd.Flags().StringVarP(&stageName, "stage", "s", "", "Name of the stage [test|stage|production]")
	fetchCmd.Flags().StringVarP(&appName, "app", "a", "", "Name of the app")

	return fetchCmd
}

func cmdFetchApp() *cobra.Command {
	var appName, stageName, pipelineName string

	fetchAppCmd := &cobra.Command{
		Use:   "app",
		Short: "Fetch an app",
		Long:  `Fetch an app`,
		RunE: func(cmd *cobra.Command, args []string) error {
			plObj := p.NewPipelineManager(pipelineName, appName, stageName)
			return plObj.FetchApp(appName, stageName, pipelineName)
		},
	}

	fetchAppCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "Name of the pipeline")
	fetchAppCmd.Flags().StringVarP(&stageName, "stage", "s", "", "Name of the stage [test|stage|production]")
	fetchAppCmd.Flags().StringVarP(&appName, "app", "a", "", "Name of the app")

	return fetchAppCmd
}

func cmdFetchPipelinePL() *cobra.Command {
	var pipelineName string

	fetchPipelineCmd := &cobra.Command{
		Use:     "pipeline",
		Aliases: []string{"pl"},
		Short:   "Fetch a pipeline",
		Long:    `Fetch a pipeline`,
		RunE: func(cmd *cobra.Command, args []string) error {
			plObj := p.NewPipelineManager(pipelineName, "", "")
			pipelinesList := plObj.GetAllRemotePipelines()
			if err := plObj.EnsurePipelineIsSet(pipelinesList); err != nil {
				return err
			}
			return plObj.FetchPipeline(pipelineName)
		},
	}

	fetchPipelineCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "Name of the pipeline")

	return fetchPipelineCmd
}
