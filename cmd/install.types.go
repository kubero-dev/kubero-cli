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
				Config string `yaml:"config"`
				Kubero struct {
					Context   string `yaml:"context"`
					Namespace string `yaml:"namespace"`
					Port      int    `yaml:"port"`
				} `yaml:"kubero"`
				Buildpacks []struct {
					Name     string `yaml:"name"`
					Language string `yaml:"language"`
					Fetch    struct {
						Repository string `yaml:"repository"`
						Tag        string `yaml:"tag"`
					} `yaml:"fetch"`
					Build struct {
						Repository string `yaml:"repository"`
						Tag        string `yaml:"tag"`
						Command    string `yaml:"command"`
					} `yaml:"build"`
					Run struct {
						Repository         string `yaml:"repository"`
						Tag                string `yaml:"tag"`
						ReadOnlyAppStorage bool   `yaml:"readOnlyAppStorage"`
						SecurityContext    struct {
							AllowPrivilegeEscalation bool `yaml:"allowPrivilegeEscalation"`
							ReadOnlyRootFilesystem   bool `yaml:"readOnlyRootFilesystem"`
						} `yaml:"securityContext"`
						Command string `yaml:"command"`
					} `yaml:"run,omitempty"`
				} `yaml:"buildpacks"`
				PodSizeList []struct {
					Name        string `yaml:"name"`
					Description string `yaml:"description"`
					Default     bool   `yaml:"default,omitempty"`
					Resources   struct {
						Requests struct {
							Memory string `yaml:"memory"`
							CPU    string `yaml:"cpu"`
						} `yaml:"requests"`
						Limits struct {
							Memory string `yaml:"memory"`
							CPU    string `yaml:"cpu"`
						} `yaml:"limits"`
					} `yaml:"resources,omitempty"`
					Active bool `yaml:"active,omitempty"`
				} `yaml:"podSizeList"`
			} `yaml:"auth"`
		} `yaml:"kubero"`
	} `yaml:"spec"`
}
