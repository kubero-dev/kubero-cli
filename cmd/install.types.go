package cmd

import "time"

type DigitalOceanKubernetesConfig struct {
	Name      string `json:"name"`
	Region    string `json:"region"`
	Version   string `json:"version"`
	NodePools []struct {
		Size  string `json:"size"`
		Count int    `json:"count"`
		Name  string `json:"name"`
	} `json:"node_pools"`
}

type DigitalOcean struct {
	KubernetesCluster struct {
		ID            string   `json:"id"`
		Name          string   `json:"name"`
		Region        string   `json:"region"`
		Version       string   `json:"version"`
		ClusterSubnet string   `json:"cluster_subnet"`
		ServiceSubnet string   `json:"service_subnet"`
		VpcUUID       string   `json:"vpc_uuid"`
		Ipv4          string   `json:"ipv4"`
		Endpoint      string   `json:"endpoint"`
		Tags          []string `json:"tags"`
		NodePools     []struct {
			ID        string        `json:"id"`
			Name      string        `json:"name"`
			Size      string        `json:"size"`
			Count     int           `json:"count"`
			Tags      []string      `json:"tags"`
			Labels    interface{}   `json:"labels"`
			Taints    []interface{} `json:"taints"`
			AutoScale bool          `json:"auto_scale"`
			MinNodes  int           `json:"min_nodes"`
			MaxNodes  int           `json:"max_nodes"`
			Nodes     []struct {
				ID     string `json:"id"`
				Name   string `json:"name"`
				Status struct {
					State string `json:"state"`
				} `json:"status"`
				DropletID string    `json:"droplet_id"`
				CreatedAt time.Time `json:"created_at"`
				UpdatedAt time.Time `json:"updated_at"`
			} `json:"nodes"`
		} `json:"node_pools"`
		MaintenancePolicy struct {
			StartTime string `json:"start_time"`
			Duration  string `json:"duration"`
			Day       string `json:"day"`
		} `json:"maintenance_policy"`
		AutoUpgrade bool `json:"auto_upgrade"`
		Status      struct {
			State   string `json:"state"`
			Message string `json:"message"`
		} `json:"status"`
		CreatedAt         time.Time `json:"created_at"`
		UpdatedAt         time.Time `json:"updated_at"`
		SurgeUpgrade      bool      `json:"surge_upgrade"`
		RegistryEnabled   bool      `json:"registry_enabled"`
		Ha                bool      `json:"ha"`
		SupportedFeatures []string  `json:"supported_features"`
	} `json:"kubernetes_cluster"`
}

type User struct {
	ID       int    `json:"id"`
	Method   string `json:"method"`
	Username string `json:"username"`
	Password string `json:"password"`
	Insecure bool   `json:"insecure"`
	Apitoken string `json:"apitoken,omitempty"`
}

type KindConfig struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
	Name       string `yaml:"name"`
	Networking struct {
		IPFamily         string `yaml:"ipFamily"`
		APIServerAddress string `yaml:"apiServerAddress"`
	} `yaml:"networking"`
	Nodes []struct {
		Role                 string   `yaml:"role"`
		Image                string   `yaml:"image,omitempty"`
		KubeadmConfigPatches []string `yaml:"kubeadmConfigPatches"`
		ExtraPortMappings    []struct {
			ContainerPort int    `yaml:"containerPort"`
			HostPort      int    `yaml:"hostPort"`
			Protocol      string `yaml:"protocol"`
		} `yaml:"extraPortMappings"`
	} `yaml:"nodes"`
}

type KuberoUIConfig struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Affinity struct {
		} `yaml:"affinity"`
		FullnameOverride string `yaml:"fullnameOverride"`
		Image            struct {
			PullPolicy string `yaml:"pullPolicy"`
			Repository string `yaml:"repository"`
			Tag        string `yaml:"tag"`
		} `yaml:"image"`
		ImagePullSecrets []interface{} `yaml:"imagePullSecrets"`
		Ingress          struct {
			Annotations struct {
			} `yaml:"annotations"`
			ClassName string `yaml:"className"`
			Enabled   bool   `yaml:"enabled"`
			Hosts     []struct {
				Host  string `yaml:"host"`
				Paths []struct {
					Path     string `yaml:"path"`
					PathType string `yaml:"pathType"`
				} `yaml:"paths"`
			} `yaml:"hosts"`
			TLS []interface{} `yaml:"tls"`
		} `yaml:"ingress"`
		NameOverride string `yaml:"nameOverride"`
		NodeSelector struct {
		} `yaml:"nodeSelector"`
		PodAnnotations struct {
		} `yaml:"podAnnotations"`
		PodSecurityContext struct {
		} `yaml:"podSecurityContext"`
		ReplicaCount int `yaml:"replicaCount"`
		Resources    struct {
		} `yaml:"resources"`
		SecurityContext struct {
		} `yaml:"securityContext"`
		Service struct {
			Port int    `yaml:"port"`
			Type string `yaml:"type"`
		} `yaml:"service"`
		ServiceAccount struct {
			Annotations struct {
			} `yaml:"annotations"`
			Create bool   `yaml:"create"`
			Name   string `yaml:"name"`
		} `yaml:"serviceAccount"`
		Tolerations []interface{} `yaml:"tolerations"`
		Kubero      struct {
			Debug      string `yaml:"debug"`
			Namespace  string `yaml:"namespace"`
			Context    string `yaml:"context"`
			WebhookURL string `yaml:"webhook_url"`
			SessionKey string `yaml:"sessionKey"`
			Auth       struct {
				Github struct {
					Enabled     bool   `yaml:"enabled"`
					ID          string `yaml:"id"`
					Secret      string `yaml:"secret"`
					CallbackURL string `yaml:"callbackUrl"`
					Org         string `yaml:"org"`
				} `yaml:"github"`
				Oauth2 struct {
					Enabled     bool   `yaml:"enabled"`
					Name        string `yaml:"name"`
					ID          string `yaml:"id"`
					AuthURL     string `yaml:"authUrl"`
					TokenURL    string `yaml:"tokenUrl"`
					Secret      string `yaml:"secret"`
					CallbackURL string `yaml:"callbackUrl"`
				} `yaml:"oauth2"`
			} `yaml:"auth"`
			Config string `yaml:"config"`
		} `yaml:"kubero"`
	} `yaml:"spec"`
}

// https://developers.scaleway.com/en/products/k8s/api/#post-612200
type ScalewayCreate struct {
	OrganizationID  string             `json:"organization_id,omitempty"` // DEPRECATED
	ProjectID       string             `json:"project_id,omitempty"`      // REQUIRED
	Type            string             `json:"type"`
	Name            string             `json:"name"` // REQUIRED
	Description     string             `json:"description"`
	Tags            []string           `json:"tags"`
	Version         string             `json:"version"`          // REQUIRED
	Cni             string             `json:"cni"`              // REQUIRED
	EnableDashboard bool               `json:"enable_dashboard"` // DEPRECATED
	Ingress         string             `json:"ingress"`          // DEPRECATED
	Pools           []ScalewayNodePool `json:"pools"`
	/*
		AutoscalerConfig struct {
			ScaleDownDisabled             bool    `json:"scale_down_disabled"`
			ScaleDownDelayAfterAdd        string  `json:"scale_down_delay_after_add"`
			Estimator                     string  `json:"estimator"`
			Expander                      string  `json:"expander"`
			IgnoreDaemonsetsUtilization   bool    `json:"ignore_daemonsets_utilization"`
			BalanceSimilarNodeGroups      bool    `json:"balance_similar_node_groups"`
			ExpendablePodsPriorityCutoff  int     `json:"expendable_pods_priority_cutoff"`
			ScaleDownUnneededTime         string  `json:"scale_down_unneeded_time"`
			ScaleDownUtilizationThreshold float32 `json:"scale_down_utilization_threshold"`
			MaxGracefulTerminationSec     int     `json:"max_graceful_termination_sec"`
		} `json:"autoscaler_config,omitempty"`
	*/
	AutoUpgrade struct {
		Enable            bool `json:"enable"`
		MaintenanceWindow struct {
			StartHour int    `json:"start_hour"`
			Day       string `json:"day"`
		} `json:"maintenance_window"`
	} `json:"auto_upgrade"`
	FeatureGates     []string `json:"feature_gates"`
	AdmissionPlugins []string `json:"admission_plugins"`
	/*
		OpenIDConnectConfig struct {
			IssuerURL      string   `json:"issuer_url"`
			ClientID       string   `json:"client_id"`
			UsernameClaim  string   `json:"username_claim"`
			UsernamePrefix string   `json:"username_prefix"`
			GroupsClaim    []string `json:"groups_claim"`
			GroupsPrefix   string   `json:"groups_prefix"`
			RequiredClaim  []string `json:"required_claim"`
		} `json:"open_id_connect_config"`
	*/
	ApiserverCertSans []string `json:"apiserver_cert_sans"`
}

type ScalewayNodePool struct {
	Name             string   `json:"name"`
	NodeType         string   `json:"node_type"`
	PlacementGroupID string   `json:"placement_group_id,omitempty"`
	Autoscaling      bool     `json:"autoscaling"`
	Size             int      `json:"size"`
	MinSize          int      `json:"min_size"`
	MaxSize          int      `json:"max_size"`
	ContainerRuntime string   `json:"container_runtime"`
	Autohealing      bool     `json:"autohealing"`
	Tags             []string `json:"tags"`
	/* fails {"details":[{"argument_name":"default_pool_config[0].kubelet_args.\u003ckubelet_argKey\u003e","help_message":"kubelet argument \u003ckubelet_argKey\u003e is not available for this version","reason":"constraint"}
	KubeletArgs      struct {
		KubeletArgKey string `json:"<kubelet_argKey>"`
	} `json:"kubelet_args"`
	*/
	/* fails with {"argument_name":"default_pool_config[0].upgrade_policy.max_unavailable","help_message":"value must be between 1 and 20","reason":"constraint"}],"message":"invalid argument(s)","type":"invalid_arguments"}
	UpgradePolicy struct {
		MaxUnavailable int `json:"max_unavailable"`
		MaxSurge       int `json:"max_surge"`
	} `json:"upgrade_policy"`
	*/
	Zone           string `json:"zone"`
	RootVolumeType string `json:"root_volume_type"`
	RootVolumeSize int    `json:"root_volume_size"`
}

type ScalewayVersionsResponse struct {
	Versions []struct {
		Name                       string   `json:"name"`
		Label                      string   `json:"label"`
		Region                     string   `json:"region"`
		AvailableCnis              []string `json:"available_cnis"`
		AvailableIngresses         []string `json:"available_ingresses"`
		AvailableContainerRuntimes []string `json:"available_container_runtimes"`
		AvailableFeatureGates      []string `json:"available_feature_gates"`
		AvailableAdmissionPlugins  []string `json:"available_admission_plugins"`
		AvailableKubeletArgs       struct {
			AvailableKubeletArgKey string `json:"<available_kubelet_argKey>"`
		} `json:"available_kubelet_args"`
	} `json:"versions"`
}

type ScalewayCreateResponse struct {
	ID               string    `json:"id"`
	Type             string    `json:"type"`
	Name             string    `json:"name"`
	Status           string    `json:"status"`
	Version          string    `json:"version"`
	Region           string    `json:"region"`
	OrganizationID   string    `json:"organization_id"`
	ProjectID        string    `json:"project_id"`
	Tags             []string  `json:"tags"`
	Cni              string    `json:"cni"`
	Description      string    `json:"description"`
	ClusterURL       string    `json:"cluster_url"`
	DNSWildcard      string    `json:"dns_wildcard"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	AutoscalerConfig struct {
		ScaleDownDisabled             bool    `json:"scale_down_disabled"`
		ScaleDownDelayAfterAdd        string  `json:"scale_down_delay_after_add"`
		Estimator                     string  `json:"estimator"`
		Expander                      string  `json:"expander"`
		IgnoreDaemonsetsUtilization   bool    `json:"ignore_daemonsets_utilization"`
		BalanceSimilarNodeGroups      bool    `json:"balance_similar_node_groups"`
		ExpendablePodsPriorityCutoff  int     `json:"expendable_pods_priority_cutoff"`
		ScaleDownUnneededTime         int     `json:"scale_down_unneeded_time"`
		ScaleDownUtilizationThreshold float32 `json:"scale_down_utilization_threshold"`
		MaxGracefulTerminationSec     int     `json:"max_graceful_termination_sec"`
	} `json:"autoscaler_config"`
	DashboardEnabled bool   `json:"dashboard_enabled"`
	Ingress          string `json:"ingress"`
	AutoUpgrade      struct {
		Enabled           bool `json:"enabled"`
		MaintenanceWindow struct {
			StartHour int    `json:"start_hour"`
			Day       string `json:"day"`
		} `json:"maintenance_window"`
	} `json:"auto_upgrade"`
	UpgradeAvailable    string   `json:"upgrade_available"`
	FeatureGates        []string `json:"feature_gates"`
	AdmissionPlugins    []string `json:"admission_plugins"`
	OpenIDConnectConfig struct {
		IssuerURL      string   `json:"issuer_url"`
		ClientID       string   `json:"client_id"`
		UsernameClaim  string   `json:"username_claim"`
		UsernamePrefix string   `json:"username_prefix"`
		GroupsClaim    []string `json:"groups_claim"`
		GroupsPrefix   string   `json:"groups_prefix"`
		RequiredClaim  []string `json:"required_claim"`
	} `json:"open_id_connect_config"`
	ApiserverCertSans []string `json:"apiserver_cert_sans"`
}

type ScalewayKubeconfigResponse struct {
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
	Content     string `json:"content"`
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

type JokeResponse struct {
	Categories []string `json:"categories"`
	CreatedAt  string   `json:"created_at"`
	IconURL    string   `json:"icon_url"`
	ID         string   `json:"id"`
	UpdatedAt  string   `json:"updated_at"`
	URL        string   `json:"url"`
	Value      string   `json:"value"`
}

type CertmanagerClusterIssuer struct {
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
