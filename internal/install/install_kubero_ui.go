package install

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/i582/cfmt/cmd/cfmt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/exec"

	"time"
)

func installKuberoUi() {
	ingressInstall := promptLine("9) Install Kubero UI", "[y,n]", "y")
	if ingressInstall != "y" {
		return
	}

	createNamespace("kubero")

	kuberoSecretInstalled, _ := exec.Command("kubectl", "get", "secret", "kubero-secrets", "-n", "kubero").Output()
	if len(kuberoSecretInstalled) > 0 {
		_, _ = cfmt.Println("{{✓ Kubero Secret exists}}::lightGreen")
	} else {
		webhookSecret := promptLine("Random string for your webhook secret", "", generateRandomString(20, ""))

		sessionKey := promptLine("Random string for your session key", "", generateRandomString(20, ""))

		if argAdminUser == "" {
			argAdminUser = promptLine("Admin User", "", "admin")
		}

		if argAdminPassword == "" {
			argAdminPassword = promptLine("Admin Password", "", generateRandomString(12, ""))
		}

		if argApiToken == "" {
			argApiToken = promptLine("Random string for admin API token", "", generateRandomString(20, ""))
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(argAdminPassword), bcrypt.DefaultCost)
		userData := map[string]string{
			"username": argAdminUser,
			"password": string(hashedPassword),
		}
		userDataJson, _ := json.Marshal(userData)
		userDataBase64 := base64.StdEncoding.EncodeToString(userDataJson)

		kuberoSecrets := map[string]string{
			"sessionKey":     sessionKey,
			"webhookSecret":  webhookSecret,
			"userDataBase64": userDataBase64,
			"apiToken":       argApiToken,
		}

		secretsYaml, _ := yaml.Marshal(kuberoSecrets)
		secretsYamlErr := os.WriteFile("kuberoSecrets.yaml", secretsYaml, 0644)
		if secretsYamlErr != nil {
			fmt.Println(secretsYamlErr)
			return
		}

		_, secretsErr := exec.Command("kubectl", "create", "secret", "generic", "kubero-secrets", "--from-file=kuberoSecrets.yaml", "-n", "kubero").Output()
		if secretsErr != nil {
			_, _ = cfmt.Println("{{✗ Failed to create Kubero secrets}}::red")
			log.Fatal(secretsErr)
		}

		e := os.Remove("kuberoSecrets.yaml")
		if e != nil {
			log.Fatal(e)
		}
	}

	kuberoUiInstalled, _ := exec.Command("kubectl", "get", "deployments.apps", "kubero-ui", "-n", "kubero").Output()
	if len(kuberoUiInstalled) > 0 {
		_, _ = cfmt.Println("{{✓ Kubero UI is already installed}}::lightGreen")
		return
	}

	kuberoUiSpinner := spinner.New("Install Kubero UI")
	kuberoUiUrl := "https://raw.githubusercontent.com/kubero-dev/kubero-ui/main/deploy/kubero-ui.yaml"
	_, _ = cfmt.Println("  run command : kubectl apply -f " + kuberoUiUrl)
	kuberoUiSpinner.Start("Installing Kubero UI")
	_, kuberoUiErr := exec.Command("kubectl", "apply", "-f", kuberoUiUrl).Output()
	if kuberoUiErr != nil {
		kuberoUiSpinner.Error("Failed to run command. Try running this command manually: kubectl apply -f " + kuberoUiUrl)
		log.Fatal(kuberoUiErr)
	}

	kuberoUiSpinner.UpdateMessage("Waiting for Kubero UI to be ready")
	time.Sleep(10 * time.Second)
	_, kuberoUiWaitErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/kubero-ui", "-n", "kubero", "--timeout=180s").Output()
	if kuberoUiWaitErr != nil {
		kuberoUiSpinner.Error("Failed to run command. Try running it manually: kubectl wait --for=condition=available deployment/kubero-ui -n kubero --timeout=180s")
		log.Fatal(kuberoUiWaitErr)
	}
	kuberoUiSpinner.Success("Kubero UI installed")
}
