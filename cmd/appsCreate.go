/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new app",
	Long:  `Create a new app in a Pipeline`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")

		createApp := appsForm()
		writeAppYaml(createApp)
		/*
			client.SetBody(createApp.Spec)
			_, appErr := client.Post("/api/cli/apps")

			if appErr != nil {
				fmt.Println(appErr)
			} else {
				cfmt.Println("{{App created successfully}}::green")
				//json.Unmarshal(app.Body(), &createApp.Spec)
				writeAppYaml(createApp)
			}
		*/
	},
}

func init() {
	appsCmd.AddCommand(createCmd)
	appsCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Skip asking for confirmation")
}

type CreateApp struct {
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

func writeAppYaml(app CreateApp) {
	// write pipeline.yaml
	yamlData, err := yaml.Marshal(&app)

	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}
	//fmt.Println(string(yamlData))

	fileName := "app." + app.Spec.Phase + ".yaml"
	err = ioutil.WriteFile(fileName, yamlData, 0644)
	if err != nil {
		panic("Unable to write data into the file")
	}
}

func appsForm() CreateApp {

	var ca CreateApp

	ca.APIVersion = "application.kubero.dev/v1alpha1"
	ca.Kind = "KuberoApp"

	ca.Spec.Pipeline = promptLine("Pipeline", "", pipelineConfig.GetString("spec.name"))

	availablePhases := getPipelinePases()
	ca.Spec.Phase = promptLine("Phase", fmt.Sprint(availablePhases), "")

	appconfig := loadAppConfig(ca.Spec.Phase)

	ca.Spec.Name = promptLine("Name", "", appconfig.GetString("spec.name"))

	ca.Spec.Domain = promptLine("Domain", "", appconfig.GetString("spec.domain"))

	gitURL := pipelineConfig.GetString("spec.git.repository.sshurl")
	ca.Spec.Gitrepo.SSHURL = promptLine("Git SSH URL", "["+getGitRemote()+"]", gitURL)

	ca.Spec.Branch = promptLine("Branch", "main", appconfig.GetString("spec.branch"))

	autodeployDefault := "n"
	if !appconfig.GetBool("spec.autodeploy") {
		autodeployDefault = "y"
	}
	autodeploy := promptLine("Autodeploy", "[y,n]", autodeployDefault)
	if autodeploy == "Y" {
		ca.Spec.Autodeploy = true
	} else {
		ca.Spec.Autodeploy = false
	}

	envCount, _ := strconv.Atoi(promptLine("Number of Env Vars", "", "0"))
	for i := 0; i < envCount; i++ {
		ca.Spec.EnvVars = append(ca.Spec.EnvVars, promptLine("Env Var", "", ""))
	}

	ca.Spec.Image.ContainerPort, _ = strconv.Atoi(promptLine("Container Port", "8080", appconfig.GetString("spec.image.containerport")))

	ca.Spec.Web.ReplicaCount, _ = strconv.Atoi(promptLine("Web Pods", "1", appconfig.GetString("spec.web.replicacount")))

	ca.Spec.Worker.ReplicaCount, _ = strconv.Atoi(promptLine("Worker Pods", "0", appconfig.GetString("spec.worker.replicacount")))

	return ca
}

func getPipelinePases() []string {
	var phases []string
	phasesList := pipelineConfig.GetStringSlice("spec.phases")

	for p := range phasesList {
		enabled := pipelineConfig.GetBool("spec.phases." + strconv.Itoa(p) + ".enabled")
		if enabled {
			phases = append(phases, pipelineConfig.GetString("spec.phases."+strconv.Itoa(p)+".name"))
		}
	}
	return phases
}

func loadAppConfig(phase string) *viper.Viper {

	appConfig := viper.New()
	appConfig.SetConfigName("app." + phase) // name of config file (without extension)
	appConfig.SetConfigType("yaml")         // REQUIRED if the config file does not have the extension in the name
	appConfig.AddConfigPath(".")            // path to look for the config file in
	appConfig.ReadInConfig()

	//fmt.Println("Using config file:", viper.ConfigFileUsed())
	//fmt.Println("Using config file:", pipelineConfig.ConfigFileUsed())

	return appConfig

}
