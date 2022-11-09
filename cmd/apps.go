/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// appsCmd represents the apps command
var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "Manage your apps",
	Long: `Manage your apps

An App runs allways in a Pipeline.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("apps called")

		fmt.Println(getGitRemote())
	},
}

func init() {
	rootCmd.AddCommand(appsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	appsCmd.Flags().StringP("pipeline", "p", "", "Name of the pipeline")
	appsCmd.MarkFlagRequired("pipeline")
}

type App struct {
	Addons   []interface{} `json:"addons"`
	Affinity struct {
	} `json:"affinity"`
	Autodeploy  bool `json:"autodeploy"`
	Autoscale   bool `json:"autoscale"`
	Autoscaling struct {
		Enabled bool `json:"enabled"`
	} `json:"autoscaling"`
	Branch             string        `json:"branch"`
	Cronjobs           []interface{} `json:"cronjobs"`
	Deploymentstrategy string        `json:"deploymentstrategy"`
	Domain             string        `json:"domain"`
	EnvVars            []interface{} `json:"envVars"`
	FullnameOverride   string        `json:"fullnameOverride"`
	Gitrepo            struct {
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
		Build struct {
			Command    string `json:"command"`
			Repository string `json:"repository"`
			Tag        string `json:"tag"`
		} `json:"build"`
		ContainerPort int `json:"containerPort"`
		Fetch         struct {
			Repository string `json:"repository"`
			Tag        string `json:"tag"`
		} `json:"fetch"`
		PullPolicy string `json:"pullPolicy"`
		Repository string `json:"repository"`
		Run        struct {
			Command    string `json:"command"`
			Repository string `json:"repository"`
			Tag        string `json:"tag"`
		} `json:"run"`
		Tag string `json:"tag"`
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
	Name         string `json:"name"`
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
}
