package kuberoApi

import "time"

type PipelineCRD struct {
	APIVersion string   `json:"apiVersion" yaml:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   metadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`
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
		Dockerimage        string `json:"dockerimage,omitempty" yaml:"dockerimage,omitempty"`
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
				NodeID        string `json:"node_id" yaml:"node_id"`
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
			Webhooks struct { //TODO: This might be a typo
			} `json:"webhooks"`
		} `json:"git,omitempty" yaml:"git,omitempty"`
		Name       string  `json:"pipelineName" yaml:"pipelineName"`
		Phases     []Phase `json:"phases"`
		Reviewapps bool    `json:"reviewapps"`
	} `json:"spec"`
}

type Phase struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Context string `json:"context"`
}

type AppCRD struct {
	APIVersion string   `json:"apiVersion" yaml:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   metadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Spec       struct {
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
		FullnameOverride string        `json:"fullnameOverride" yaml:"fullnameOverride"`
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
			PullPolicy    string `json:"pullPolicy" yaml:"pullPolicy"`
			Repository    string `json:"repository"`
			Tag           string `json:"tag"`
		} `json:"image"`
		ImagePullSecrets []interface{} `json:"imagePullSecrets" yaml:"imagePullSecrets"`
		Ingress          struct {
			Annotations struct {
			} `json:"annotations"`
			ClassName string `json:"className"`
			Enabled   bool   `json:"enabled"`
			Hosts     []struct {
				Host  string `json:"host"`
				Paths []struct {
					Path     string `json:"path"`
					PathType string `json:"pathType" yaml:"pathType"`
				} `json:"paths"`
			} `json:"hosts"`
			TLS []interface{} `json:"tls"`
		} `json:"ingress"`
		Name         string `json:"appname" yaml:"appname"`
		NameOverride string `json:"nameOverride" yaml:"nameOverride"`
		NodeSelector struct {
		} `json:"nodeSelector" yaml:"NodeSelector"`
		Phase          string `json:"phase"`
		Pipeline       string `json:"pipeline"`
		PodAnnotations struct {
		} `json:"podAnnotations" yaml:"podAnnotations"`
		PodSecurityContext struct {
		} `json:"podSecurityContext" yaml:"podSecurityContext"`
		Podsize      string `json:"podsize"`
		ReplicaCount int    `json:"replicaCount" yaml:"replicaCount"`
		Service      struct {
			Port int    `json:"port"`
			Type string `json:"type"`
		} `json:"service"`
		ServiceAccount struct {
			Annotations struct {
			} `json:"annotations"`
			Create bool   `json:"create"`
			Name   string `json:"name"`
		} `json:"serviceAccount" yaml:"serviceAccount"`
		Tolerations []interface{} `json:"tolerations"`
		Web         struct {
			Autoscaling struct {
				MaxReplicas                       int `json:"maxReplicas" yaml:"maxReplicas"`
				MinReplicas                       int `json:"minReplicas" yaml:"minReplicas"`
				TargetCPUUtilizationPercentage    int `json:"targetCPUUtilizationPercentage" yaml:"targetCPUUtilizationPercentage"`
				TargetMemoryUtilizationPercentage int `json:"targetMemoryUtilizationPercentage" yaml:"targetMemoryUtilizationPercentage"`
			} `json:"autoscaling"`
			ReplicaCount int `json:"replicaCount yaml:pathType"`
		} `json:"web"`
		Worker struct {
			Autoscaling struct {
				MaxReplicas                       int `json:"maxReplicas" yaml:"maxReplicas"`
				MinReplicas                       int `json:"minReplicas" yaml:"minReplicas"`
				TargetCPUUtilizationPercentage    int `json:"targetCPUUtilizationPercentage" yaml:"targetCPUUtilizationPercentage"`
				TargetMemoryUtilizationPercentage int `json:"targetMemoryUtilizationPercentage" yaml:"targetMemoryUtilizationPercentage"`
			} `json:"autoscaling"`
			ReplicaCount int `json:"replicaCount yaml:pathType"`
		} `json:"worker"`
		Security struct {
			VulnerabilityScans bool `json:"vulnerabilityScans,omitempty" yaml:"vulnerabilityScans,omitempty"`
		} `json:"security,omitempty" yaml:"security,omitempty"`
	} `json:"spec"`
}

type metadata struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	//Namespace string      `json:"namespace,omitempty" yaml:"namespace,omitempty"` // we want to left this empty
	Labels interface{} `json:"labels,omitempty" yaml:"labels,omitempty"`
}
