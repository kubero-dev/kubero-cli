package kuberoApi

import "time"

type PipelineCRD struct {
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
		Domain             string `json:"domain"`
		Dockerimage        string `json:"dockerimage,omitempty"`
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
		} `json:"git,omitempty"`
		Name       string  `json:"pipelineName"`
		Phases     []Phase `json:"phases"`
		Reviewapps bool    `json:"reviewapps"`
	} `json:"spec"`
}

type Phase struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Context string `json:"context"`
}
