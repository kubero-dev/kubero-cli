/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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

		client.SetBody(createPipeline)
		pipeline, _ := client.Post("/api/cli/pipelines/")
		fmt.Println(pipeline)

	},
}

func init() {
	pipelinesCmd.AddCommand(PipelineCreateCmd)
}

type CreatePipeline struct {
	PipelineName  string `json:"pipelineName"`
	RepoProvider  string `json:"repoprovider"`
	RepositoryURL string `json:"repositoryURL"`
	Phases        []struct {
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`
		Context string `json:"context"`
	} `json:"phases"`
	Reviewapps         bool   `json:"reviewapps"`
	Dockerimage        string `json:"dockerimage"`
	Deploymentstrategy string `json:"deploymentstrategy"`
	Buildpack          string `json:"buildpack"`
}

type Phase struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Context string `json:"context"`
}

func pipelinesForm() CreatePipeline {

	var cp CreatePipeline

	// those fields are deprecated and may be removed in the future
	cp.Dockerimage = ""
	cp.Deploymentstrategy = "git"

	if pipeline == "" {
		cp.PipelineName = promptLine("Pipeline Name")
	} else {
		cp.PipelineName = pipeline
	}

	cp.RepoProvider = promptLine("Repository Provider " + fmt.Sprint(repoSimpleList))
	//cp.RepoProvider = "Github"

	cp.RepositoryURL = promptLine("Repository URL")
	//cp.RepositoryURL = "git@github.com:kubero-dev/template-nodeapp.git"

	cp.Buildpack = promptLine("Buildpack " + fmt.Sprint(buildPacksSimpleList))
	//cp.Buildpack = "NodeJS"

	phaseReview := promptLine("enable reviewapps [y,n]")
	if phaseReview == "y" {
		cp.Reviewapps = true
		cp.Phases = append(cp.Phases, Phase{
			Name:    "review",
			Enabled: true,
			Context: promptLine("Context for reviewapps " + fmt.Sprint(contextSimpleList)),
		})
	} else {
		cp.Reviewapps = false
		cp.Phases = append(cp.Phases, Phase{
			Name:    "review",
			Enabled: false,
			Context: "",
		})
	}

	phaseTest := promptLine("enable test [y,n]")
	if phaseTest == "y" {
		cp.Phases = append(cp.Phases, Phase{
			Name:    "test",
			Enabled: true,
			Context: promptLine("Context for test " + fmt.Sprint(contextSimpleList)),
		})
	} else {
		cp.Phases = append(cp.Phases, Phase{
			Name:    "test",
			Enabled: false,
			Context: "",
		})
	}

	phaseStage := promptLine("enable stage [y,n]")
	if phaseStage == "y" {
		cp.Phases = append(cp.Phases, Phase{
			Name:    "stage",
			Enabled: true,
			Context: promptLine("Context for stage " + fmt.Sprint(contextSimpleList)),
		})
	} else {
		cp.Phases = append(cp.Phases, Phase{
			Name:    "stage",
			Enabled: false,
			Context: "",
		})
	}

	phaseProduction := promptLine("enable production [y,n]")
	//var phaseProductionContext string = ""
	if phaseProduction != "n" {
		cp.Phases = append(cp.Phases, Phase{
			Name:    "production",
			Enabled: true,
			Context: promptLine("Context for production " + fmt.Sprint(contextSimpleList)),
		})
	} else {
		cp.Phases = append(cp.Phases, Phase{
			Name:    "production",
			Enabled: false,
			Context: "",
		})
	}

	return cp
}

//{"pipelineName":"aaaa","gitrepo":"git@github.com:kubero-dev/template-nodeapp.git","phases":[{"name":"review","enabled":false,"context":""},{"name":"test","enabled":false,"context":""},{"name":"stage","enabled":true,"context":"inClusterContext"},{"name":"production","enabled":true,"context":"inClusterContext"}],"reviewapps":true,"git":{"keys":{"id":73211402,"title":"bot@kubero","verified":true,"created_at":"2022-11-02T16:09:54Z","url":"https://api.github.com/repos/kubero-dev/template-nodeapp/keys/73211402","read_only":true,"pub":"c3NoLWVkMjU1MTkgQUFBQUMzTnphQzFsWkRJMU5URTVBQUFBSUdueHBQT0tXV3J2S0x6TGNoa2h6L3AreE54ZlhWbHlWektyaGFnUzJzeDIgKHVubmFtZWQp","priv":"LS0tLS1CRUdJTiBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0KYjNCbGJuTnphQzFyWlhrdGRqRUFBQUFBQkc1dmJtVUFBQUFFYm05dVpRQUFBQUFBQUFBQkFBQUFNd0FBQUF0emMyZ3RaVwpReU5UVXhPUUFBQUNCcDhhVHppbGxxN3lpOHkzSVpJYy82ZnNUY1gxMVpjbGN5cTRXb0V0ck1kZ0FBQUpES2xEZTV5cFEzCnVRQUFBQXR6YzJndFpXUXlOVFV4T1FBQUFDQnA4YVR6aWxscTd5aTh5M0laSWMvNmZzVGNYMTFaY2xjeXE0V29FdHJNZGcKQUFBRUJ2K0krTkp1alVwcmZ4QlBtdFBDWjhKak5teWtidnVWbXA3Ym84Wko3eU9HbnhwUE9LV1dydktMekxjaGtoei9wKwp4TnhmWFZseVZ6S3JoYWdTMnN4MkFBQUFDU2gxYm01aGJXVmtLUUVDQXdRPQotLS0tLUVORCBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0K"},"repository":{"id":501665730,"node_id":"R_kgDOHebPwg","name":"template-nodeapp","description":"Simple example Node app","owner":"kubero-dev","private":false,"ssh_url":"git@github.com:kubero-dev/template-nodeapp.git","clone_url":"https://github.com/kubero-dev/template-nodeapp.git","language":"JavaScript","homepage":"","admin":true,"push":true,"visibility":"public","default_branch":"main"},"webhooks":{},"webhook":{"id":386238670,"active":true,"created_at":"2022-10-30T20:26:50Z","url":"https://kubero6d304b3d55d72.loca.lt/api/repo/webhooks/github","insecure":"0","events":["pull_request","push"]}},"dockerimage":"","deploymentstrategy":"git","buildpack":{"name":"NodeJS","language":"JavaScript","fetch":{"repository":"ghcr.io/kubero-dev/buildpacks/fetch","tag":"main"},"build":{"repository":"node","tag":"latest","command":"npm install"},"run":{"repository":"node","tag":"latest","command":"node index.js"}}}
