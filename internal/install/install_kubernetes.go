package install

import (
	"github.com/i582/cfmt/cmd/cfmt"
	"os"
)

func installKubernetes() {
	kubernetesInstall := promptLine("1) Create a kubernetes cluster", "[y,n]", "y")
	if kubernetesInstall != "y" {
		return
	}

	clusterType = selectFromList("Select a Kubernetes provider", clusterTypeList, "")

	switch clusterType {
	case "scaleway":
		installScaleway()
	case "linode":
		installLinode()
	case "gke":
		installGKE()
	case "digitalocean":
		installDigitalOcean()
	case "kind":
		installKind()
	default:
		_, _ = cfmt.Println("{{✗ Unknown cluster type}}::red")
		os.Exit(1)
	}
}
