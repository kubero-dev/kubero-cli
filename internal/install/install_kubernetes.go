package install

import (
	"fmt"
	"github.com/faelmori/kubero-cli/internal/log"
)

func installKubernetes() error {
	kubernetesInstall := promptLine("1) Create a kubernetes cluster", "[y,n]", "y")
	if kubernetesInstall != "y" {
		log.Println("Skipping Kubernetes cluster installation")
		return nil
	}
	clusterType = selectFromList("Select a Kubernetes provider", clusterTypeList, "")

	switch clusterType {
	case "scaleway":
		return installScaleway()
	case "linode":
		return installLinode()
	case "gke":
		return installGKE()
	case "digitalocean":
		return installDigitalOcean()
	case "kind":
		return installKind()
	default:
		return fmt.Errorf("invalid cluster type: %s", clusterType)
	}
}
