package kuberoApi

import (
	"time"

	"gorm.io/gorm"
)

type Metadata struct {
	gorm.Model
	Name   string      `json:"name,omitempty" yaml:"name,omitempty"`
	Labels interface{} `json:"labels,omitempty" yaml:"labels,omitempty"`
}

type GitKeys struct {
	gorm.Model
	ReadOnly bool   `json:"read_only"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Verified bool   `json:"verified"`
}

type GitRepository struct {
	gorm.Model
	Provider      string `json:"provider"`
	Admin         bool   `json:"admin"`
	CloneURL      string `json:"clone_url"`
	DefaultBranch string `json:"default_branch"`
	Description   string `json:"description"`
	Homepage      string `json:"homepage"`
	Language      string `json:"language"`
	Name          string `json:"name"`
	NodeID        string `json:"node_id" yaml:"node_id"`
	Owner         string `json:"owner"`
	Private       bool   `json:"private"`
	Push          bool   `json:"push"`
	SshUrl        string `json:"ssh_url"`
	Visibility    string `json:"visibility"`
}

type GitWebhook struct {
	gorm.Model
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	Events    []string  `json:"events"`
	Insecure  string    `json:"insecure"`
	URL       string    `json:"url"`
}

type Git struct {
	gorm.Model
	Keys       GitKeys       `json:"keys" gorm:"embedded"`
	Repository GitRepository `json:"repository" gorm:"embedded"`
	Webhook    GitWebhook    `json:"webhook" gorm:"embedded"`
}

type Build struct {
	gorm.Model
	Command    string `json:"command"`
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
}

type Fetch struct {
	gorm.Model
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
}

type Run struct {
	gorm.Model
	Command    string `json:"command"`
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
}

type Buildpack struct {
	gorm.Model
	Build    Build  `json:"build" gorm:"embedded"`
	Fetch    Fetch  `json:"fetch" gorm:"embedded"`
	Language string `json:"language"`
	Name     string `json:"name"`
	Run      Run    `json:"run" gorm:"embedded"`
}

type Phase struct {
	gorm.Model
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Context string `json:"context"`
}

type PipelineSpec struct {
	gorm.Model
	Buildpack          Buildpack `json:"buildpack" gorm:"embedded"`
	DeploymentStrategy string    `json:"deploymentstrategy"`
	Domain             string    `json:"domain"`
	DockerImage        string    `json:"dockerimage,omitempty" yaml:"dockerimage,omitempty"`
	Git                Git       `json:"git,omitempty" yaml:"git,omitempty" gorm:"embedded"`
	Name               string    `json:"pipelineName" yaml:"pipelineName"`
	Phases             []Phase   `json:"phases" gorm:"foreignKey:PipelineID"`
	ReviewApps         bool      `json:"reviewapps"`
}

type PipelineCRD struct {
	gorm.Model
	APIVersion string       `json:"apiVersion" yaml:"apiVersion"`
	Kind       string       `json:"kind"`
	Metadata   Metadata     `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Spec       PipelineSpec `json:"spec" gorm:"embedded"`
}

type AppSpec struct {
	gorm.Model
	Autodeploy  bool `json:"autodeploy"`
	Autoscale   bool `json:"autoscale"`
	Autoscaling struct {
		Enabled bool `json:"enabled"`
	} `json:"autoscaling" gorm:"embedded"`
	Branch           string        `json:"branch"`
	Buildpack        string        `json:"buildpack"`
	Domain           string        `json:"domain"`
	FullnameOverride string        `json:"fullnameOverride" yaml:"fullnameOverride"`
	Gitrepo          GitRepository `json:"gitrepo" gorm:"embedded"`
	Image            struct {
		Fetch struct {
			Repository string `json:"repository"`
			Tag        string `json:"tag"`
		} `json:"fetch" gorm:"embedded"`
		Build struct {
			Command    string `json:"command"`
			Repository string `json:"repository"`
			Tag        string `json:"tag"`
		} `json:"build" gorm:"embedded"`
		Run struct {
			Command    string `json:"command"`
			Repository string `json:"repository"`
			Tag        string `json:"tag"`
		} `json:"run" gorm:"embedded"`
		ContainerPort int    `json:"containerPort"`
		PullPolicy    string `json:"pullPolicy" yaml:"pullPolicy"`
		Repository    string `json:"repository"`
		Tag           string `json:"tag"`
	} `json:"image" gorm:"embedded"`
	Ingress struct {
		ClassName string `json:"className"`
		Enabled   bool   `json:"enabled"`
	} `json:"ingress" gorm:"embedded"`
	Name         string `json:"appname" yaml:"appname"`
	NameOverride string `json:"nameOverride" yaml:"nameOverride"`
	Phase        string `json:"phase"`
	Pipeline     string `json:"pipeline"`
	PodSize      string `json:"podsize"`
	ReplicaCount int    `json:"replicaCount" yaml:"replicaCount"`
	Service      struct {
		Port int    `json:"port"`
		Type string `json:"type"`
	} `json:"service" gorm:"embedded"`
	ServiceAccount struct {
		Create bool   `json:"create"`
		Name   string `json:"name"`
	} `json:"serviceAccount" yaml:"serviceAccount" gorm:"embedded"`
	EnvVars  []string `json:"envVars" gorm:"-"`
	Web      Web      `json:"web" gorm:"embedded"`
	Worker   Worker   `json:"worker" gorm:"embedded"`
	Security struct {
		VulnerabilityScans bool `json:"vulnerabilityScans,omitempty" yaml:"vulnerabilityScans,omitempty"`
	} `json:"security,omitempty" yaml:"security,omitempty"`
}

type Web struct {
	Autoscaling struct {
		MaxReplicas                       int `json:"maxReplicas" gorm:"column:maxReplicas"`
		MinReplicas                       int `json:"minReplicas" gorm:"column:minReplicas"`
		TargetCPUUtilizationPercentage    int `json:"targetCPUUtilizationPercentage" gorm:"column:targetCPUUtilizationPercentage"`
		TargetMemoryUtilizationPercentage int `json:"targetMemoryUtilizationPercentage" gorm:"column:targetMemoryUtilizationPercentage"`
	} `json:"autoscaling" gorm:"embedded"`
	ReplicaCount int `json:"replicaCount" gorm:"column:replicaCount"`
}

type Worker struct {
	Autoscaling struct {
		MaxReplicas                       int `json:"maxReplicas" gorm:"column:maxReplicas"`
		MinReplicas                       int `json:"minReplicas" gorm:"column:minReplicas"`
		TargetCPUUtilizationPercentage    int `json:"targetCPUUtilizationPercentage" gorm:"column:targetCPUUtilizationPercentage"`
		TargetMemoryUtilizationPercentage int `json:"targetMemoryUtilizationPercentage" gorm:"column:targetMemoryUtilizationPercentage"`
	} `json:"autoscaling" gorm:"embedded"`
	ReplicaCount int `json:"replicaCount" gorm:"column:replicaCount"`
}

type AppCRD struct {
	gorm.Model
	APIVersion string   `json:"apiVersion" yaml:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Spec       AppSpec  `json:"spec" gorm:"embedded"`
}
