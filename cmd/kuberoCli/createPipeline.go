/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package kuberoCli

import (
	"fmt"
	"os"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var pipelineName string

// pipelineCmd represents the pipeline command
var createPipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Create a new pipeline",
	Long:  `Create a new Pipeline`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create a new pipeline")

		_ = createPipeline()
	},
}

func init() {
	createCmd.AddCommand(createPipelineCmd)
}

func createPipeline() PipelineCRD {

	if pipelineName == "" {
		pipelineName = promptLine("Pipeline Name", "", "")
	}

	loadConfigs("/.kubero/", pipelineName)

	loadRepositories()
	loadContexts()
	loadBuildpacks()
	pipelineYaml := pipelinesForm()

	writePipelineYaml(pipelineYaml)

	cfmt.Println("\n\n{{Created pipeline.yaml}}::green")
	cfmt.Println(pipelineName)

	return pipelineYaml
}

func writePipelineYaml(pipeline PipelineCRD) {
	basePath := "/.kubero/" //TODO Make it dynamic

	gitdir := getGitdir()
	dir := gitdir + basePath + pipelineName
	err := os.MkdirAll(dir, 0755)

	if err != nil {
		fmt.Println(err)
		panic("Unable to create directory")
	}

	yamlData, err := yaml.Marshal(&pipeline)

	// iterate over phases to create the directory
	for _, phase := range pipeline.Spec.Phases {
		if phase.Enabled {
			err := os.MkdirAll(dir+"/"+phase.Name, 0755)
			if err != nil {
				fmt.Println(err)
				panic("Unable to create directory")
			}
		}
	}

	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}
	//fmt.Println(string(yamlData))

	fileName := dir + "/pipeline.yaml"
	err = os.WriteFile(fileName, yamlData, 0644)
	if err != nil {
		panic("Unable to write data into the file")
	}

}

func pipelinesForm() PipelineCRD {

	var cp PipelineCRD

	cp.APIVersion = "application.kubero.dev/v1alpha1"
	cp.Kind = "KuberoPipeline"

	selectedBuildpack := pipelineConfig.GetString("spec.buildpack.name")
	cp.Spec.Buildpack.Name = promptLine("Buildpack ", fmt.Sprint(buildPacksSimpleList), selectedBuildpack)

	domain := pipelineConfig.GetString("spec.domain")
	cp.Spec.Domain = promptLine("FQDN Domain ", fmt.Sprint(buildPacksSimpleList), domain)

	// those fields are deprecated and may be removed in the future
	cp.Spec.Dockerimage = ""
	cp.Spec.Deploymentstrategy = "git"

	gitconnection := promptLine("Connect pipeline to a Git repository (GitOps)", "[y,n]", "n")

	if gitconnection == "y" {
		gitPrivider := pipelineConfig.GetString("spec.git.repository.provider")
		cp.Spec.Git.Repository.Provider = promptLine("Repository Provider", fmt.Sprint(repoSimpleList), gitPrivider)

		gitURL := pipelineConfig.GetString("spec.git.repository.sshurl")
		cp.Spec.Git.Repository.SSHURL = promptLine("Repository URL", "["+getGitRemote()+"]", gitURL)
	}

	phaseReview := promptLine("enable reviewapps", "[y,n]", "n")
	if phaseReview == "y" {
		cp.Spec.Reviewapps = true
		contextDefault := pipelineConfig.GetString("spec.phases.0.context")
		cp.Spec.Phases = append(cp.Spec.Phases, Phase{
			Name:    "review",
			Enabled: true,
			Context: promptLine("Context for reviewapps", fmt.Sprint(contextSimpleList), contextDefault),
		})
	} else {
		cp.Spec.Reviewapps = false
		cp.Spec.Phases = append(cp.Spec.Phases, Phase{
			Name:    "review",
			Enabled: false,
			Context: "",
		})
	}

	phaseTest := promptLine("enable test", "[y,n]", "n")
	if phaseTest == "y" {
		contextDefault := pipelineConfig.GetString("spec.phases.1.context")
		cp.Spec.Phases = append(cp.Spec.Phases, Phase{
			Name:    "test",
			Enabled: true,
			Context: promptLine("Context for test", fmt.Sprint(contextSimpleList), contextDefault),
		})
	} else {
		cp.Spec.Phases = append(cp.Spec.Phases, Phase{
			Name:    "test",
			Enabled: false,
			Context: "",
		})
	}

	phaseStage := promptLine("enable stage", "[y,n]", "n")
	if phaseStage == "y" {
		contextDefault := pipelineConfig.GetString("spec.phases.2.context")
		cp.Spec.Phases = append(cp.Spec.Phases, Phase{
			Name:    "stage",
			Enabled: true,
			Context: promptLine("Context for stage", fmt.Sprint(contextSimpleList), contextDefault),
		})
	} else {
		cp.Spec.Phases = append(cp.Spec.Phases, Phase{
			Name:    "stage",
			Enabled: false,
			Context: "",
		})
	}

	phaseProduction := promptLine("enable production", "[y,n]", "y")
	//var phaseProductionContext string = ""
	if phaseProduction != "n" {
		contextDefault := pipelineConfig.GetString("spec.phases.3.context")
		cp.Spec.Phases = append(cp.Spec.Phases, Phase{
			Name:    "production",
			Enabled: true,
			Context: promptLine("Context for production ", fmt.Sprint(contextSimpleList), contextDefault),
		})
	} else {
		cp.Spec.Phases = append(cp.Spec.Phases, Phase{
			Name:    "production",
			Enabled: false,
			Context: "",
		})
	}

	return cp
}
