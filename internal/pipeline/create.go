package pipeline

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	t "github.com/faelmori/kubero-cli/types"
	"github.com/i582/cfmt/cmd/cfmt"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
)

func (m *PipelineManager) CreatePipelineAndApp() error {
	createPipelineAndAppArg := promptLine("Create a new pipeline", "[y,n]", "y")
	if createPipelineAndAppArg == "y" {
		if _, createPlErr := m.CreatePipeline(); createPlErr != nil {
			return createPlErr
		}
	}

	pipelinesList := m.GetAllLocalApps()
	if ensurePlIsSetErr := m.EnsurePipelineIsSet(pipelinesList); ensurePlIsSetErr != nil {
		return ensurePlIsSetErr
	}
	if ensurePlNameIsSetErr := m.EnsureStageNameIsSet(); ensurePlNameIsSetErr != nil {
		m.EnsureAppNameIsSet()
	}
	if createAppErr := m.CreateApp(); createAppErr != nil {
		return createAppErr
	}

	return nil
}

func (m *PipelineManager) appForm() (*t.AppCRD, error) {
	var appCRD t.AppCRD

	plConfig := m.LoadAppConfig(m.pipelineName, m.stageName, m.appName)

	appCRD.APIVersion = "application.kubero.dev/v1alpha1"
	appCRD.Kind = "KuberoApp"

	appCRD.Spec.Name = m.appName
	appCRD.Spec.Pipeline = m.pipelineName
	appCRD.Spec.Phase = m.stageName

	appCRD.Spec.Domain = promptLine("Domain", "", plConfig.GetString("spec.domain"))

	unmarshalKeyErr := plConfig.UnmarshalKey("spec.git.repository", &appCRD.Spec.Gitrepo)
	if unmarshalKeyErr != nil {
		fmt.Println(unmarshalKeyErr)
		return nil, unmarshalKeyErr
	}

	gitURL := m.loadPipelineConfig(m.pipelineName).GetString("spec.git.repository.sshurl")
	appCRD.Spec.Branch = promptLine("Branch", gitURL+":", plConfig.GetString("spec.branch"))

	appCRD.Spec.Buildpack = m.loadPipelineConfig(m.pipelineName).GetString("spec.buildpack.name")

	autodeployDefault := "n"
	if !plConfig.GetBool("spec.autodeploy") {
		autodeployDefault = "y"
	}
	autodeploy := promptLine("Autodeploy", "[y,n]", autodeployDefault)
	if autodeploy == "Y" {
		appCRD.Spec.Autodeploy = true
	} else {
		appCRD.Spec.Autodeploy = false
	}

	envCount, _ := strconv.Atoi(promptLine("Number of Env Vars", "", "0"))
	appCRD.Spec.EnvVars = []string{}
	for i := 0; i < envCount; i++ {
		appCRD.Spec.EnvVars = append(appCRD.Spec.EnvVars, promptLine("Env Var", "", ""))
	}

	appCRD.Spec.Image.ContainerPort, _ = strconv.Atoi(promptLine("Container Port", "8080", plConfig.GetString("spec.image.containerport")))

	appCRD.Spec.Web = t.Web{}
	appCRD.Spec.Web.ReplicaCount, _ = strconv.Atoi(promptLine("Web Pods", "1", plConfig.GetString("spec.web.replicacount")))

	appCRD.Spec.Worker = t.Worker{}
	appCRD.Spec.Worker.ReplicaCount, _ = strconv.Atoi(promptLine("Worker Pods", "0", plConfig.GetString("spec.worker.replicacount")))

	return &appCRD, nil
}

func (m *PipelineManager) CreateApp() error {
	appCRD, appCRDErr := m.appForm()
	if appCRDErr != nil {
		return appCRDErr
	}

	if writeYamlAppErr := m.writeAppYaml(*appCRD); writeYamlAppErr != nil {
		return writeYamlAppErr
	}

	_, _ = cfmt.Println("\n\n{{Created appCRD.yaml}}::green")

	return nil
}

func (m *PipelineManager) writeAppYaml(appCRD t.AppCRD) error {
	yamlData, err := yaml.Marshal(&appCRD)

	if err != nil || appCRD.Spec.Name == "" {
		return err
	}

	fileName := ".kubero/" + appCRD.Spec.Pipeline + "/" + appCRD.Spec.Phase + "/" + appCRD.Spec.Name + ".yaml"

	err = os.WriteFile(fileName, yamlData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (m *PipelineManager) CreatePipeline() (*t.PipelineCRD, error) {
	_, plConfigsErr := m.LoadPLConfigs(m.pipelineName)
	if plConfigsErr != nil {
		return nil, plConfigsErr
	}

	m.LoadRepositories()
	m.LoadContexts()
	m.LoadBuildpacks()

	pipelineCRD := m.pipelinesForm()

	m.writePipelineYaml(pipelineCRD)

	_, _ = cfmt.Println("\n\n{{Created pipeline.yaml}}::green")
	_, _ = cfmt.Println(m.pipelineName)

	return &pipelineCRD, nil
}

func (m *PipelineManager) writePipelineYaml(pipeline t.PipelineCRD) {
	basePath := "/.kubero/" //TODO Make it dynamic

	gitdir := m.GetConfigManager().GetGitDir()
	dir := gitdir + basePath + m.pipelineName
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

func (m *PipelineManager) pipelinesForm() t.PipelineCRD {

	var pipelineCRD t.PipelineCRD

	if m.pipelineName == "" {
		m.pipelineName = promptLine("Define a PipelineName name", "", "")
	}
	pipelineCRD.Spec.Name = m.pipelineName

	pipelineCRD.APIVersion = "application.kubero.dev/v1alpha1"
	pipelineCRD.Kind = "KuberoPipeline"

	fmt.Println("")
	prompt := &survey.Select{
		Message: "Select a buildpack",
		Options: nil, //m.GetBuildPacksSimpleList(),
	}
	askOneErr := survey.AskOne(prompt, &pipelineCRD.Spec.Buildpack.Name)
	if askOneErr != nil {
		fmt.Println(askOneErr.Error())
		return t.PipelineCRD{}
	}

	domain := m.loadPipelineConfig(m.pipelineName).GetString("spec.domain")
	pipelineCRD.Spec.Domain = promptLine("FQDN Domain ", "", domain)

	// those fields are deprecated and may be removed in the future
	pipelineCRD.Spec.DockerImage = ""
	pipelineCRD.Spec.DeploymentStrategy = "git"

	gitconnection := promptLine("Connect pipeline to a Git repository (GitOps)", "[y,n]", "n")

	contextDefault := m.GetContextSimpleList()[0]
	if gitconnection == "y" {
		repoProviderList := t.Repositories{}
		gitProvider := m.loadPipelineConfig(m.pipelineName).GetString("spec.git.repository.provider")
		pipelineCRD.Spec.Git.Repository.Provider = promptLine("Repository Provider", fmt.Sprint(repoProviderList), gitProvider)

		gitURL := m.loadPipelineConfig(m.pipelineName).GetString("spec.git.repository.sshurl")
		girRemote := m.GetGitRemote()
		pipelineCRD.Spec.Git.Repository.SshUrl = promptLine("Repository URL", "["+girRemote+"]", gitURL)

		phaseReview := promptLine("enable reviewapps", "[y,n]", "n")
		if phaseReview == "y" {
			pipelineCRD.Spec.ReviewApps = true
			pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, t.Phase{
				Name:    "review",
				Enabled: true,
				Context: promptLine("Context for reviewapps", fmt.Sprint(m.GetContextSimpleList()), contextDefault),
			})
		} else {
			pipelineCRD.Spec.ReviewApps = false
			pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, t.Phase{
				Name:    "review",
				Enabled: false,
				Context: "",
			})
		}
	}

	phaseTest := promptLine("enable test", "[y,n]", "n")
	if phaseTest == "y" {
		pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, t.Phase{
			Name:    "test",
			Enabled: true,
			Context: promptLine("Context for test", fmt.Sprint(m.GetContextSimpleList()), contextDefault),
		})
	} else {
		pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, t.Phase{
			Name:    "test",
			Enabled: false,
			Context: "",
		})
	}

	phaseStage := promptLine("enable stage", "[y,n]", "n")
	if phaseStage == "y" {
		pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, t.Phase{
			Name:    "stage",
			Enabled: true,
			Context: promptLine("Context for stage", fmt.Sprint(m.GetContextSimpleList()), contextDefault),
		})
	} else {
		pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, t.Phase{
			Name:    "stage",
			Enabled: false,
			Context: "",
		})
	}

	phaseProduction := promptLine("enable production", "[y,n]", "y")
	if phaseProduction != "n" {
		pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, t.Phase{
			Name:    "production",
			Enabled: true,
			Context: promptLine("Context for production ", fmt.Sprint(m.GetContextSimpleList()), contextDefault),
		})
	} else {
		pipelineCRD.Spec.Phases = append(pipelineCRD.Spec.Phases, t.Phase{
			Name:    "production",
			Enabled: false,
			Context: "",
		})
	}

	return pipelineCRD
}

func (m *PipelineManager) GetContextSimpleList() []string {
	m.LoadContexts()
	if m.contexts == nil {
		return []string{}
	}
	contexts := m.contexts
	contextList := make([]string, 0)
	for _, context := range contexts {
		ctx := *context
		contextList = append(contextList, ctx.GetName())
	}

	return contextList
}
