package install

import (
	"fmt"
	l "github.com/kubero-dev/kubero-cli/internal/log"
)

func (m *ManagerInstall) InstallKubernetes() error {
	kubernetesInstall := promptLine("1) Create a kubernetes cluster", "[y,n]", "y")
	if kubernetesInstall != "y" {
		l.Println("Skipping Kubernetes cluster installation")
		return nil
	}
	m.clusterType = selectFromList("Select a Kubernetes provider", m.clusterTypeList, "")

	switch m.clusterType {
	case "scaleway":
		return m.installScaleway()
	case "linode":
		return m.installLinode()
	case "gke":
		return m.installGKE()
	case "digitalocean":
		return m.installDigitalOcean()
	case "kind":
		return m.installKind()
	default:
		return fmt.Errorf("invalid cluster type: %s", m.clusterType)
	}
}
