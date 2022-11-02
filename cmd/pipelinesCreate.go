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
		pipelinesForm()
	},
}

func init() {
	pipelinesCmd.AddCommand(PipelineCreateCmd)
}

type CreatePipeline struct {
	PipelineName string `json:"pipelineName"`
	Gitrepo      string `json:"gitrepo"`
	Phases       []struct {
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`
		Context string `json:"context"`
	} `json:"phases"`
	Reviewapps         bool   `json:"reviewapps"`
	RepositoryURL      string `json:"repositoryURL"`
	Dockerimage        string `json:"dockerimage"`
	Deploymentstrategy string `json:"deploymentstrategy"`
	Buildpack          string `json:"buildpack"`
}

func pipelinesForm() {
	if pipeline == "" {
		pipeline = promptLine("Pipeline Name")
	}

	repository := promptLine("Repository Type " + fmt.Sprint(repoSimpleList))
	repositoryURL := promptLine("Repository URL")
	Buildpack := promptLine("Buildpack " + fmt.Sprint(buildPacksSimpleList))

	phaseReview := promptLine("enable reviewapps [yes,no]")
	var phaseReviewContext string = ""
	if phaseReview == "yes" {
		phaseReviewContext = promptLine("Context for reviewapps " + fmt.Sprint(contextSimpleList))
	}

	phaseTest := promptLine("enable test [yes,no]")
	var phaseTestContext string = ""
	if phaseReview == "yes" {
		phaseTestContext = promptLine("Context for reviewapps " + fmt.Sprint(contextSimpleList))
	}

	phaseStage := promptLine("enable stage [yes,no]")
	var phaseStageContext string = ""
	if phaseReview == "yes" {
		phaseStageContext = promptLine("Context for reviewapps " + fmt.Sprint(contextSimpleList))
	}

	phaseProduction := promptLine("enable production [yes,no]")
	var phaseProductionContext string = ""
	if phaseReview == "yes" {
		phaseProductionContext = promptLine("Context for reviewapps " + fmt.Sprint(contextSimpleList))
	}

	fmt.Println(pipeline)
	fmt.Println(repository)
	fmt.Println(repositoryURL)
	fmt.Println(Buildpack)

	fmt.Println(phaseReview)
	fmt.Println(phaseReviewContext)
	fmt.Println(phaseTest)
	fmt.Println(phaseTestContext)
	fmt.Println(phaseStage)
	fmt.Println(phaseStageContext)
	fmt.Println(phaseProduction)
	fmt.Println(phaseProductionContext)

	CreatePipeline := CreatePipeline{
		PipelineName:       pipeline,
		Gitrepo:            repositoryURL,
		Reviewapps:         false,
		RepositoryURL:      repositoryURL,
		Dockerimage:        "",
		Deploymentstrategy: "",
		Buildpack:          Buildpack,
	}
	fmt.Println(CreatePipeline)
}

//{"pipelineName":"aaaa","gitrepo":"git@github.com:kubero-dev/template-nodeapp.git","phases":[{"name":"review","enabled":false,"context":""},{"name":"test","enabled":false,"context":""},{"name":"stage","enabled":true,"context":"inClusterContext"},{"name":"production","enabled":true,"context":"inClusterContext"}],"reviewapps":true,"git":{"keys":{"id":73211402,"title":"bot@kubero","verified":true,"created_at":"2022-11-02T16:09:54Z","url":"https://api.github.com/repos/kubero-dev/template-nodeapp/keys/73211402","read_only":true,"pub":"c3NoLWVkMjU1MTkgQUFBQUMzTnphQzFsWkRJMU5URTVBQUFBSUdueHBQT0tXV3J2S0x6TGNoa2h6L3AreE54ZlhWbHlWektyaGFnUzJzeDIgKHVubmFtZWQp","priv":"LS0tLS1CRUdJTiBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0KYjNCbGJuTnphQzFyWlhrdGRqRUFBQUFBQkc1dmJtVUFBQUFFYm05dVpRQUFBQUFBQUFBQkFBQUFNd0FBQUF0emMyZ3RaVwpReU5UVXhPUUFBQUNCcDhhVHppbGxxN3lpOHkzSVpJYy82ZnNUY1gxMVpjbGN5cTRXb0V0ck1kZ0FBQUpES2xEZTV5cFEzCnVRQUFBQXR6YzJndFpXUXlOVFV4T1FBQUFDQnA4YVR6aWxscTd5aTh5M0laSWMvNmZzVGNYMTFaY2xjeXE0V29FdHJNZGcKQUFBRUJ2K0krTkp1alVwcmZ4QlBtdFBDWjhKak5teWtidnVWbXA3Ym84Wko3eU9HbnhwUE9LV1dydktMekxjaGtoei9wKwp4TnhmWFZseVZ6S3JoYWdTMnN4MkFBQUFDU2gxYm01aGJXVmtLUUVDQXdRPQotLS0tLUVORCBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0K"},"repository":{"id":501665730,"node_id":"R_kgDOHebPwg","name":"template-nodeapp","description":"Simple example Node app","owner":"kubero-dev","private":false,"ssh_url":"git@github.com:kubero-dev/template-nodeapp.git","clone_url":"https://github.com/kubero-dev/template-nodeapp.git","language":"JavaScript","homepage":"","admin":true,"push":true,"visibility":"public","default_branch":"main"},"webhooks":{},"webhook":{"id":386238670,"active":true,"created_at":"2022-10-30T20:26:50Z","url":"https://kubero6d304b3d55d72.loca.lt/api/repo/webhooks/github","insecure":"0","events":["pull_request","push"]}},"dockerimage":"","deploymentstrategy":"git","buildpack":{"name":"NodeJS","language":"JavaScript","fetch":{"repository":"ghcr.io/kubero-dev/buildpacks/fetch","tag":"main"},"build":{"repository":"node","tag":"latest","command":"npm install"},"run":{"repository":"node","tag":"latest","command":"node index.js"}}}
