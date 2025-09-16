package kuberoCli

import (
	"time"

	"gorm.io/gorm"
)

type DigitalOceanKubernetesConfig struct {
	Name      string `json:"name" gorm:"column:name"`
	Region    string `json:"region" gorm:"column:region"`
	Version   string `json:"version" gorm:"column:version"`
	NodePools []struct {
		Size  string `json:"size" gorm:"column:size"`
		Count int    `json:"count" gorm:"column:count"`
		Name  string `json:"name" gorm:"column:name"`
	} `json:"node_pools" gorm:"embedded"`
}

type DigitalOcean struct {
	KubernetesCluster struct {
		ID            string   `json:"id" gorm:"column:id"`
		Name          string   `json:"name" gorm:"column:name"`
		Region        string   `json:"region" gorm:"column:region"`
		Version       string   `json:"version" gorm:"column:version"`
		ClusterSubnet string   `json:"cluster_subnet" gorm:"column:cluster_subnet"`
		ServiceSubnet string   `json:"service_subnet" gorm:"column:service_subnet"`
		VpcUUID       string   `json:"vpc_uuid" gorm:"column:vpc_uuid"`
		Ipv4          string   `json:"ipv4" gorm:"column:ipv4"`
		Endpoint      string   `json:"endpoint" gorm:"column:endpoint"`
		Tags          []string `json:"tags" gorm:"column:tags"`
		NodePools     []struct {
			ID        string        `json:"id" gorm:"column:id"`
			Name      string        `json:"name" gorm:"column:name"`
			Size      string        `json:"size" gorm:"column:size"`
			Count     int           `json:"count" gorm:"column:count"`
			Tags      []string      `json:"tags" gorm:"column:tags"`
			Labels    interface{}   `json:"labels" gorm:"column:labels"`
			Taints    []interface{} `json:"taints" gorm:"column:taints"`
			AutoScale bool          `json:"auto_scale" gorm:"column:auto_scale"`
			MinNodes  int           `json:"min_nodes" gorm:"column:min_nodes"`
			MaxNodes  int           `json:"max_nodes" gorm:"column:max_nodes"`
			Nodes     []struct {
				ID     string `json:"id" gorm:"column:id"`
				Name   string `json:"name" gorm:"column:name"`
				Status struct {
					State string `json:"state" gorm:"column:state"`
				} `json:"status" gorm:"embedded"`
				DropletID string    `json:"droplet_id" gorm:"column:droplet_id"`
				CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
				UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
			} `json:"nodes" gorm:"embedded"`
		} `json:"node_pools" gorm:"embedded"`
		MaintenancePolicy struct {
			StartTime string `json:"start_time" gorm:"column:start_time"`
			Duration  string `json:"duration" gorm:"column:duration"`
			Day       string `json:"day" gorm:"column:day"`
		} `json:"maintenance_policy" gorm:"embedded"`
		AutoUpgrade bool `json:"auto_upgrade" gorm:"column:auto_upgrade"`
		Status      struct {
			State   string `json:"state" gorm:"column:state"`
			Message string `json:"message" gorm:"column:message"`
		} `json:"status" gorm:"embedded"`
		CreatedAt         time.Time `json:"created_at" gorm:"column:created_at"`
		UpdatedAt         time.Time `json:"updated_at" gorm:"column:updated_at"`
		SurgeUpgrade      bool      `json:"surge_upgrade" gorm:"column:surge_upgrade"`
		RegistryEnabled   bool      `json:"registry_enabled" gorm:"column:registry_enabled"`
		Ha                bool      `json:"ha" gorm:"column:ha"`
		SupportedFeatures []string  `json:"supported_features" gorm:"column:supported_features"`
	} `json:"kubernetes_cluster" gorm:"embedded"`
}

type DigitaloceanOptions struct {
	Options struct {
		Regions []struct {
			Name string `json:"name" gorm:"column:name"`
			Slug string `json:"slug" gorm:"column:slug"`
		} `json:"regions" gorm:"embedded"`
		Versions []struct {
			Slug              string   `json:"slug" gorm:"column:slug"`
			KubernetesVersion string   `json:"kubernetes_version" gorm:"column:kubernetes_version"`
			SupportedFeatures []string `json:"supported_features" gorm:"column:supported_features"`
		} `json:"versions" gorm:"embedded"`
		Sizes []struct {
			Name string `json:"name" gorm:"column:name"`
			Slug string `json:"slug" gorm:"column:slug"`
		} `json:"sizes" gorm:"embedded"`
	} `json:"options" gorm:"embedded"`
}

type User struct {
	gorm.Model
	Method   string `json:"method" gorm:"column:method"`
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
	Insecure bool   `json:"insecure" gorm:"column:insecure"`
	ApiToken string `json:"apiToken,omitempty" gorm:"column:apiToken"`
}

type KindConfig struct {
	Kind       string `yaml:"kind" gorm:"column:kind"`
	APIVersion string `yaml:"apiVersion" gorm:"column:apiVersion"`
	Name       string `yaml:"name" gorm:"column:name"`
	Networking struct {
		IPFamily         string `yaml:"ipFamily" gorm:"column:ipFamily"`
		APIServerAddress string `yaml:"apiServerAddress" gorm:"column:apiServerAddress"`
	} `yaml:"networking" gorm:"embedded"`
	Nodes []struct {
		Role                 string   `yaml:"role" gorm:"column:role"`
		Image                string   `yaml:"image,omitempty" gorm:"column:image"`
		KubeadmConfigPatches []string `yaml:"kubeadmConfigPatches" gorm:"column:kubeadmConfigPatches"`
		ExtraPortMappings    []struct {
			ContainerPort int    `yaml:"containerPort" gorm:"column:containerPort"`
			HostPort      int    `yaml:"hostPort" gorm:"column:hostPort"`
			Protocol      string `yaml:"protocol" gorm:"column:protocol"`
		} `yaml:"extraPortMappings" gorm:"embedded"`
	} `yaml:"nodes" gorm:"embedded"`
	ContainerdConfigPatches []string `yaml:"containerdConfigPatches" gorm:"column:containerdConfigPatches"`
}

type KuberoUIConfig struct {
	APIVersion string `yaml:"apiVersion" gorm:"column:apiVersion"`
	Kind       string `yaml:"kind" gorm:"column:kind"`
	Metadata   struct {
		Name string `yaml:"name" gorm:"column:name"`
	} `yaml:"metadata" gorm:"embedded"`
	Spec struct {
		FullnameOverride string `yaml:"fullnameOverride" gorm:"column:fullnameOverride"`
		Image            struct {
			PullPolicy string `yaml:"pullPolicy" gorm:"column:pullPolicy"`
			Repository string `yaml:"repository" gorm:"column:repository"`
			Tag        string `yaml:"tag" gorm:"column:tag"`
		} `yaml:"image" gorm:"embedded"`
		ImagePullSecrets []interface{} `yaml:"imagePullSecrets" gorm:"-"`
		Ingress          struct {
			Annotations struct {
				KubernetesIoIngressClass string `yaml:"cert-manager.io/cluster-issuer,omitempty" gorm:"column:cert_manager_io_cluster_issuer"`
				KubernetesIoTlsAcme      string `yaml:"kubernetes.io/tls-acme,omitempty" gorm:"column:kubernetes_io_tls_acme"`
			} `yaml:"annotations" gorm:"embedded"`
			ClassName string `yaml:"className" gorm:"column:className"`
			Enabled   bool   `yaml:"enabled" gorm:"column:enabled"`
			Hosts     []struct {
				Host  string `yaml:"host" gorm:"column:host"`
				Paths []struct {
					Path     string `yaml:"path" gorm:"column:path"`
					PathType string `yaml:"pathType" gorm:"column:pathType"`
				} `yaml:"paths" gorm:"embedded"`
			} `yaml:"hosts" gorm:"embedded"`
			TLS []KuberoUITls `yaml:"tls" gorm:"-"`
		} `yaml:"ingress" gorm:"embedded"`
		NameOverride string `yaml:"nameOverride" gorm:"column:nameOverride"`
		NodeSelector struct {
		} `yaml:"nodeSelector" gorm:"-"`
		PodAnnotations struct {
		} `yaml:"podAnnotations" gorm:"-"`
		PodSecurityContext struct {
		} `yaml:"podSecurityContext" gorm:"-"`
		Prometheus struct {
			Enabled  bool   `yaml:"enabled" gorm:"column:enabled"`
			Endpoint string `yaml:"endpoint" gorm:"column:endpoint"`
		} `yaml:"prometheus,omitempty" gorm:"embedded"`
		Registry struct {
			Enabled bool   `yaml:"enabled" gorm:"column:enabled"`
			Create  bool   `yaml:"create" gorm:"column:create"`
			Host    string `yaml:"host" gorm:"column:host"`
			SubPath string `yaml:"subPath" gorm:"column:subPath"`
			Account struct {
				Username string `yaml:"username" gorm:"column:username"`
				Password string `yaml:"password" gorm:"column:password"`
				Hash     string `yaml:"hash" gorm:"column:hash"`
			} `yaml:"account" gorm:"embedded"`
			Port             int         `yaml:"port" gorm:"column:port"`
			Storage          string      `yaml:"storage" gorm:"column:storage"`
			StorageClassName interface{} `yaml:"storageClassName" gorm:"column:storageClassName"`
		} `yaml:"registry" gorm:"embedded"`
		ReplicaCount int `yaml:"replicaCount" gorm:"column:replicaCount"`
		Resources    struct {
		} `yaml:"resources" gorm:"-"`
		SecurityContext struct {
		} `yaml:"securityContext" gorm:"-"`
		Service struct {
			Port int    `yaml:"port" gorm:"column:port"`
			Type string `yaml:"type" gorm:"column:type"`
		} `yaml:"service" gorm:"embedded"`
		ServiceAccount struct {
			Annotations struct {
			} `yaml:"annotations" gorm:"-"`
			Create bool   `yaml:"create" gorm:"column:create"`
			Name   string `yaml:"name" gorm:"column:name"`
		} `yaml:"serviceAccount" gorm:"embedded"`
		Tolerations []interface{} `yaml:"tolerations" gorm:"-"`
		Kubero      struct {
			Debug      string `yaml:"debug" gorm:"column:debug"`
			Namespace  string `yaml:"namespace" gorm:"column:namespace"`
			Context    string `yaml:"context" gorm:"column:context"`
			WebhookURL string `yaml:"webhook_url" gorm:"column:webhook_url"`
			SessionKey string `yaml:"sessionKey" gorm:"column:sessionKey"`
			Auth       struct {
				Github struct {
					Enabled     bool   `yaml:"enabled" gorm:"column:enabled"`
					ID          string `yaml:"id,omitempty" gorm:"column:id"`
					Secret      string `yaml:"secret,omitempty" gorm:"column:secret"`
					CallbackURL string `yaml:"callbackUrl,omitempty" gorm:"column:callbackUrl"`
					Org         string `yaml:"org,omitempty" gorm:"column:org"`
				} `yaml:"github" gorm:"embedded"`
				Oauth2 struct {
					Enabled     bool   `yaml:"enabled" gorm:"column:enabled"`
					Name        string `yaml:"name,omitempty" gorm:"column:name"`
					ID          string `yaml:"id,omitempty" gorm:"column:id"`
					AuthURL     string `yaml:"authUrl,omitempty" gorm:"column:authUrl"`
					TokenURL    string `yaml:"tokenUrl,omitempty" gorm:"column:tokenUrl"`
					Secret      string `yaml:"secret,omitempty" gorm:"column:secret"`
					CallbackURL string `yaml:"callbackUrl,omitempty" gorm:"column:callbackUrl"`
				} `yaml:"oauth2" gorm:"embedded"`
			} `yaml:"auth" gorm:"embedded"`
			AuditLogs struct {
				Enabled          bool     `yaml:"enabled" gorm:"column:enabled"`
				StorageClassName string   `yaml:"storageClassName" gorm:"column:storageClassName"`
				Size             string   `yaml:"size" gorm:"column:size"`
				AccessModes      []string `yaml:"accessModes" gorm:"column:accessModes"`
				Limit            string   `yaml:"limit" gorm:"column:limit"`
			} `yaml:"auditLogs" gorm:"embedded"`
			DataBase struct {
				StorageClassName string   `yaml:"storageClassName" gorm:"column:storageClassName"`
				Size             string   `yaml:"size" gorm:"column:size"`
				AccessModes      []string `yaml:"accessModes" gorm:"column:accessModes"`
			} `yaml:"database" gorm:"embedded"`
			Config KuberoConfigfile `yaml:"config,omitempty" gorm:"embedded"`
		} `yaml:"kubero" gorm:"embedded"`
	} `yaml:"spec" gorm:"embedded"`
}

type NotificationsConfig struct {
	Enabled   bool     `yaml:"enabled" gorm:"column:enabled"`
	Name      string   `yaml:"name" gorm:"column:name"`
	Type      string   `yaml:"type" gorm:"column:type"`
	Pipelines []string `yaml:"pipelines" gorm:"column:pipelines"`
	Events    []string `yaml:"events" gorm:"column:events"`
	Config    struct {
		Url     string `yaml:"url" gorm:"column:url"`
		Channel string `yaml:"channel,omitempty" gorm:"column:channel"`
		Secret  string `yaml:"secret,omitempty" gorm:"column:secret"`
	} `yaml:"config" gorm:"embedded"`
}

type KuberoConfigfile struct {
	Kubero struct {
		Readonly bool `yaml:"readonly" gorm:"column:readonly"`
		Admin    struct {
			Disabled bool `yaml:"disabled" gorm:"column:enabled"`
		} `yaml:"admin" gorm:"embedded"`
		Console struct {
			Enabled bool `yaml:"enabled" gorm:"column:enabled"`
		} `yaml:"console" gorm:"embedded"`
		Banner struct {
			Show      bool   `yaml:"show" gorm:"column:show"`
			Message   string `yaml:"message" gorm:"column:message"`
			BgColor   string `yaml:"bgColor" gorm:"column:bgColor"`
			Fontcolor string `yaml:"fontcolor" gorm:"column:fontcolor"`
		} `yaml:"banner" gorm:"embedded"`
	} `yaml:"kubero" gorm:"embedded"`
	Notifications []NotificationsConfig `yaml:"notifications" gorm:"embedded"`
	ClusterIssuer string                `yaml:"clusterIssuer" gorm:"column:clusterIssuer"`
	Templates     struct {
		Enabled  bool `yaml:"enabled" gorm:"column:enabled"`
		Catalogs []struct {
			Name             string `yaml:"name" gorm:"column:name"`
			Description      string `yaml:"description" gorm:"column:description"`
			TemplateBasePath string `yaml:"templateBasePath" gorm:"column:templateBasePath"`
			Index            struct {
				URL    string `yaml:"url" gorm:"column:url"`
				Format string `yaml:"format" gorm:"column:format"`
			} `yaml:"index" gorm:"embedded"`
		} `yaml:"catalogs" gorm:"embedded"`
	} `yaml:"templates" gorm:"embedded"`
	BuildPacks []struct {
		Name     string `yaml:"name" gorm:"column:name"`
		Language string `yaml:"language" gorm:"column:language"`
		Fetch    struct {
			Repository      string `yaml:"repository" gorm:"column:repository"`
			Tag             string `yaml:"tag" gorm:"column:tag"`
			SecurityContext struct {
				RunAsUser int `yaml:"runAsUser" gorm:"column:runAsUser"`
			} `yaml:"securityContext" gorm:"embedded"`
		} `yaml:"fetch" gorm:"embedded"`
		Build struct {
			Repository      string `yaml:"repository" gorm:"column:repository"`
			Tag             string `yaml:"tag" gorm:"column:tag"`
			Command         string `yaml:"command" gorm:"column:command"`
			SecurityContext struct {
				RunAsUser int `yaml:"runAsUser" gorm:"column:runAsUser"`
			} `yaml:"securityContext" gorm:"embedded"`
		} `yaml:"build,omitempty" gorm:"embedded"`
		Run struct {
			Repository         string `yaml:"repository" gorm:"column:repository"`
			Tag                string `yaml:"tag" gorm:"column:tag"`
			ReadOnlyAppStorage bool   `yaml:"readOnlyAppStorage" gorm:"column:readOnlyAppStorage"`
			SecurityContext    struct {
				AllowPrivilegeEscalation bool `yaml:"allowPrivilegeEscalation" gorm:"column:allowPrivilegeEscalation"`
				ReadOnlyRootFilesystem   bool `yaml:"readOnlyRootFilesystem" gorm:"column:readOnlyRootFilesystem"`
			} `yaml:"securityContext" gorm:"embedded"`
			Command string `yaml:"command" gorm:"column:command"`
		} `yaml:"run,omitempty" gorm:"embedded"`
	} `yaml:"buildpacks" gorm:"embedded"`
	PodSizeList []struct {
		Active      bool   `yaml:"active,omitempty" gorm:"column:active"`
		Name        string `yaml:"name" gorm:"column:name"`
		Description string `yaml:"description" gorm:"column:description"`
		Default     bool   `yaml:"default,omitempty" gorm:"column:default"`
		Resources   struct {
			Requests struct {
				Memory string `yaml:"memory" gorm:"column:memory"`
				CPU    string `yaml:"cpu" gorm:"column:cpu"`
			} `yaml:"requests" gorm:"embedded"`
			Limits struct {
				Memory string `yaml:"memory" gorm:"column:memory"`
				CPU    string `yaml:"cpu" gorm:"column:cpu"`
			} `yaml:"limits" gorm:"embedded"`
		} `yaml:"resources,omitempty" gorm:"embedded"`
	} `yaml:"podSizeList" gorm:"embedded"`
}

type KuberoUITls struct {
	SecretName string   `yaml:"secretName" gorm:"column:secretName"`
	Hosts      []string `yaml:"hosts" gorm:"column:hosts"`
}

// ScalewayCreate https://developers.scaleway.com/en/products/k8s/api/#post-612200
type ScalewayCreate struct {
	OrganizationID string             `json:"organization_id,omitempty" gorm:"column:organization_id"`
	ProjectID      string             `json:"project_id,omitempty" gorm:"column:project_id"`
	Type           string             `json:"type" gorm:"column:type"`
	Name           string             `json:"name" gorm:"column:name"`
	Description    string             `json:"description" gorm:"column:description"`
	Tags           []string           `json:"tags" gorm:"column:tags"`
	Version        string             `json:"version" gorm:"column:version"`
	Cni            string             `json:"cni" gorm:"column:cni"`
	Pools          []ScalewayNodePool `json:"pools" gorm:"embedded"`
	AutoUpgrade    struct {
		Enable            bool `json:"enable" gorm:"column:enable"`
		MaintenanceWindow struct {
			StartHour int    `json:"start_hour" gorm:"column:start_hour"`
			Day       string `json:"day" gorm:"column:day"`
		} `json:"maintenance_window" gorm:"embedded"`
	} `json:"auto_upgrade" gorm:"embedded"`
	FeatureGates      []string `json:"feature_gates" gorm:"column:feature_gates"`
	AdmissionPlugins  []string `json:"admission_plugins" gorm:"column:admission_plugins"`
	ApiServerCertSans []string `json:"apiServer_cert_sans" gorm:"column:apiServer_cert_sans"`
	Ingress           string   `json:"ingress" gorm:"column:ingress"`
}

type ScalewayNodePool struct {
	Name             string   `json:"name" gorm:"column:name"`
	NodeType         string   `json:"node_type" gorm:"column:node_type"`
	PlacementGroupID string   `json:"placement_group_id,omitempty" gorm:"column:placement_group_id"`
	Autoscaling      bool     `json:"autoscaling" gorm:"column:autoscaling"`
	Size             int      `json:"size" gorm:"column:size"`
	MinSize          int      `json:"min_size" gorm:"column:min_size"`
	MaxSize          int      `json:"max_size" gorm:"column:max_size"`
	ContainerRuntime string   `json:"container_runtime" gorm:"column:container_runtime"`
	AutoHealing      bool     `json:"autoHealing" gorm:"column:autoHealing"`
	Tags             []string `json:"tags" gorm:"column:tags"`
	Zone             string   `json:"zone" gorm:"column:zone"`
	RootVolumeType   string   `json:"root_volume_type" gorm:"column:root_volume_type"`
	RootVolumeSize   int      `json:"root_volume_size" gorm:"column:root_volume_size"`
}

type ScalewayVersionsResponse struct {
	Versions []struct {
		Name                       string   `json:"name" gorm:"column:name"`
		Label                      string   `json:"label" gorm:"column:label"`
		Region                     string   `json:"region" gorm:"column:region"`
		AvailableCnis              []string `json:"available_cnis" gorm:"column:available_cnis"`
		AvailableIngresses         []string `json:"available_ingresses" gorm:"column:available_ingresses"`
		AvailableContainerRuntimes []string `json:"available_container_runtimes" gorm:"column:available_container_runtimes"`
		AvailableFeatureGates      []string `json:"available_feature_gates" gorm:"column:available_feature_gates"`
		AvailableAdmissionPlugins  []string `json:"available_admission_plugins" gorm:"column:available_admission_plugins"`
		AvailableKubeletArgs       struct {
			AvailableKubeletArgKey string `json:"<available_kubelet_argKey>" gorm:"column:available_kubelet_arg_key"`
		} `json:"available_kubelet_args" gorm:"embedded"`
	} `json:"versions" gorm:"embedded"`
}

type ScalewayCreateResponse struct {
	ID               string    `json:"id" gorm:"column:id"`
	Type             string    `json:"type" gorm:"column:type"`
	Name             string    `json:"name" gorm:"column:name"`
	Status           string    `json:"status" gorm:"column:status"`
	Version          string    `json:"version" gorm:"column:version"`
	Region           string    `json:"region" gorm:"column:region"`
	OrganizationID   string    `json:"organization_id" gorm:"column:organization_id"`
	ProjectID        string    `json:"project_id" gorm:"column:project_id"`
	Tags             []string  `json:"tags" gorm:"column:tags"`
	Cni              string    `json:"cni" gorm:"column:cni"`
	Description      string    `json:"description" gorm:"column:description"`
	ClusterURL       string    `json:"cluster_url" gorm:"column:cluster_url"`
	DNSWildcard      string    `json:"dns_wildcard" gorm:"column:dns_wildcard"`
	CreatedAt        time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"column:updated_at"`
	AutoscalerConfig struct {
		ScaleDownDisabled             bool    `json:"scale_down_disabled" gorm:"column:scale_down_disabled"`
		ScaleDownDelayAfterAdd        string  `json:"scale_down_delay_after_add" gorm:"column:scale_down_delay_after_add"`
		Estimator                     string  `json:"estimator" gorm:"column:estimator"`
		Expander                      string  `json:"expander" gorm:"column:expander"`
		IgnoreDaemonSetsUtilization   bool    `json:"ignore_daemonSets_utilization" gorm:"column:ignore_daemonSets_utilization"`
		BalanceSimilarNodeGroups      bool    `json:"balance_similar_node_groups" gorm:"column:balance_similar_node_groups"`
		ExpendablePodsPriorityCutoff  int     `json:"expendable_pods_priority_cutoff" gorm:"column:expendable_pods_priority_cutoff"`
		ScaleDownUnneededTime         int     `json:"scale_down_unneeded_time" gorm:"column:scale_down_unneeded_time"`
		ScaleDownUtilizationThreshold float32 `json:"scale_down_utilization_threshold" gorm:"column:scale_down_utilization_threshold"`
		MaxGracefulTerminationSec     int     `json:"max_graceful_termination_sec" gorm:"column:max_graceful_termination_sec"`
	} `json:"autoscaler_config" gorm:"embedded"`
	DashboardEnabled bool   `json:"dashboard_enabled" gorm:"column:dashboard_enabled"`
	Ingress          string `json:"ingress" gorm:"column:ingress"`
	AutoUpgrade      struct {
		Enabled           bool `json:"enabled" gorm:"column:enabled"`
		MaintenanceWindow struct {
			StartHour int    `json:"start_hour" gorm:"column:start_hour"`
			Day       string `json:"day" gorm:"column:day"`
		} `json:"maintenance_window" gorm:"embedded"`
	} `json:"auto_upgrade" gorm:"embedded"`
	UpgradeAvailable    string   `json:"upgrade_available" gorm:"column:upgrade_available"`
	FeatureGates        []string `json:"feature_gates" gorm:"column:feature_gates"`
	AdmissionPlugins    []string `json:"admission_plugins" gorm:"column:admission_plugins"`
	OpenIDConnectConfig struct {
		IssuerURL      string   `json:"issuer_url" gorm:"column:issuer_url"`
		ClientID       string   `json:"client_id" gorm:"column:client_id"`
		UsernameClaim  string   `json:"username_claim" gorm:"column:username_claim"`
		UsernamePrefix string   `json:"username_prefix" gorm:"column:username_prefix"`
		GroupsClaim    []string `json:"groups_claim" gorm:"column:groups_claim"`
		GroupsPrefix   string   `json:"groups_prefix" gorm:"column:groups_prefix"`
		RequiredClaim  []string `json:"required_claim" gorm:"column:required_claim"`
	} `json:"open_id_connect_config" gorm:"embedded"`
	ApiServerCertSans []string `json:"apiServer_cert_sans" gorm:"column:apiServer_cert_sans"`
}

type ScalewayKubeconfigResponse struct {
	Name        string `json:"name" gorm:"column:name"`
	ContentType string `json:"content_type" gorm:"column:content_type"`
	Content     string `json:"content" gorm:"column:content"`
}

type JokeResponse struct {
	Categories []string `json:"categories"`
	CreatedAt  string   `json:"created_at"`
	IconURL    string   `json:"icon_url"`
	ID         string   `json:"id"`
	UpdatedAt  string   `json:"updated_at"`
	URL        string   `json:"url"`
	Value      string   `json:"value"`
}

type KuberoIngress struct {
	Items []struct {
		Kind string `json:"kind"`
		Spec struct {
			IngressClassName string `json:"ingressClassName"`
			Rules            []struct {
				Host string `json:"host"`
				HTTP struct {
					Paths []struct {
						Backend struct {
							Service struct {
								Name string `json:"name"`
								Port struct {
									Number int `json:"number"`
								} `json:"port"`
							} `json:"service"`
						} `json:"backend"`
						Path string `json:"path"`
					} `json:"paths"`
				} `json:"http"`
			} `json:"rules"`
		} `json:"spec"`
		Status struct {
			LoadBalancer struct {
				Ingress []struct {
					Hostname string `json:"hostname"`
					IP       string `json:"ip"`
				} `json:"ingress"`
			} `json:"loadBalancer"`
		} `json:"status"`
	} `json:"items"`
}

type LinodeCreateClusterRequest struct {
	Label        string   `json:"label"`
	Region       string   `json:"region"`
	K8SVersion   string   `json:"k8s_version"`
	Tags         []string `json:"tags"`
	ControlPlane struct {
		HighAvailability bool `json:"high_availability"`
	} `json:"control_plane"`
	NodePools []LinodeNodepool `json:"node_pools"`
}

type LinodeCreateClusterResponse struct {
	ControlPlane struct {
		HighAvailability bool `json:"high_availability"`
	} `json:"control_plane"`
	Created    time.Time `json:"created"`
	ID         int       `json:"id"`
	K8SVersion string    `json:"k8s_version"`
	Label      string    `json:"label"`
	Region     string    `json:"region"`
	Tags       []string  `json:"tags"`
	Updated    time.Time `json:"updated"`
}

type LinodeNodepool struct {
	Type       string `json:"type"`
	Count      int    `json:"count"`
	Autoscaler struct {
		Enabled bool `json:"enabled"`
		Max     int  `json:"max,omitempty"`
		Min     int  `json:"min,omitempty"`
	} `json:"autoscaler,omitempty"`
}

type CertManagerClusterIssuer struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Acme struct {
			Server              string `yaml:"server"`
			Email               string `yaml:"email"`
			PrivateKeySecretRef struct {
				Name string `yaml:"name"`
			} `yaml:"privateKeySecretRef"`
			Solvers []struct {
				HTTP01 struct {
					Ingress struct {
						Class string `yaml:"class"`
					} `yaml:"ingress"`
				} `yaml:"http01"`
			} `yaml:"solvers"`
		} `yaml:"acme"`
	} `yaml:"spec"`
}
