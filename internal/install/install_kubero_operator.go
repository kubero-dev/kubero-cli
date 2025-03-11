package install

import (
	"fmt"
	"github.com/i582/cfmt/cmd/cfmt"
	"log"
	"os/exec"
	"time"
)

func installKuberoOperator() {
	_, _ = cfmt.Println("\n  {{3) Install Kubero Operator}}::bold")

	kuberoInstalled, _ := exec.Command("kubectl", "get", "operator", "kubero-operator.operators").Output()
	if len(kuberoInstalled) > 0 {
		_, _ = cfmt.Println("{{✓ Kubero Operator is already installed}}::lightGreen")
		return
	}

	if installOlm {
		installKuberoOLMOperator()
	} else {
		installKuberoOperatorSlim()
	}
}

func installKuberoOLMOperator() {
	kuberoSpinner := spinner.New("Install Kubero Operator")
	_, _ = cfmt.Println("  run command : kubectl apply -f https://operatorhub.io/install/kubero-operator.yaml")
	kuberoSpinner.Start("Install Kubero Operator")
	_, kuberoErr := exec.Command("kubectl", "apply", "-f", "https://operatorhub.io/install/kubero-operator.yaml").Output()
	if kuberoErr != nil {
		fmt.Println("")
		kuberoSpinner.Error("Failed to run command to install the Operator. Try running this command manually: kubectl apply -f https://operatorhub.io/install/kubero-operator.yaml")
		log.Fatal(kuberoErr)
	}

	kuberoSpinner.UpdateMessage("Wait for Kubero Operator to be ready")
	var kuberoWait []byte
	for len(kuberoWait) == 0 {
		kuberoWait, _ = exec.Command("kubectl", "api-resources", "--api-group=application.kubero.dev", "--no-headers=true").Output()
		time.Sleep(1 * time.Second)
	}

	kuberoSpinner.Success("Kubero Operator installed successfully")
}

func installKuberoOperatorSlim() {
	kuberoSpinner := spinner.New("Install Kubero Operator")
	_, _ = cfmt.Println("  run command : kubectl apply -f https://raw.githubusercontent.com/kubero-dev/kubero-operator/main/deploy/operator.yaml")
	kuberoSpinner.Start("Install Kubero Operator")
	_, kuberoErr := exec.Command("kubectl", "apply", "-f", "https://raw.githubusercontent.com/kubero-dev/kubero-operator/main/deploy/operator.yaml").Output()
	if kuberoErr != nil {
		fmt.Println("")
		kuberoSpinner.Error("Failed to run command to install the Operator. Try running this command manually: kubectl apply -f https://raw.githubusercontent.com/kubero-dev/kubero-operator/main/deploy/operator.yaml")
		log.Fatal(kuberoErr)
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
		log.Fatal(olmWaitErr)
	}
	kuberoSpinner.Success("Kubero Operator installed successfully")
}
