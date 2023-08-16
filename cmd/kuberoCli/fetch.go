package kuberoCli

import (
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Sync youre remote pipelines and apps to your local repository",
	Long:  `Use the pipeline or app subcommand to sync your pipelines and apps to your local repository`,
	Run: func(cmd *cobra.Command, args []string) {
		if pipelineName != "" && appName == "" {
			fetchPipeline()
		} else if appName != "" {
			fetchApp()
		} else {
			fetchAllPipelines()
		}
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
	fetchCmd.Flags().StringVarP(&appName, "app", "a", "", "name of the app")
}

func fetchPipeline() {
	confirmation := promptLine("Are you sure you want to fetch the pipeline "+pipelineName+"?", "[y,n]", "y")
	if confirmation == "y" {
		cfmt.Println("{{Fetching pipeline}} " + pipelineName + "::yellow")
	} else {
		cfmt.Println("{{Aborted}}::red")
		return
	}
}

func fetchApp() {

	if pipelineName == "" {
		cfmt.Println("{{Please specify a pipeline}}::red")
		return
	}

	confirmation := promptLine("Are you sure you want to fetch the app "+appName+" from "+pipelineName+"?", "[y,n]", "y")
	if confirmation == "y" {
		cfmt.Println("{{Fetching app}} " + appName + "::yellow")
	} else {
		cfmt.Println("{{Aborted}}::red")
		return
	}
}

func fetchAllPipelines() {
	confirmation := promptLine("Are you sure you want to fetch all pipelines?", "[y,n]", "n")
	if confirmation == "y" {
		cfmt.Println("{{Fetching all pipelines}}::yellow")
	} else {
		cfmt.Println("{{Aborted}}::red")
		return
	}
}
