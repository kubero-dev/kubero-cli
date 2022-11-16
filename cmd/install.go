/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/go-resty/resty/v2"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/leaanthony/spinner"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install all required components for kubero",
	Long: `This command will try to install all required components for kubero.

required binaries:
 - kubectl
 - kind (optional)`,
	Run: func(cmd *cobra.Command, args []string) {

		checkAllBinaries()

		installKind()
		checkCluster()
		installOLM()
		installIngress()
		/*

			ingress := promptLine("Install ingress-nginx", "[y,n]", "y")
			if ingress == "y" {
				installIngressNginx()
			}

			kuberoOperator := promptLine("Install kubero-operator", "[y,n]", "y")
			if kuberoOperator == "y" {
				installKuberoOperator()
			}

			kuberoUi := promptLine("Install kubero-ui", "[y,n]", "y")
			if kuberoUi == "y" {
				installKuberoUi()
			}
		*/
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func checkAllBinaries() {
	if !checkBinary("kubectl") {
		cfmt.Println("{{✗ kubectl is not installed}}::red")
	} else {
		cfmt.Println("{{✓ kubectl is installed}}::green")
	}

	if !checkBinary("kind") {
		cfmt.Println("{{⚠ kind is not installed}}::yellow (only required if you want to install a local kind cluster)")
	} else {
		cfmt.Println("{{✓ kind is installed}}::green")
	}
}

func checkBinary(binary string) bool {
	_, err := exec.LookPath(binary)
	return err == nil
}

func installKind() {
	kindInstall := promptLine("Start a local kubernetes kind cluster", "[y,n]", "n")
	if kindInstall != "y" {
		return
	}

	if !checkBinary("kind") {
		log.Fatal("kind is not installed")
	}

	installer := resty.New()
	// TODO : installing the binaries needs to respect the OS and architecture
	//	installer.SetBaseURL("https://kind.sigs.k8s.io")
	//	installer.R().Get("/dl/v0.17.0/kind-linux-amd64")

	installer.SetBaseURL("https://raw.githubusercontent.com")
	kindConfig, _ := installer.R().Get("/kubero-dev/kubero/main/kind.yaml")

	fmt.Println("-------------- kind.yaml ---------------")
	fmt.Println(kindConfig.String())
	fmt.Println("----------------------------------------")
	/* SKIPED for development
	kindConfigErr := os.WriteFile("kind.yaml", kindConfig.Body(), 0644)
	if kindConfigErr != nil {
		fmt.Println(kindConfigErr)
		return
	}
	*/

	kindSpinner := spinner.New("Spin up a local Kind cluster")
	kindSpinner.Start("run command : kind create cluster --config kind.yaml")
	out, err := exec.Command("kind", "create", "cluster", "--config", "kind.yaml").Output()
	if err != nil {
		kindSpinner.Error("Failed to run command. Try runnig it manually and skip this step")
		log.Fatal(err)
	}
	kindSpinner.Success("Kind cluster started sucessfully")

	fmt.Println(string(out))
}

func checkCluster() {
	var outb, errb bytes.Buffer

	clusterInfo := exec.Command("kubectl", "cluster-info")
	clusterInfo.Stdout = &outb
	clusterInfo.Stderr = &errb
	err := clusterInfo.Run()
	if err != nil {
		fmt.Println(errb.String())
		fmt.Println(outb.String())
		log.Fatal("command failed : kubectl cluster-info")
	}
	fmt.Println(outb.String())

	out, _ := exec.Command("kubectl", "config", "get-contexts").Output()
	fmt.Println(string(out))

	clusterselect := promptLine("Is the CURRENT cluster the one you wish to install Kubero?", "[y,n]", "y")
	if clusterselect == "n" {
		os.Exit(0)
	}
}

func installOLM() {

	openshiftInstalled, _ := exec.Command("kubectl", "get", "deployment", "olm-operator", "-n", "openshift-operator-lifecycle-manager").Output()
	if len(openshiftInstalled) > 0 {
		cfmt.Println("{{✓ OLM is allredy installed}}::green")
		return
	}

	//namespace := promptLine("Install OLM in which namespace?", "[openshift-operator-lifecycle-manager,olm]", "olm")
	namespace := "olm"
	olmInstalled, _ := exec.Command("kubectl", "get", "deployment", "olm-operator", "-n", namespace).Output()
	if len(olmInstalled) > 0 {
		cfmt.Println("{{✓ OLM is allredy installed}}::green")
		return
	}

	olmInstall := promptLine("Install OLM", "[y,n]", "y")
	if olmInstall != "y" {
		log.Fatal("OLM is required to install Kubero")
	}

	olmRelease := promptLine("Install OLM from which release?", "[0.19.0,0.20.0,0.21.0,0.22.0]", "0.22.0")
	olmURL := "https://github.com/operator-framework/operator-lifecycle-manager/releases/download/v" + olmRelease

	olmCRDInstalled, _ := exec.Command("kubectl", "get", "crd", "subscriptions.operators.coreos.com").Output()
	if len(olmCRDInstalled) > 0 {
		cfmt.Println("{{✓ OLM CRD's allredy installed}}::green")
	} else {
		fmt.Println("run command : kubectl create -f " + olmURL + "/crds.yaml")
		_, olmCRDErr := exec.Command("kubectl", "create", "-f", olmURL+"/crds.yaml").Output()
		if olmCRDErr != nil {
			cfmt.Println("{{✗ OLM CRD installation failed }}::red")
			log.Fatal(olmCRDErr)
		} else {
			cfmt.Println("{{✓ OLM CRDs installed}}::green")
		}
	}

	olmSpinner := spinner.New("Install OLM")
	olmSpinner.Start("run command : kubectl create -f " + olmURL + "/olm.yaml")

	_, olmOLMErr := exec.Command("kubectl", "create", "-f", olmURL+"/olm.yaml").Output()
	if olmOLMErr != nil {
		fmt.Println("")
		olmSpinner.Error("Failed to run command. Try runnig it manually")
		log.Fatal(olmOLMErr)
	}
	olmSpinner.Success("OLM installed sucessfully")

	olmWaitSpinner := spinner.New("Wait for OLM to be ready")
	olmWaitSpinner.Start("run command : kubectl wait --for=condition=available deployment/olm-operator -n " + namespace + " --timeout=60s")
	_, olmWaitErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/olm-operator", "-n", namespace, "--timeout=60s").Output()
	if olmWaitErr != nil {
		olmWaitSpinner.Error("Failed to run command. Try runnig it manually")
		log.Fatal(olmWaitErr)
	}
	olmWaitSpinner.Success("OLM is ready")

	olmWaitCatalogSpinner := spinner.New("Wait for OLM Catalog to be ready")
	olmWaitCatalogSpinner.Start("run command : kubectl wait --for=condition=available deployment/catalog-operator -n " + namespace + " --timeout=60s")
	_, olmWaitCatalogErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/catalog-operator", "-n", namespace, "--timeout=60s").Output()
	if olmWaitCatalogErr != nil {
		olmWaitCatalogSpinner.Error("Failed to run command. Try runnig it manually")
		log.Fatal(olmWaitCatalogErr)
	}
	olmWaitCatalogSpinner.Success("OLM Catalog is ready")
}

func installIngress() {

	ingressInstalled, _ := exec.Command("kubectl", "api-resources", "--api-group=networking.k8s.io", "--namespaced=false", "--no-headers=true").Output()
	if len(ingressInstalled) > 0 {
		cfmt.Println("{{✓ Ingress is allredy installed}}::green")
		return
	}

	ingressInstall := promptLine("Install Ingress", "[y,n]", "y")
	if ingressInstall != "y" {
		log.Fatal("Ingress is required to install Kubero")
	} else {
		ingressSpinner := spinner.New("Install Ingress")
		ingressSpinner.Start("run command : kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.41.2/deploy/static/provider/cloud/deploy.yaml")
		_, ingressErr := exec.Command("kubectl", "apply", "-f", "https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.41.2/deploy/static/provider/cloud/deploy.yaml").Output()
		if ingressErr != nil {
			ingressSpinner.Error("Failed to run command. Try runnig it manually")
			log.Fatal(ingressErr)
		}
		ingressSpinner.Success("Ingress installed sucessfully")
	}

}

/*
func installKuberoOperator() {
	fmt.Println("install kubero-operator")
}

func installKuberoUi() {
	fmt.Println("install kubero-ui")
}
*/
