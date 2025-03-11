package install

import (
	"github.com/faelmori/kubero-cli/internal/log"
	"github.com/go-resty/resty/v2"
	"github.com/leaanthony/spinner"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"time"
)

func (m *ManagerInstall) InstallCertManager() error {
	install := promptLine("6) Install SSL CertManager", "[y,n]", "y")
	if install != "y" {
		log.Info("Skipping CertManager installation")
		return nil
	}

	if m.installOlm {
		return m.installOLMCertManager()
	} else {
		return m.installCertManagerSlim()
	}
}

func (m *ManagerInstall) installCertManagerSlim() error {
	kuberoUIInstalled, _ := exec.Command("kubectl", "get", "crd", "certificates.cert-manager.io").Output()
	if len(kuberoUIInstalled) > 0 {
		log.Info("CertManager already installed")
		return nil
	}

	certManagerSpinner := spinner.New("Install Cert Manager")
	certManagerUrl := "https://github.com/cert-manager/cert-manager/releases/latest/download/cert-manager.yaml"
	log.Info("run command : kubectl create -f " + certManagerUrl)
	certManagerSpinner.Start("Installing Cert Manager")
	_, certManagerErr := exec.Command("kubectl", "create", "-f", certManagerUrl).Output()
	if certManagerErr != nil {
		certManagerSpinner.Error("Failed to run command. Try running this command manually: kubectl create -f " + certManagerUrl)
		return certManagerErr
	}

	certManagerSpinner.UpdateMessage("Waiting for Cert Manager to be ready")
	time.Sleep(10 * time.Second)
	_, certManagerWaitErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/cert-manager-webhook", "-n", "cert-manager", "--timeout=180s", "-n", "cert-manager").Output()
	if certManagerWaitErr != nil {
		certManagerSpinner.Error("Failed to run command. Try running it manually: kubectl wait --for=condition=available deployment/cert-manager-webhook -n cert-manager --timeout=180s -n cert-manager")
		return certManagerWaitErr
	}
	certManagerSpinner.Success("Cert Manager installed")

	return m.installCertManagerClusterIssuer("cert-manager")
}

func (m *ManagerInstall) installCertManagerClusterIssuer(namespace string) error {
	installer := resty.New()

	installer.SetBaseURL("https://raw.githubusercontent.com")
	kf, _ := installer.R().Get("kubero-dev/kubero-cli/main/templates/certManagerClusterIssuer.prod.yaml")

	var certManagerClusterIssuer CertManagerClusterIssuer
	_ = yaml.Unmarshal(kf.Body(), &certManagerClusterIssuer)

	argCertManagerContact := promptLine("6.1) LetsEncrypt ACME contact email", "", "noreply@yourdomain.com")
	certManagerClusterIssuer.Spec.Acme.Email = argCertManagerContact

	clusterIssuer := promptLine("6.2) ClusterIssuer Name", "", "letsencrypt-prod")
	certManagerClusterIssuer.Metadata.Name = clusterIssuer

	certManagerClusterIssuerYaml, _ := yaml.Marshal(certManagerClusterIssuer)
	certManagerClusterIssuerYamlErr := os.WriteFile("kuberoCertManagerClusterIssuer.yaml", certManagerClusterIssuerYaml, 0644)
	if certManagerClusterIssuerYamlErr != nil {
		log.Error("Failed to write CertManager ClusterIssuer yaml file")
		return certManagerClusterIssuerYamlErr
	}

	_, certManagerClusterIssuerErr := exec.Command("kubectl", "apply", "-f", "kuberoCertManagerClusterIssuer.yaml", "-n", namespace).Output()
	if certManagerClusterIssuerErr != nil {
		log.Error("Failed to create CertManager ClusterIssuer. Try running this command manually: kubectl apply -f kuberoCertManagerClusterIssuer.yaml -n cert-manager")
		return certManagerClusterIssuerErr
	} else {
		e := os.Remove("kuberoCertManagerClusterIssuer.yaml")
		if e != nil {
			log.Error("Failed to remove CertManager ClusterIssuer yaml file")
			return e
		}

		log.Info("Cert Manager Cluster Issuer created")

		return nil
	}
}

func (m *ManagerInstall) installOLMCertManager() error {
	certManagerInstalled, _ := exec.Command("kubectl", "get", "deployment", "cert-manager-webhook", "-n", "operators").Output()
	if len(certManagerInstalled) > 0 {
		log.Info("Cert Manager already installed")
		return nil
	}

	certManagerSpinner := spinner.New("Install Cert Manager")
	log.Info("run command : kubectl create -f https://operatorhub.io/install/cert-manager.yaml")
	certManagerSpinner.Start("Installing Cert Manager")
	_, certManagerErr := exec.Command("kubectl", "create", "-f", "https://operatorhub.io/install/cert-manager.yaml").Output()
	if certManagerErr != nil {
		certManagerSpinner.Error("Failed to run command. Try running this command manually: kubectl create -f https://operatorhub.io/install/cert-manager.yaml")
		return certManagerErr
	}
	certManagerSpinner.Success("Cert Manager installed")

	certManagerSpinner = spinner.New("Wait for Cert Manager to be ready")
	certManagerSpinner.Start("installing Cert Manager")

	log.Info("run command : kubectl wait --for=condition=available deployment/cert-manager-webhook -n cert-manager --timeout=180s -n operators")

	_, certManagerWaitErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/cert-manager-webhook", "-n", "cert-manager", "--timeout=180s", "-n", "operators").Output()
	if certManagerWaitErr != nil {
		certManagerSpinner.Error("Failed to run command. Try running it manually: kubectl wait --for=condition=available deployment/cert-manager-webhook -n cert-manager --timeout=180s -n operators")
		return certManagerWaitErr
	}
	certManagerSpinner.Success("Cert Manager is ready")

	return m.installCertManagerClusterIssuer("default")
}
