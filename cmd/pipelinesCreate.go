/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// createCmd represents the create command
var PipelineCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new pipeline",
	Long:  `Create a new Pipeline`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create a new pipeline")

		loadRepositories()
		loadContexts()
		loadBuildpacks()
		createPipeline := pipelinesForm()

		client.SetBody(createPipeline.Spec)
		pipeline, pipelineErr := client.Post("/api/cli/pipelines/")

		if pipelineErr != nil {
			fmt.Println(pipelineErr)
		} else {
			cfmt.Println("{{Pipeline created successfully}}::green")
			json.Unmarshal(pipeline.Body(), &createPipeline.Spec)
			writePipelineYaml(createPipeline)
		}

	},
}

func init() {
	pipelinesCmd.AddCommand(PipelineCreateCmd)
}

type CreatePipeline struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Spec       struct {
		Buildpack struct {
			Build struct {
				Command    string `json:"command"`
				Repository string `json:"repository"`
				Tag        string `json:"tag"`
			} `json:"build"`
			Fetch struct {
				Repository string `json:"repository"`
				Tag        string `json:"tag"`
			} `json:"fetch"`
			Language string `json:"language"`
			Name     string `json:"name"`
			Run      struct {
				Command    string `json:"command"`
				Repository string `json:"repository"`
				Tag        string `json:"tag"`
			} `json:"run"`
		} `json:"buildpack"`
		Deploymentstrategy string `json:"deploymentstrategy"`
		Dockerimage        string `json:"dockerimage"`
		Git                struct {
			Keys struct {
				CreatedAt time.Time `json:"created_at"`
				ID        int       `json:"id"`
				//Priv      string    `json:"priv"`
				//Pub       string    `json:"pub"`
				ReadOnly bool   `json:"read_only"`
				Title    string `json:"title"`
				URL      string `json:"url"`
				Verified bool   `json:"verified"`
			} `json:"keys"`
			Repository struct {
				Provider      string `json:"provider"`
				Admin         bool   `json:"admin"`
				CloneURL      string `json:"clone_url"`
				DefaultBranch string `json:"default_branch"`
				Description   string `json:"description"`
				Homepage      string `json:"homepage"`
				ID            int    `json:"id"`
				Language      string `json:"language"`
				Name          string `json:"name"`
				NodeID        string `json:"node_id"`
				Owner         string `json:"owner"`
				Private       bool   `json:"private"`
				Push          bool   `json:"push"`
				SSHURL        string `json:"ssh_url"`
				Visibility    string `json:"visibility"`
			} `json:"repository"`
			Webhook struct {
				Active    bool      `json:"active"`
				CreatedAt time.Time `json:"created_at"`
				Events    []string  `json:"events"`
				ID        int       `json:"id"`
				Insecure  string    `json:"insecure"`
				URL       string    `json:"url"`
			} `json:"webhook"`
			Webhooks struct {
			} `json:"webhooks"`
		} `json:"git"`
		Name       string  `json:"pipelineName"`
		Phases     []Phase `json:"phases"`
		Reviewapps bool    `json:"reviewapps"`
	} `json:"spec"`
}

type CreatePipelines struct {
	PipelineName       string  `json:"pipelineName"`
	RepoProvider       string  `json:"repoprovider"`
	RepositoryURL      string  `json:"repositoryURL"`
	Phases             []Phase `json:"phases"`
	Reviewapps         bool    `json:"reviewapps"`
	Dockerimage        string  `json:"dockerimage"`
	Deploymentstrategy string  `json:"deploymentstrategy"`
	Buildpack          string  `json:"buildpack"`
}

type Phase struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Context string `json:"context"`
}

func writePipelineYaml(pipeline CreatePipeline) {
	// write pipeline.yaml
	yamlData, err := yaml.Marshal(&pipeline)

	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}
	//fmt.Println(string(yamlData))

	fileName := "pipeline.yaml"
	err = os.WriteFile(fileName, yamlData, 0644)
	if err != nil {
		panic("Unable to write data into the file")
	}
}

func pipelinesForm() CreatePipeline {

	var cp CreatePipeline

	cp.APIVersion = "application.kubero.dev/v1alpha1"
	cp.Kind = "KuberoPipeline"

	// those fields are deprecated and may be removed in the future
	cp.Spec.Dockerimage = ""
	cp.Spec.Deploymentstrategy = "git"

	if pipeline == "" {
		pipeline = pipelineConfig.GetString("spec.name")
		cp.Spec.Name = promptLine("Pipeline Name", "", pipeline)
	} else {
		cp.Spec.Name = pipeline
	}

	gitPrivider := pipelineConfig.GetString("spec.git.repository.provider")
	cp.Spec.Git.Repository.Provider = promptLine("Repository Provider", fmt.Sprint(repoSimpleList), gitPrivider)

	gitURL := pipelineConfig.GetString("spec.git.repository.sshurl")
	cp.Spec.Git.Repository.SSHURL = promptLine("Repository URL", "["+getGitRemote()+"]", gitURL)
	//cp.RepositoryURL = "git@github.com:kubero-dev/template-nodeapp.git"

	selectedBuildpack := pipelineConfig.GetString("spec.buildpack.name")
	cp.Spec.Buildpack.Name = promptLine("Buildpack ", fmt.Sprint(buildPacksSimpleList), selectedBuildpack)

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
