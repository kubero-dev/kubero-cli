package install

import (
	"github.com/kubero-dev/kubero-cli/internal/log"
	"os/exec"
)

func (m *ManagerInstall) InstallMetrics() error {
	installed, _ := exec.Command("kubectl", "get", "deployments.apps", "metrics-server", "-n", "kube-system").Output()
	if len(installed) > 0 {
		log.Info("Metrics server already installed")
		return nil
	}
	install := promptLine("5) Install Kubernetes internal metrics service (required for HPA, Horizontal Pod Autoscaling)", "[y,n]", "y")
	if install != "y" {
		log.Println("Skipping metrics server installation")
		return nil
	}
	components := "https://raw.githubusercontent.com/kubero-dev/kubero-operator/main/deploy/metrics-server.yaml"
	_, installErr := exec.Command("kubectl", "apply", "-f", components).Output()
	if installErr != nil {
		log.Fatal("Failed to install metrics server", installErr)
		return installErr
	}
	log.Info("Metrics server installed successfully")
	return nil
}
