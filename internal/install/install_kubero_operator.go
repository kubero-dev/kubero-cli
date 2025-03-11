package install

import (
	"fmt"
	"github.com/faelmori/kubero-cli/internal/log"
	"github.com/leaanthony/spinner"
	"os/exec"
	"time"
)

func (m *ManagerInstall) InstallKuberoOperator() error {
	log.Info("Installing KuberOS Operator")

	kuberoInstalled, _ := exec.Command("kubectl", "get", "operator", "kubero-operator.operators").Output()
	if len(kuberoInstalled) > 0 {
		log.Info("Kubero Operator is already installed")
		return nil
	}

	if m.installOlm {
		return m.InstallKuberoOLMOperator()
	} else {
		return m.installKuberoOperatorSlim()
	}
}

func (m *ManagerInstall) InstallKuberoOLMOperator() error {
	kuberoSpinner := spinner.New("Install Kubero Operator")
	log.Info("run command : kubectl apply -f https://operatorhub.io/install/kubero-operator.yaml")
	kuberoSpinner.Start("Install Kubero Operator")
	_, kuberoErr := exec.Command("kubectl", "apply", "-f", "https://operatorhub.io/install/kubero-operator.yaml").Output()
	if kuberoErr != nil {
		fmt.Println("")
		kuberoSpinner.Error("Failed to run command to install the Operator. Try running this command manually: kubectl apply -f https://operatorhub.io/install/kubero-operator.yaml")
		return kuberoErr
	}

	kuberoSpinner.UpdateMessage("Wait for Kubero Operator to be ready")
	var kuberoWait []byte
	for len(kuberoWait) == 0 {
		kuberoWait, _ = exec.Command("kubectl", "api-resources", "--api-group=application.kubero.dev", "--no-headers=true").Output()
		time.Sleep(1 * time.Second)
	}

	kuberoSpinner.Success("Kubero Operator installed successfully")

	return nil
}

func (m *ManagerInstall) installKuberoOperatorSlim() error {
	kuberoSpinner := spinner.New("Install Kubero Operator")
	log.Info("run command : kubectl apply -f https://raw.githubusercontent.com/kubero-dev/kubero-operator/main/deploy/operator.yaml")
	kuberoSpinner.Start("Install Kubero Operator")
	_, kuberoErr := exec.Command("kubectl", "apply", "-f", "https://raw.githubusercontent.com/kubero-dev/kubero-operator/main/deploy/operator.yaml").Output()
	if kuberoErr != nil {
		fmt.Println("")
		kuberoSpinner.Error("Failed to run command to install the Operator. Try running this command manually: kubectl apply -f https://raw.githubusercontent.com/kubero-dev/kubero-operator/main/deploy/operator.yaml")
		return kuberoErr
	}

	kuberoSpinner.UpdateMessage("Wait for Kubero Operator CRD's to be installed")
	var kuberoWait []byte
	for len(kuberoWait) == 0 {
		kuberoWait, _ = exec.Command("kubectl", "api-resources", "--api-group=application.kubero.dev", "--no-headers=true").Output()
		time.Sleep(1 * time.Second)
	}
	kuberoSpinner.UpdateMessage("Kubero Operator CRD's installed")

	time.Sleep(5 * time.Second)
	kuberoSpinner.UpdateMessage("Wait for Kubero Operator to be ready")
	_, olmWaitErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/kubero-operator-controller-manager", "-n", "kubero-operator-system", "--timeout=300s").Output()
	if olmWaitErr != nil {
		kuberoSpinner.Error("Failed to wait for Kubero UI to become ready")
		return olmWaitErr
	}
	kuberoSpinner.Success("Kubero Operator installed successfully")

	return nil
}
