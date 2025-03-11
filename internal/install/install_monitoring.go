package install

import (
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/kubero-dev/kubero-cli/internal/log"
	"github.com/leaanthony/spinner"
	"os/exec"
	"time"
)

func (m *ManagerInstall) InstallMonitoring() error {
	install := promptLine("7) Install Monitoring", "[y,n]", "y")
	if install != "y" {
		log.Info("Skipping monitoring installation")
		return nil
	}

	monitoringInstalledCmd, _ := exec.Command("kubectl", "get", "deployments.apps", "monitoring-stack-operators").Output()
	if len(monitoringInstalledCmd) > 0 {
		log.Info("Monitoring is already installed")
		return nil
	}

	monitoringSpinner := spinner.New("Install Monitoring")
	monitoringUrl := "https://raw.githubusercontent.com/kubero-dev/kubero-operator/main/deploy/monitoring.yaml"
	_, _ = cfmt.Println("  run command : kubectl apply -f " + monitoringUrl)
	monitoringSpinner.Start("Installing Monitoring")
	_, monitoringErr := exec.Command("kubectl", "apply", "-f", monitoringUrl).Output()
	if monitoringErr != nil {
		monitoringSpinner.Error("Failed to run command. Try running this command manually: kubectl apply -f " + monitoringUrl)
		return monitoringErr
	}

	monitoringSpinner.UpdateMessage("Waiting for Monitoring to be ready")
	time.Sleep(10 * time.Second)
	_, monitoringWaitErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/monitoring-stack-operators", "--timeout=180s").Output()
	if monitoringWaitErr != nil {
		monitoringSpinner.Error("Failed to run command. Try running it manually: kubectl wait --for=condition=available deployment/monitoring-stack-operators --timeout=180s")
		return monitoringWaitErr
	}
	monitoringSpinner.Success("Monitoring installed")

	return nil
}
