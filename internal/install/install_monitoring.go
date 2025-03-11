package install

import (
	"github.com/i582/cfmt/cmd/cfmt"
	"log"
	"os/exec"
	"time"
)

func installMonitoring() {
	install := promptLine("7) Install Monitoring", "[y,n]", "y")
	if install != "y" {
		return
	}

	monitoringInstalled, _ := exec.Command("kubectl", "get", "deployments.apps", "monitoring-stack-operators").Output()
	if len(monitoringInstalled) > 0 {
		_, _ = cfmt.Println("{{✓ Monitoring is already installed}}::lightGreen")
		return
	}

	monitoringSpinner := spinner.New("Install Monitoring")
	monitoringUrl := "https://raw.githubusercontent.com/kubero-dev/kubero-operator/main/deploy/monitoring.yaml"
	_, _ = cfmt.Println("  run command : kubectl apply -f " + monitoringUrl)
	monitoringSpinner.Start("Installing Monitoring")
	_, monitoringErr := exec.Command("kubectl", "apply", "-f", monitoringUrl).Output()
	if monitoringErr != nil {
		monitoringSpinner.Error("Failed to run command. Try running this command manually: kubectl apply -f " + monitoringUrl)
		log.Fatal(monitoringErr)
	}

	monitoringSpinner.UpdateMessage("Waiting for Monitoring to be ready")
	time.Sleep(10 * time.Second)
	_, monitoringWaitErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/monitoring-stack-operators", "--timeout=180s").Output()
	if monitoringWaitErr != nil {
		monitoringSpinner.Error("Failed to run command. Try running it manually: kubectl wait --for=condition=available deployment/monitoring-stack-operators --timeout=180s")
		log.Fatal(monitoringWaitErr)
	}
	monitoringSpinner.Success("Monitoring installed")
}
