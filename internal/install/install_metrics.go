package install

import (
	"fmt"
	"github.com/i582/cfmt/cmd/cfmt"
	"log"
	"os/exec"
)

func installMetrics() {
	installed, _ := exec.Command("kubectl", "get", "deployments.apps", "metrics-server", "-n", "kube-system").Output()
	if len(installed) > 0 {
		_, _ = cfmt.Println("{{✓ Metrics is already enabled}}::lightGreen")
		return
	}
	install := promptLine("5) Install Kubernetes internal metrics service (required for HPA, Horizontal Pod Autoscaling)", "[y,n]", "y")
	if install != "y" {
		return
	}

	components := "https://raw.githubusercontent.com/kubero-dev/kubero-operator/main/deploy/metrics-server.yaml"
	_, installErr := exec.Command("kubectl", "apply", "-f", components).Output()

	if installErr != nil {
		fmt.Println("failed to install metrics server")
		log.Fatal(installErr)
	}
	_, _ = cfmt.Println("{{✓ Metrics server installed}}::lightGreen")
}
