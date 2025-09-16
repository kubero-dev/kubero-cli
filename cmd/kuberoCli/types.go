package kuberoCli

import (
	"kubero/pkg/kuberoApi"
	"time"

	"gorm.io/gorm"
)

type Pipeline struct {
	gorm.Model
	BuildPack struct {
		Build struct {
			Command    string `json:"command" gorm:"column:command"`
			Repository string `json:"repository" gorm:"column:repository"`
			Tag        string `json:"tag" gorm:"column:tag"`
		} `json:"build" gorm:"embedded"`
		Fetch struct {
			Repository string `json:"repository" gorm:"column:repository"`
			Tag        string `json:"tag" gorm:"column:tag"`
		} `json:"fetch" gorm:"embedded"`
		Language string `json:"language" gorm:"column:language"`
		Name     string `json:"name" gorm:"column:name"`
		Run      struct {
			Command    string `json:"command" gorm:"column:command"`
			Repository string `json:"repository" gorm:"column:repository"`
			Tag        string `json:"tag" gorm:"column:tag"`
		} `json:"run" gorm:"embedded"`
	} `json:"buildpack" gorm:"embedded"`
	DeploymentStrategy string `json:"deploymentstrategy" gorm:"column:deploymentstrategy"`
	DockerImage        string `json:"dockerimage" gorm:"column:dockerimage"`
	Git                struct {
		Keys struct {
			ReadOnly bool   `json:"read_only" gorm:"column:read_only"`
			Title    string `json:"title" gorm:"column:title"`
			URL      string `json:"url" gorm:"column:url"`
			Verified bool   `json:"verified" gorm:"column:verified"`
		} `json:"keys" gorm:"embedded"`
		Repository struct {
			Admin         bool   `json:"admin" gorm:"column:admin"`
			CloneURL      string `json:"clone_url" gorm:"column:clone_url"`
			DefaultBranch string `json:"default_branch" gorm:"column:default_branch"`
			Description   string `json:"description" gorm:"column:description"`
			Homepage      string `json:"homepage" gorm:"column:homepage"`
			ID            int    `json:"id" gorm:"column:id"`
			Language      string `json:"language" gorm:"column:language"`
			Name          string `json:"name" gorm:"column:name"`
			NodeID        string `json:"node_id" gorm:"column:node_id"`
			Owner         string `json:"owner" gorm:"column:owner"`
			Private       bool   `json:"private" gorm:"column:private"`
			Push          bool   `json:"push" gorm:"column:push"`
			SshUrl        string `json:"ssh_url" gorm:"column:ssh_url"`
			Visibility    string `json:"visibility" gorm:"column:visibility"`
		} `json:"repository" gorm:"embedded"`
		Webhook struct {
			Active    bool      `json:"active" gorm:"column:active"`
			CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
			Events    []string  `json:"events" gorm:"column:events"`
			Insecure  string    `json:"insecure" gorm:"column:insecure"`
			URL       string    `json:"url" gorm:"column:url"`
		} `json:"webhook" gorm:"embedded"`
	} `json:"git" gorm:"embedded"`
	Name       string  `json:"name" gorm:"column:name"`
	Phases     []Phase `json:"phases" gorm:"foreignKey:PipelineID"`
	ReviewApps bool    `json:"reviewapps" gorm:"column:reviewapps"`
}

type Phase struct {
	gorm.Model
	PipelineID uint   `gorm:"column:pipeline_id"`
	Context    string `json:"context" gorm:"column:context"`
	Enabled    bool   `json:"enabled" gorm:"column:enabled"`
	Name       string `json:"name" gorm:"column:name"`
	Apps       []App  `json:"apps" gorm:"foreignKey:PhaseID"`
}

type PipelinesList struct {
	Items []Pipeline `json:"items"`
}

type Contexts []struct {
	Cluster string `json:"cluster"`
	Name    string `json:"name"`
	User    string `json:"user"`
}

type Repositories struct {
	Github    bool `json:"github"`
	Gitea     bool `json:"gitea"`
	Gitlab    bool `json:"gitlab"`
	Bitbucket bool `json:"bitbucket"`
	Docker    bool `json:"docker"`
}

type App struct {
	gorm.Model
	PhaseID     uint          `gorm:"column:phase_id"`
	Addons      []interface{} `json:"addons" gorm:"-"`
	Autodeploy  bool          `json:"autodeploy" gorm:"column:autodeploy"`
	Autoscale   bool          `json:"autoscale" gorm:"column:autoscale"`
	Autoscaling struct {
		Enabled bool `json:"enabled" gorm:"column:enabled"`
	} `json:"autoscaling" gorm:"embedded"`
	Branch             string        `json:"branch" gorm:"column:branch"`
	CronJobs           []interface{} `json:"cronjobs" gorm:"-"`
	DeploymentStrategy string        `json:"deploymentstrategy" gorm:"column:deploymentstrategy"`
	Domain             string        `json:"domain" gorm:"column:domain"`
	EnvVars            []interface{} `json:"envVars" gorm:"-"`
	FullnameOverride   string        `json:"fullnameOverride" gorm:"column:fullnameOverride"`
	Gitrepo            struct {
		Admin         bool   `json:"admin" gorm:"column:admin"`
		CloneURL      string `json:"clone_url" gorm:"column:clone_url"`
		DefaultBranch string `json:"default_branch" gorm:"column:default_branch"`
		Description   string `json:"description" gorm:"column:description"`
		Homepage      string `json:"homepage" gorm:"column:homepage"`
		ID            int    `json:"id" gorm:"column:id"`
		Language      string `json:"language" gorm:"column:language"`
		Name          string `json:"name" gorm:"column:name"`
		NodeID        string `json:"node_id" gorm:"column:node_id"`
		Owner         string `json:"owner" gorm:"column:owner"`
		Private       bool   `json:"private" gorm:"column:private"`
		Push          bool   `json:"push" gorm:"column:push"`
		SshUrl        string `json:"ssh_url" gorm:"column:ssh_url"`
		Visibility    string `json:"visibility" gorm:"column:visibility"`
	} `json:"gitrepo" gorm:"embedded"`
	Image struct {
		Build struct {
			Command    string `json:"command" gorm:"column:command"`
			Repository string `json:"repository" gorm:"column:repository"`
			Tag        string `json:"tag" gorm:"column:tag"`
		} `json:"build" gorm:"embedded"`
		ContainerPort string `json:"containerPort" gorm:"column:containerPort"`
		Fetch         struct {
			Repository string `json:"repository" gorm:"column:repository"`
			Tag        string `json:"tag" gorm:"column:tag"`
		} `json:"fetch" gorm:"embedded"`
		PullPolicy string `json:"pullPolicy" gorm:"column:pullPolicy"`
		Repository string `json:"repository" gorm:"column:repository"`
		Run        struct {
			Command    string `json:"command" gorm:"column:command"`
			Repository string `json:"repository" gorm:"column:repository"`
			Tag        string `json:"tag" gorm:"column:tag"`
		} `json:"run" gorm:"embedded"`
		Tag string `json:"tag" gorm:"column:tag"`
	} `json:"image" gorm:"embedded"`
	ImagePullSecrets []interface{} `json:"imagePullSecrets" gorm:"-"`
	Ingress          struct {
		Annotations struct {
		} `json:"annotations" gorm:"-"`
		ClassName string `json:"className" gorm:"column:className"`
		Enabled   bool   `json:"enabled" gorm:"column:enabled"`
		Hosts     []struct {
			Host  string `json:"host" gorm:"column:host"`
			Paths []struct {
				Path     string `json:"path" gorm:"column:path"`
				PathType string `json:"pathType" gorm:"column:pathType"`
			} `json:"paths" gorm:"-"`
		} `json:"hosts" gorm:"-"`
		TLS []interface{} `json:"tls" gorm:"-"`
	} `json:"ingress" gorm:"embedded"`
	Name         string `json:"name" gorm:"column:name"`
	NameOverride string `json:"nameOverride" gorm:"column:nameOverride"`
	NodeSelector struct {
	} `json:"nodeSelector" gorm:"-"`
	Phase          string `json:"phase" gorm:"column:phase"`
	Pipeline       string `json:"pipeline" gorm:"column:pipeline"`
	PodAnnotations struct {
	} `json:"podAnnotations" gorm:"-"`
	PodSecurityContext struct {
	} `json:"podSecurityContext" gorm:"-"`
	PodSize      PodSize `json:"podsize" gorm:"column:podsize"`
	ReplicaCount int     `json:"replicaCount" gorm:"column:replicaCount"`
	Service      struct {
		Port int    `json:"port" gorm:"column:port"`
		Type string `json:"type" gorm:"column:type"`
	} `json:"service" gorm:"embedded"`
	ServiceAccount struct {
		Annotations struct {
		} `json:"annotations" gorm:"-"`
		Create bool   `json:"create" gorm:"column:create"`
		Name   string `json:"name" gorm:"column:name"`
	} `json:"serviceAccount" gorm:"embedded"`
	Tolerations []interface{} `json:"tolerations" gorm:"-"`
	Web         Web           `json:"web" gorm:"embedded"`
	Worker      Worker        `json:"worker" gorm:"embedded"`
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

type Addon struct {
	gorm.Model
	ID      string `json:"id" gorm:"column:id"`
	Enabled bool   `json:"enabled" gorm:"column:enabled"`
	Version struct {
		Latest    string `json:"latest" gorm:"column:latest"`
		Installed string `json:"installed" gorm:"column:installed"`
	} `json:"version,omitempty" gorm:"embedded"`
	Description string `json:"description,omitempty" gorm:"column:description"`
	Readme      string `json:"readme,omitempty" gorm:"column:readme"`
	ArtifactURL string `json:"artifact_url" gorm:"column:artifact_url"`
	Kind        string `json:"kind" gorm:"column:kind"`
	Install     string `json:"install" gorm:"column:install"`
	Beta        bool   `json:"beta" gorm:"column:beta"`
}

type buildPacks []struct {
	Name     string `json:"name" gorm:"column:name"`
	Language string `json:"language" gorm:"column:language"`
	Fetch    struct {
		Repository string `json:"repository" gorm:"column:repository"`
		Tag        string `json:"tag" gorm:"column:tag"`
	} `json:"fetch" gorm:"embedded"`
	Build struct {
		Repository string `json:"repository" gorm:"column:repository"`
		Tag        string `json:"tag" gorm:"column:tag"`
		Command    string `json:"command" gorm:"column:command"`
	} `json:"build" gorm:"embedded"`
	Run struct {
		Repository         string `json:"repository" gorm:"column:repository"`
		Tag                string `json:"tag" gorm:"column:tag"`
		ReadOnlyAppStorage bool   `json:"readOnlyAppStorage" gorm:"column:readOnlyAppStorage"`
		SecurityContext    *struct {
			AllowPrivilegeEscalation *bool `json:"allowPrivilegeEscalation" gorm:"column:allowPrivilegeEscalation"`
			ReadOnlyRootFilesystem   *bool `json:"readOnlyRootFilesystem" gorm:"column:readOnlyRootFilesystem"`
		} `json:"securityContext" gorm:"embedded"`
		Command string `json:"command" gorm:"column:command"`
	} `json:"run,omitempty" gorm:"embedded"`
}

type PodSize struct {
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Default     bool   `json:"default,omitempty" gorm:"column:default"`
	Resources   struct {
		Requests struct {
			Memory string `json:"memory" gorm:"column:memory"`
			CPU    string `json:"cpu" gorm:"column:cpu"`
		} `json:"requests" gorm:"embedded"`
		Limits struct {
			Memory string `json:"memory" gorm:"column:memory"`
			CPU    string `json:"cpu" gorm:"column:cpu"`
		} `json:"limits,omitempty" gorm:"embedded"`
	} `json:"resources,omitempty" gorm:"embedded"`
	Active bool `json:"active,omitempty" gorm:"column:active"`
}

type pipelinesConfigsList map[string]kuberoApi.PipelineCRD

type appShort struct {
	Name     string `json:"name" gorm:"column:name"`
	Phase    string `json:"phase" gorm:"column:phase"`
	Pipeline string `json:"pipeline" gorm:"column:pipeline"`
}

type Instance struct {
	Name       string `json:"-" yaml:"-"`
	ApiUrl     string `json:"apiurl" yaml:"apiurl" gorm:"column:apiurl"`
	IacBaseDir string `json:"iacBaseDir,omitempty" yaml:"iacBaseDir,omitempty" gorm:"column:iacBaseDir"`
	ConfigPath string `json:"-" yaml:"-"`
	Tunnel     struct {
		Subdomain string `json:"subdomain" yaml:"subdomain" gorm:"column:subdomain"`
		Port      int    `json:"port" yaml:"port" gorm:"column:port"`
		Host      string `json:"host" yaml:"host" gorm:"column:host"`
	} `json:"tunnel,omitempty" yaml:"tunnel,omitempty"`
}

type Config struct {
	gorm.Model
	Api struct {
		Url   string `json:"url" yaml:"url" gorm:"column:url"`
		Token string `json:"token" yaml:"token" gorm:"column:token"`
	} `json:"api" yaml:"api" gorm:"embedded"`
}

type GithubVersion struct {
	gorm.Model
	Name       string `json:"name" gorm:"column:name"`
	ZipballUrl string `json:"zipball_url" gorm:"column:zipball_url"`
	TarballURL string `json:"tarball_url" gorm:"column:tarball_url"`
	Commit     struct {
		Sha string `json:"sha" gorm:"column:sha"`
		URL string `json:"url" gorm:"column:url"`
	} `json:"commit" gorm:"embedded"`
	NodeID string `json:"node_id" gorm:"column:node_id"`
}
