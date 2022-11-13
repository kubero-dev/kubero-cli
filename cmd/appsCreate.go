/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new app",
	Long:  `Create a new app in a Pipeline`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")

		appPipeline := appsForm()
		fmt.Println(appPipeline)
	},
}

func init() {
	appsCmd.AddCommand(createCmd)
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

func appsForm() CreateApp {

	var ca CreateApp

	ca.Spec.Pipeline = promptLine("Pipeline", "", "")

	ca.Spec.Name = promptLine("Name", "", "")

	ca.Spec.Domain = promptLine("Domain", "", "")

	ca.Spec.Gitrepo.SSHURL = promptLine("Git SSH URL", "", "")

	ca.Spec.Branch = promptLine("Branch", "main", "main")

	autodeploy := promptLine("Audtodeploy", "Y,n", "Y")
	if autodeploy == "Y" {
		ca.Spec.Autodeploy = true
	} else {
		ca.Spec.Autodeploy = false
	}

	envCount, _ := strconv.Atoi(promptLine("Env Vars", "number", "0"))
	for i := 0; i < envCount; i++ {
		ca.Spec.EnvVars = append(ca.Spec.EnvVars, promptLine("Env Var", "", ""))
	}

	ca.Spec.Image.ContainerPort, _ = strconv.Atoi(promptLine("Container Port", "8080", "8080"))

	ca.Spec.Web.ReplicaCount, _ = strconv.Atoi(promptLine("Web Pods", "1", "1"))

	ca.Spec.Worker.ReplicaCount, _ = strconv.Atoi(promptLine("Worker Pods", "0", "0"))

	return ca
}
