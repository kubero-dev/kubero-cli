/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"

	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/leaanthony/spinner"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install all required components for kubero",
	Long: `This command will try to install all required components for kubero on a kubernetes cluster.
It is allso possible to install a local kind cluster.

required binaries:
 - kubectl
 - kind (optional)`,
	Run: func(cmd *cobra.Command, args []string) {

		rand.Seed(time.Now().UnixNano())
		checkAllBinaries()
		installKind()
		checkCluster()
		installOLM()
		installIngress()
		installKuberoOperator()
		installKuberoUi()
		finalMessage()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func checkAllBinaries() {
	cfmt.Println("{{\n  Check for required binaries}}::lightWhite")
	if !checkBinary("kubectl") {
		cfmt.Println("{{✗ kubectl is not installed}}::red")
	} else {
		cfmt.Println("{{✓ kubectl is installed}}::lightGreen")
	}

	if !checkBinary("kind") {
		cfmt.Println("{{⚠ kind is not installed}}::yellow (only required if you want to install a local kind cluster)")
	} else {
		cfmt.Println("{{✓ kind is installed}}::lightGreen")
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
		log.Fatal("kind binary is not installed")
	}

	installer := resty.New()
	// TODO : installing the binaries needs to respect the OS and architecture
	//	installer.SetBaseURL("https://kind.sigs.k8s.io")
	//	installer.R().Get("/dl/v0.17.0/kind-linux-amd64")

	installer.SetBaseURL("https://raw.githubusercontent.com")
	kf, _ := installer.R().Get("/kubero-dev/kubero/main/kind.yaml")

	var kindConfig KindConfig
	yaml.Unmarshal(kf.Body(), &kindConfig)

	kindConfig.Name = promptLine("Kind Cluster Name", "", "kubero-"+strconv.Itoa(rand.Intn(1000)))
	kindConfig.Nodes[0].Image = "kindest/node:v1.25.3"
	kindConfig.Nodes[0].ExtraPortMappings[0].HostPort, _ = strconv.Atoi(promptLine("Local HTTP Port", "", "80"))
	kindConfig.Nodes[0].ExtraPortMappings[1].HostPort, _ = strconv.Atoi(promptLine("Local HTTPS Port", "", "443"))

	kindConfigYaml, _ := yaml.Marshal(kindConfig)
	fmt.Println("-------------- kind.yaml ---------------")
	fmt.Println(string(kindConfigYaml))
	fmt.Println("----------------------------------------")

	kindConfigErr := os.WriteFile("kind.yaml", kindConfigYaml, 0644)
	if kindConfigErr != nil {
		fmt.Println(kindConfigErr)
		return
	}

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
		cfmt.Println("{{✓ OLM is allredy installed}}::lightGreen")
		return
	}

	//namespace := promptLine("Install OLM in which namespace?", "[openshift-operator-lifecycle-manager,olm]", "olm")
	namespace := "olm"
	olmInstalled, _ := exec.Command("kubectl", "get", "deployment", "olm-operator", "-n", namespace).Output()
	if len(olmInstalled) > 0 {
		cfmt.Println("{{✓ OLM is allredy installed}}::lightGreen")
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
		cfmt.Println("{{✓ OLM CRD's allredy installed}}::lightGreen")
	} else {
		//fmt.Println("  run command : kubectl create -f " + olmURL + "/crds.yaml")
		_, olmCRDErr := exec.Command("kubectl", "create", "-f", olmURL+"/crds.yaml").Output()
		if olmCRDErr != nil {
			cfmt.Println("{{✗ OLM CRD installation failed }}::red")
			log.Fatal(olmCRDErr)
		} else {
			cfmt.Println("{{✓ OLM CRDs installed}}::lightGreen")
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

	ingressInstalled, _ := exec.Command("kubectl", "get", "ns", "ingress-nginx").Output()
	if len(ingressInstalled) > 0 {
		cfmt.Println("{{✓ Ingress is allredy installed}}::lightGreen")
		return
	}

	ingressInstall := promptLine("Install Ingress", "[y,n]", "y")
	if ingressInstall != "y" {
		log.Fatal("Ingress is required to install Kubero")
	} else {
		ingressProvider := promptLine("Provider", "[kind,aws,baremetal,cloud(Azure,Google,Oracle),do(digital ocean),exoscale,scw(scaleway)]", "kind")
		ingressSpinner := spinner.New("Install Ingress")
		ingressSpinner.Start("run command : kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/" + ingressProvider + "/deploy.yaml")
		_, ingressErr := exec.Command("kubectl", "apply", "-f", "https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/"+ingressProvider+"/deploy.yaml").Output()
		if ingressErr != nil {
			ingressSpinner.Error("Failed to run command. Try runnig it manually")
			log.Fatal(ingressErr)
		}
		ingressSpinner.Success("Ingress installed sucessfully")
	}

}

func installKuberoOperator() {

	cfmt.Println("{{\n  Install Kubero Operator}}::lightWhite")

	kuberoInstalled, _ := exec.Command("kubectl", "get", "operator", "kubero-operator.operators").Output()
	if len(kuberoInstalled) > 0 {
		cfmt.Println("{{✓ Kubero Operator is allredy installed}}::lightGreen")
		return
	}

	kuberoSpinner := spinner.New("Install Kubero Operator")
	kuberoSpinner.Start("run command : kubectl apply -f https://operatorhub.io/install/kubero-operator.yaml")
	_, kuberoErr := exec.Command("kubectl", "apply", "-f", "https://operatorhub.io/install/kubero-operator.yaml").Output()
	if kuberoErr != nil {
		fmt.Println("")
		kuberoSpinner.Error("Failed to run command to install the Operator. Try runnig it manually and then rerun the installation")
		log.Fatal(kuberoErr)
	}

	kuberoSpinner.UpdateMessage("Wait for Kubero Operator to be ready")
	var kuberoWait []byte
	for len(kuberoWait) == 0 {
		// kubectl api-resources --api-group=application.kubero.dev --no-headers=true
		kuberoWait, _ = exec.Command("kubectl", "api-resources", "--api-group=application.kubero.dev", "--no-headers=true").Output()
		time.Sleep(1 * time.Second)
	}

	kuberoSpinner.Success("Kubero Operator installed sucessfully")

}

func installKuberoUi() {

	ingressInstall := promptLine("Install Kubero UI", "[y,n]", "y")
	if ingressInstall != "y" {
		return
	}

	kuberoNSinstalled, _ := exec.Command("kubectl", "get", "ns", "kubero").Output()
	if len(kuberoNSinstalled) > 0 {
		cfmt.Println("{{✓ Kubero Namespace exists}}::lightGreen")
	} else {
		_, kuberoNSErr := exec.Command("kubectl", "create", "namespace", "kubero").Output()
		if kuberoNSErr != nil {
			fmt.Println("Failed to run command to create the namespace. Try runnig it manually")
			log.Fatal(kuberoNSErr)
		} else {
			cfmt.Println("{{✓ Kubero Namespace created}}::lightGreen")
		}
	}

	kuberoSecretInstalled, _ := exec.Command("kubectl", "get", "secret", "kubero-secrets", "-n", "kubero").Output()
	if len(kuberoSecretInstalled) > 0 {
		cfmt.Println("{{✓ Kubero Secret exists}}::lightGreen")
	} else {

		webhookSecret := promptLine("Random string for your webhook secret", "", generatePassword(20))

		sessionKey := promptLine("Random string for your session key", "", generatePassword(20))

		githubPersonalAccessToken := promptLine("Github personal access token (empty=disabled)", "", "")

		giteaPersonalAccessToken := promptLine("Gitea personal access token (empty=disabled)", "", "")

		adminUser := promptLine("Admin User", "", "admin")

		adminPass := promptLine("Admin Password", "", generatePassword(12))

		adminToken := promptLine("Random string for admin API token", "", generatePassword(20))

		var userDB []User
		userDB = append(userDB, User{Username: adminUser, Password: adminPass, Insecure: true, Apitoken: adminToken})
		userDBjson, _ := json.Marshal(userDB)
		userDBencoded := base64.StdEncoding.EncodeToString(userDBjson)

		_, kuberoErr := exec.Command("kubectl", "create", "secret", "generic", "kubero-secrets",
			"--from-literal=KUBERO_WEBHOOK_SECRET="+webhookSecret,
			"--from-literal=KUBERO_SESSION_KEY="+sessionKey,
			"--from-literal=KUBERO_GITHUB_PERSONAL_ACCESS_TOKEN="+githubPersonalAccessToken,
			"--from-literal=KUBERO_GITEA_PERSONAL_ACCESS_TOKEN="+giteaPersonalAccessToken,
			"--from-literal=KUBERO_USERS="+userDBencoded,
			"-n", "kubero",
		).Output()

		if kuberoErr != nil {
			cfmt.Println("{{✗ Failed to run command to create the secret. Try runnig it manually}}::red")
			log.Fatal(kuberoErr)
		} else {
			cfmt.Println("{{✓ Kubero Secret created}}::lightGreen")
		}
	}

	kuberoUIInstalled, _ := exec.Command("kubectl", "get", "kuberoes.application.kubero.dev", "-n", "kubero").Output()
	if len(kuberoUIInstalled) > 0 {
		cfmt.Println("{{✓ Kubero UI allready installed}}::lightGreen")
	} else {
		var outb, errb bytes.Buffer

		kuberoUI := exec.Command("kubectl", "apply", "-f", "https://raw.githubusercontent.com/kubero-dev/kubero-operator/main/config/samples/application_v1alpha1_kubero.yaml", "-n", "kubero")
		kuberoUI.Stdout = &outb
		kuberoUI.Stderr = &errb
		err := kuberoUI.Run()
		if err != nil {
			fmt.Println(errb.String())
			fmt.Println(outb.String())
			cfmt.Println("{{✗ Failed to run command to install Kubero UI. Try runnig it manually}}::red")
			log.Fatal()
		} else {
			cfmt.Println("{{✓ Kubero UI installed}}::lightGreen")
		}

		time.Sleep(1 * time.Second)
		kuberoUISpinner := spinner.New("Wait for Kubero UI to be ready")
		kuberoUISpinner.Start("run command : kubectl wait --for=condition=available deployment/kubero-sample -n kubero --timeout=60s")
		_, olmWaitErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/kubero-sample", "-n", "kubero", "--timeout=60s").Output()
		if olmWaitErr != nil {
			fmt.Println("") // keeps the spinner from overwriting the last line
			kuberoUISpinner.Error("Failed to run command. Try runnig it manually")
			log.Fatal(olmWaitErr)
		}
		kuberoUISpinner.Success("Kubero UI is ready")
	}

}

func finalMessage() {
	cfmt.Println(`

{{⚠ make sure your DNS is pointing to your Kubernetes cluster}}::yellow

	,--. ,--.        ,--.
	|  .'   /,--.,--.|  |-.  ,---. ,--.--. ,---.
	|  .   ' |  ||  || .-. '| .-. :|  .--'| .-. |
	|  |\   \'  ''  '| '-' |\   --.|  |   ' '-' '
	'--' '--' '----'  '---'  '----''--'    '---'

Your Kubero UI :
  {{http://kubero.lacolhost.com:80}}::blue

Documentation:
  https://github.com/kubero-dev/kubero/wiki
`)
}

func generatePassword(length int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!$?.-%")
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type User struct {
	ID       int    `json:"id"`
	Method   string `json:"method"`
	Username string `json:"username"`
	Password string `json:"password"`
	Insecure bool   `json:"insecure"`
	Apitoken string `json:"apitoken,omitempty"`
}

type KindConfig struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
	Name       string `yaml:"name"`
	Networking struct {
		IPFamily         string `yaml:"ipFamily"`
		APIServerAddress string `yaml:"apiServerAddress"`
	} `yaml:"networking"`
	Nodes []struct {
		Role                 string   `yaml:"role"`
		Image                string   `yaml:"image,omitempty"`
		KubeadmConfigPatches []string `yaml:"kubeadmConfigPatches"`
		ExtraPortMappings    []struct {
			ContainerPort int    `yaml:"containerPort"`
			HostPort      int    `yaml:"hostPort"`
			Protocol      string `yaml:"protocol"`
		} `yaml:"extraPortMappings"`
	} `yaml:"nodes"`
}
