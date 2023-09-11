package kuberoCli

import (
	"log"
	"os"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Undeploy your pipelines and apps from the cluster",
	Long: `Use the pipeline or app subcommand to undeploy your pipelines and apps from the cluster
Subcommands:
  kubero down [pipeline|app]`,
	Run: func(cmd *cobra.Command, args []string) {
		if pipelineName != "" && appName == "" {
			downPipeline()
		} else if appName != "" {
			downApp()
		} else {
			downAllPipelines()
		}
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
	downCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "name of the pipeline")
	//downCmd.MarkFlagRequired("pipeline")
	downCmd.Flags().StringVarP(&appName, "app", "a", "", "name of the app")
}

func downPipeline() {
	if pipelineName == "" {
		pipelineName = promptLine("Please define a pipeline ", "", "")
		return
	}
	downPipelineByName(pipelineName)
}

func downPipelineByName(pipelineName string) {
	confirmation := promptLine("Are you sure you want to undeploy the pipeline "+pipelineName+"?", "[y,n]", "y")
	if confirmation == "y" {
		cfmt.Println("{{Undeploying pipeline}}::yellow " + pipelineName)

		_, err := api.UnDeployPipeline(pipelineName)
		if err != nil {
			panic("Unable to undeploy Pipeline")
		}

	} else {
		cfmt.Println("{{Aborted}}::red")
		return
	}
}

func downApp() {

	if pipelineName == "" {
		cfmt.Println("{{Please specify a pipeline}}::red")
		return
	}

	if stageName == "" {
		cfmt.Println("{{Please specify a stage}}::red")
		return
	}

	confirmation := promptLine("Are you sure you want to undeploy the app "+appName+" from "+stageName+" in "+pipelineName+"?", "[y,n]", "y")
	if confirmation == "y" {
		cfmt.Println("{{Undeploying app}} " + appName + "::yellow")

		_, err := client.Delete("/api/cli/pipelines/" + pipelineName + "/" + stageName + "/" + appName)
		if err != nil {
			panic("Unable to undeploy App")
		}

	} else {
		cfmt.Println("{{Aborted}}::red")
		return
	}
}

func downAllPipelines() {
	confirmation := promptLine("Are you sure you want to undeploy all pipelines?", "[y,n]", "n")
	if confirmation == "y" {
		cfmt.Println("{{Undeploying all pipelines}}::yellow")
		pipelinesList := getAllLocalPipelines()
		for _, pipeline := range pipelinesList {
			downPipelineByName(pipeline)
		}

	} else {
		cfmt.Println("{{Aborted}}::red")
		return
	}
}

func getAllLocalPipelines() []string {

	basePath := "/.kubero/"
	gitdir := getGitdir()
	dir := gitdir + basePath + pipelineName

	pipelineNames := []string{}
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			if _, err := os.Stat(dir + "/" + f.Name() + "/pipeline.yaml"); err == nil {
				pipelineNames = append(pipelineNames, f.Name())
			}
		}
	}

	return pipelineNames
}
