package install

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/leaanthony/spinner"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/exec"
	"time"
)

func installCertManager() {
	install := promptLine("6) Install SSL CertManager", "[y,n]", "y")
	if install != "y" {
		return
	}

	if installOlm {
		installOLMCertManager()
	} else {
		installCertManagerSlim()
	}
}

func installCertManagerSlim() {
	kuberoUIInstalled, _ := exec.Command("kubectl", "get", "crd", "certificates.cert-manager.io").Output()
	if len(kuberoUIInstalled) > 0 {
		_, _ = cfmt.Println("{{✓ CertManager already installed}}::lightGreen")
		return
	}

	certManagerSpinner := spinner.New("Install Cert Manager")
	certManagerUrl := "https://github.com/cert-manager/cert-manager/releases/latest/download/cert-manager.yaml"
	_, _ = cfmt.Println("  run command : kubectl create -f " + certManagerUrl)
	certManagerSpinner.Start("Installing Cert Manager")
	_, certManagerErr := exec.Command("kubectl", "create", "-f", certManagerUrl).Output()
	if certManagerErr != nil {
		certManagerSpinner.Error("Failed to run command. Try running this command manually: kubectl create -f " + certManagerUrl)
		log.Fatal(certManagerErr)
	}

	certManagerSpinner.UpdateMessage("Waiting for Cert Manager to be ready")
	time.Sleep(10 * time.Second)
	_, certManagerWaitErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/cert-manager-webhook", "-n", "cert-manager", "--timeout=180s", "-n", "cert-manager").Output()
	if certManagerWaitErr != nil {
		certManagerSpinner.Error("Failed to run command. Try running it manually: kubectl wait --for=condition=available deployment/cert-manager-webhook -n cert-manager --timeout=180s -n cert-manager")
		log.Fatal(certManagerWaitErr)
	}
	certManagerSpinner.Success("Cert Manager installed")

	installCertManagerClusterIssuer("cert-manager")
}

func installCertManagerClusterIssuer(namespace string) {
	installer := resty.New()

	installer.SetBaseURL("https://raw.githubusercontent.com")
	kf, _ := installer.R().Get("kubero-dev/kubero-cli/main/templates/certManagerClusterIssuer.prod.yaml")

	var certManagerClusterIssuer CertManagerClusterIssuer
	_ = yaml.Unmarshal(kf.Body(), &certManagerClusterIssuer)

	argCertManagerContact := promptLine("6.1) Letsencrypt ACME contact email", "", "noreply@yourdomain.com")
	certManagerClusterIssuer.Spec.Acme.Email = argCertManagerContact

	clusterIssuer := promptLine("6.2) ClusterIssuer Name", "", "letsencrypt-prod")
	certManagerClusterIssuer.Metadata.Name = clusterIssuer

	certManagerClusterIssuerYaml, _ := yaml.Marshal(certManagerClusterIssuer)
	certManagerClusterIssuerYamlErr := os.WriteFile("kuberoCertManagerClusterIssuer.yaml", certManagerClusterIssuerYaml, 0644)
	if certManagerClusterIssuerYamlErr != nil {
		fmt.Println(certManagerClusterIssuerYamlErr)
		return
	}

	_, certManagerClusterIssuerErr := exec.Command("kubectl", "apply", "-f", "kuberoCertManagerClusterIssuer.yaml", "-n", namespace).Output()
	if certManagerClusterIssuerErr != nil {
		_, _ = cfmt.Println("{{✗ Failed to create CertManager ClusterIssuer. Try running this command manually: kubectl apply -f kuberoCertManagerClusterIssuer.yaml -n cert-manager}}::red")
		return
	} else {
		e := os.Remove("kuberoCertManagerClusterIssuer.yaml")
		if e != nil {
			log.Fatal(e)
		}
		_, _ = cfmt.Println("{{✓ Cert Manager Cluster Issuer created}}::lightGreen")
	}
}

func installOLMCertManager() {
	certManagerInstalled, _ := exec.Command("kubectl", "get", "deployment", "cert-manager-webhook", "-n", "operators").Output()
	if len(certManagerInstalled) > 0 {
		_, _ = cfmt.Println("{{✓ Cert Manager already installed}}::lightGreen")
		return
	}

	certManagerSpinner := spinner.New("Install Cert Manager")
	_, _ = cfmt.Println("  run command : kubectl create -f https://operatorhub.io/install/cert-manager.yaml")
	certManagerSpinner.Start("Installing Cert Manager")
	_, certManagerErr := exec.Command("kubectl", "create", "-f", "https://operatorhub.io/install/cert-manager.yaml").Output()
	if certManagerErr != nil {
		certManagerSpinner.Error("Failed to run command. Try running this command manually: kubectl create -f https://operatorhub.io/install/cert-manager.yaml")
		log.Fatal(certManagerErr)
	}
	certManagerSpinner.Success("Cert Manager installed")

	certManagerSpinner = spinner.New("Wait for Cert Manager to be ready")
	certManagerSpinner.Start("installing Cert Manager")

	_, _ = cfmt.Println("\r  run command : kubectl wait --for=condition=available deployment/cert-manager-webhook -n cert-manager --timeout=180s -n operators")
	_, _ = cfmt.Println("\r  This might take a while. Time enough for a joke:")
	//for i := 0; i < 4; i++ {
	//    tellAChucknorrisJoke()
	//    time.Sleep(15 * time.Second)
	//}
	_, certManagerWaitErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/cert-manager-webhook", "-n", "cert-manager", "--timeout=180s", "-n", "operators").Output()
	if certManagerWaitErr != nil {
		certManagerSpinner.Error("Failed to run command. Try running it manually: kubectl wait --for=condition=available deployment/cert-manager-webhook -n cert-manager --timeout=180s -n operators")
		log.Fatal(certManagerWaitErr)
	}
	certManagerSpinner.Success("Cert Manager is ready")

	installCertManagerClusterIssuer("default")
}
