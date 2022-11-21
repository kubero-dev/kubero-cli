package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

// pipelinesCmd represents the pipelines command
var pipelinesCmd = &cobra.Command{
	Use:   "pipelines",
	Short: "Manage your pipelines",
	Long: `List your pipelines
An App runs allways in a Pipeline. A Pipeline is a collection of Apps.`,
}

type Pipeline struct {
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
			Priv      string    `json:"priv"`
			Pub       string    `json:"pub"`
			ReadOnly  bool      `json:"read_only"`
			Title     string    `json:"title"`
			URL       string    `json:"url"`
			Verified  bool      `json:"verified"`
		} `json:"keys"`
		Repository struct {
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
	Name   string `json:"name"`
	Phases []struct {
		Context string `json:"context"`
		Enabled bool   `json:"enabled"`
		Name    string `json:"name"`
		Apps    []App  `json:"apps"`
	} `json:"phases"`
	Reviewapps bool `json:"reviewapps"`
}
type PipelinesList struct {
	Items []Pipeline `json:"items"`
}

var pipeline string

func init() {
	rootCmd.AddCommand(pipelinesCmd)
	pipelinesCmd.PersistentFlags().StringVarP(&pipeline, "pipeline", "p", "", "name of the pipeline")
}
