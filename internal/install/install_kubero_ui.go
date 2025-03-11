package install

import (
	"encoding/base64"
	"encoding/json"
	"github.com/faelmori/kubero-cli/internal/log"
	"github.com/leaanthony/spinner"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"

	"time"
)

func (m *ManagerInstall) InstallKuberoUi() error {
	ingressInstall := promptLine("9) Install Kubero UI", "[y,n]", "y")
	if ingressInstall != "y" {
		log.Info("Skipping Kubero UI installation")
		return nil
	}

	if createNamespaceErr := utils.CreateNamespace("kubero"); createNamespaceErr != nil {
		return createNamespaceErr
	}

	kuberoSecretInstalled, _ := exec.Command("kubectl", "get", "secret", "kubero-secrets", "-n", "kubero").Output()
	if len(kuberoSecretInstalled) > 0 {
		log.Info("Kubero Secret exists")
	} else {
		webhookSecret := promptLine("Random string for your webhook secret", "", utils.GenerateRandomString(20, ""))

		sessionKey := promptLine("Random string for your session key", "", utils.GenerateRandomString(20, ""))

		if m.argAdminUser == "" {
			m.argAdminUser = promptLine("Admin User", "", "admin")
		}

		if m.argAdminPassword == "" {
			m.argAdminPassword = promptLine("Admin Password", "", utils.GenerateRandomString(12, ""))
		}

		if m.argApiToken == "" {
			m.argApiToken = promptLine("Random string for admin API token", "", utils.GenerateRandomString(20, ""))
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(m.argAdminPassword), bcrypt.DefaultCost)
		userData := map[string]string{
			"username": m.argAdminUser,
			"password": string(hashedPassword),
		}
		userDataJson, _ := json.Marshal(userData)
		userDataBase64 := base64.StdEncoding.EncodeToString(userDataJson)

		kuberoSecrets := map[string]string{
			"sessionKey":     sessionKey,
			"webhookSecret":  webhookSecret,
			"userDataBase64": userDataBase64,
			"apiToken":       m.argApiToken,
		}

		secretsYaml, _ := yaml.Marshal(kuberoSecrets)
		secretsYamlErr := os.WriteFile("kuberoSecrets.yaml", secretsYaml, 0644)
		if secretsYamlErr != nil {
			log.Error("Failed to write Kubero secrets to file")
			return secretsYamlErr
		}

		_, secretsErr := exec.Command("kubectl", "create", "secret", "generic", "kubero-secrets", "--from-file=kuberoSecrets.yaml", "-n", "kubero").Output()
		if secretsErr != nil {
			log.Error("Failed to create Kubero secrets")
			return secretsErr
		}

		e := os.Remove("kuberoSecrets.yaml")
		if e != nil {
			log.Error("Failed to remove kuberoSecrets.yaml")
			return e
		}
	}

	kuberoUiInstalled, _ := exec.Command("kubectl", "get", "deployments.apps", "kubero-ui", "-n", "kubero").Output()
	if len(kuberoUiInstalled) > 0 {
		log.Info("Kubero UI is already installed")
		return nil
	}

	kuberoUiSpinner := spinner.New("Install Kubero UI")
	kuberoUiUrl := "https://raw.githubusercontent.com/kubero-dev/kubero-ui/main/deploy/kubero-ui.yaml"

	log.Info("run command : kubectl apply -f " + kuberoUiUrl)
	kuberoUiSpinner.Start("Installing Kubero UI")
	_, kuberoUiErr := exec.Command("kubectl", "apply", "-f", kuberoUiUrl).Output()
	if kuberoUiErr != nil {
		kuberoUiSpinner.Error("Failed to run command. Try running this command manually: kubectl apply -f " + kuberoUiUrl)
		return kuberoUiErr
	}

	kuberoUiSpinner.UpdateMessage("Waiting for Kubero UI to be ready")
	time.Sleep(10 * time.Second)
	_, kuberoUiWaitErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/kubero-ui", "-n", "kubero", "--timeout=180s").Output()
	if kuberoUiWaitErr != nil {
		kuberoUiSpinner.Error("Failed to run command. Try running it manually: kubectl wait --for=condition=available deployment/kubero-ui -n kubero --timeout=180s")
		return kuberoUiWaitErr
	}
	kuberoUiSpinner.Success("Kubero UI installed")

	return nil
}
