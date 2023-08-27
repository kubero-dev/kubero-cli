package kuberoCli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// appCmd represents the app command
var createAppCmd = &cobra.Command{
	Use:   "app",
	Short: "Create a new app in a Pipeline",
	Long: `Create a new app in a Pipeline.

If called without arguments, it will ask for all the required information`,
	Run: func(cmd *cobra.Command, args []string) {

		createApp := appForm(appName, pipelineName)
		writeAppYaml(createApp)

	},
}

func init() {
	createCmd.AddCommand(createAppCmd)
	createAppCmd.Flags().StringVarP(&pipelineName, "pipeline", "p", "", "Name of the pipeline")
	createAppCmd.Flags().StringVarP(&stageName, "stage", "s", "", "Name of the stage")
	createAppCmd.Flags().StringVarP(&appName, "app", "a", "", "Name of the app")
}

func appForm(AppName string, pipelineName string) AppCRD {

	var app AppCRD

	app.APIVersion = "application.kubero.dev/v1alpha1"
	app.Kind = "KuberoApp"

	if appName == "" {
		app.Spec.Name = promptLine("App Name", "", appName)
	} else {
		app.Spec.Name = appName
	}

	if pipelineName == "" {
		app.Spec.Pipeline = promptLine("Pipeline Name", "", pipelineName)
	} else {
		app.Spec.Pipeline = pipelineName
	}

	pipelineConfig := getPipelineConfig(pipelineName)
	availablePhases := getPipelinePhases(pipelineConfig)
	if stageName == "" {
		app.Spec.Phase = promptLine("Phase", fmt.Sprint(availablePhases), stageName)
	} else {
		app.Spec.Phase = stageName
	}

	app.Spec.Domain = promptLine("Domain", "", "")

	gitURL := pipelineConfig.GetString("spec.git.repository.sshurl")
	//ca.Spec.Gitrepo.SSHURL = promptLine("Git SSH URL", "["+getGitRemote()+"]", gitURL)

	//ca.Spec.Gitrepo.SSHURL = pipelineConfig.GetString("spec.git.repository")
	pipelineConfig.UnmarshalKey("spec.git.repository", &app.Spec.Gitrepo)
	app.Spec.Branch = promptLine("Branch", gitURL+":", "")

	app.Spec.Buildpack = pipelineConfig.GetString("spec.buildpack.name")

	autodeploy := promptLine("Autodeploy", "[y,n]", "n")
	if autodeploy == "Y" {
		app.Spec.Autodeploy = true
	} else {
		app.Spec.Autodeploy = false
	}

	envCount, _ := strconv.Atoi(promptLine("Number of Env Vars", "", "0"))
	for i := 0; i < envCount; i++ {
		app.Spec.EnvVars = append(app.Spec.EnvVars, promptLine("Env Var", "", ""))
	}

	app.Spec.Image.ContainerPort, _ = strconv.Atoi(promptLine("Container Port", "8080", ""))

	app.Spec.Web.ReplicaCount, _ = strconv.Atoi(promptLine("Web Pods", "1", ""))

	app.Spec.Worker.ReplicaCount, _ = strconv.Atoi(promptLine("Worker Pods", "0", ""))

	return app
}

type AppCRD struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
	} `json:"metadata"`
	Spec struct {
		Addons   []interface{} `json:"addons"`
		Affinity struct {
		} `json:"affinity"`
		Autodeploy  bool `json:"autodeploy"`
		Autoscale   bool `json:"autoscale"`
		Autoscaling struct {
			Enabled bool `json:"enabled"`
		} `json:"autoscaling"`
		Branch           string        `json:"branch"`
		Buildpack        string        `json:"buildpack"`
		Cronjobs         []interface{} `json:"cronjobs"`
		Domain           string        `json:"domain"`
		EnvVars          []interface{} `json:"envvars"`
		FullnameOverride string        `json:"fullnameOverride"`
		Gitrepo          struct {
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
		} `json:"gitrepo"`
		Image struct {
			Fetch struct {
				Repository string `json:"repository"`
				Tag        string `json:"tag"`
			} `json:"fetch"`
			Build struct {
				Command    string `json:"command"`
				Repository string `json:"repository"`
				Tag        string `json:"tag"`
			} `json:"build"`
			Run struct {
				Command    string `json:"command"`
				Repository string `json:"repository"`
				Tag        string `json:"tag"`
			} `json:"run"`
			ContainerPort int    `json:"containerPort"`
			PullPolicy    string `json:"pullPolicy"`
			Repository    string `json:"repository"`
			Tag           string `json:"tag"`
		} `json:"image"`
		ImagePullSecrets []interface{} `json:"imagePullSecrets"`
		Ingress          struct {
			Annotations struct {
			} `json:"annotations"`
			ClassName string `json:"className"`
			Enabled   bool   `json:"enabled"`
			Hosts     []struct {
				Host  string `json:"host"`
				Paths []struct {
					Path     string `json:"path"`
					PathType string `json:"pathType"`
				} `json:"paths"`
			} `json:"hosts"`
			TLS []interface{} `json:"tls"`
		} `json:"ingress"`
		Name         string `json:"appname"`
		NameOverride string `json:"nameOverride"`
		NodeSelector struct {
		} `json:"nodeSelector"`
		Phase          string `json:"phase"`
		Pipeline       string `json:"pipeline"`
		PodAnnotations struct {
		} `json:"podAnnotations"`
		PodSecurityContext struct {
		} `json:"podSecurityContext"`
		Podsize      string `json:"podsize"`
		ReplicaCount int    `json:"replicaCount"`
		Service      struct {
			Port int    `json:"port"`
			Type string `json:"type"`
		} `json:"service"`
		ServiceAccount struct {
			Annotations struct {
			} `json:"annotations"`
			Create bool   `json:"create"`
			Name   string `json:"name"`
		} `json:"serviceAccount"`
		Tolerations []interface{} `json:"tolerations"`
		Web         struct {
			Autoscaling struct {
				MaxReplicas                       int `json:"maxReplicas"`
				MinReplicas                       int `json:"minReplicas"`
				TargetCPUUtilizationPercentage    int `json:"targetCPUUtilizationPercentage"`
				TargetMemoryUtilizationPercentage int `json:"targetMemoryUtilizationPercentage"`
			} `json:"autoscaling"`
			ReplicaCount int `json:"replicaCount"`
		} `json:"web"`
		Worker struct {
			Autoscaling struct {
				MaxReplicas                       int `json:"maxReplicas"`
				MinReplicas                       int `json:"minReplicas"`
				TargetCPUUtilizationPercentage    int `json:"targetCPUUtilizationPercentage"`
				TargetMemoryUtilizationPercentage int `json:"targetMemoryUtilizationPercentage"`
			} `json:"autoscaling"`
			ReplicaCount int `json:"replicaCount"`
		} `json:"worker"`
	} `json:"spec"`
}

func writeAppYaml(app AppCRD) {
	// write pipeline.yaml
	yamlData, err := yaml.Marshal(&app)

	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}
	//fmt.Println(string(yamlData))

	fileName := ".kubero/" + app.Spec.Pipeline + "/" + app.Spec.Phase + "/" + app.Spec.Name + ".yaml"
	fmt.Println(fileName)
	err = os.WriteFile(fileName, yamlData, 0644)
	if err != nil {
		panic("Unable to write data into the file")
	}
}

func getPipelinePhases(pipelineConfig *viper.Viper) []string {
	var phases []string

	//pipelineConfig := getPipelineConfig(pipelineName)

	phasesList := pipelineConfig.GetStringSlice("spec.phases")

	for p := range phasesList {
		enabled := pipelineConfig.GetBool("spec.phases." + strconv.Itoa(p) + ".enabled")
		if enabled {
			phases = append(phases, pipelineConfig.GetString("spec.phases."+strconv.Itoa(p)+".name"))
		}
	}
	return phases
}

func getPipelineConfig(pipelineName string) *viper.Viper {

	basePath := ".kubero/"
	gitdir := getGitdir()
	dir := gitdir + basePath + pipelineName
	fmt.Println(dir)

	pipelineConfig := viper.New()
	pipelineConfig.SetConfigName("pipeline") // name of config file (without extension)
	pipelineConfig.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name
	pipelineConfig.AddConfigPath(dir)        // path to look for the config file in
	pipelineConfig.ReadInConfig()

	return pipelineConfig
}
