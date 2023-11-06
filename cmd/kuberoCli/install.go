package kuberoCli

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/creasty/defaults"
	"k8s.io/utils/strings/slices"
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
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/tools/clientcmd"
)

var kuberoConfig KuberoConfig

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Create a Kubernetes cluster and install all required components for kubero",
	Long: `This command will create a kubernetes cluster and install all required components 
for kubero on any kubernetes cluster.

required binaries:
 - kubectl
 - kind (optional)`,
	Run: func(cmd *cobra.Command, args []string) {

		rand.Seed(time.Now().UnixNano())

		if err := defaults.Set(&kuberoConfig); err != nil {
			panic(err)
		}

		if arg_generate_config {
			writeInstallCLIConfig()
			return
		}

		checkAllBinaries()

		// Check config file
		if arg_config != "" {
			// Read yaml
			yamlFile, err := os.ReadFile(arg_config)
			if err != nil {
				log.Fatal(err)
			}
			yaml.Unmarshal(yamlFile, &kuberoConfig)
			kuberoConfig.configLoaded = true
			clusterType = kuberoConfig.ClusterType
			cfmt.Println("{{✓ Kubero install config loaded}}::lightGreen")
		}

		switch arg_component {
		case "metrics":
			installMetrics()
			return
		case "certmanager":
			installCertManager()
			return
		case "olm":
			installOLM()
			return
		case "kubero-operator":
			installKuberoOperator()
			return
		case "kubero-ui":
			installKuberoUi()
			return
		case "ingress":
			installIngress()
			return
		case "kubernetes":
			installKubernetes()
			checkCluster()
			return
		case "":
			printInstallSteps()
			installKubernetes()
			checkCluster()
			installOLM()
			installIngress()
			installMetrics()
			installCertManager()
			installKuberoOperator()
			installKuberoUi()
			writeCLIconfig()
			printDNSinfo()
			finalMessage()
			return
		default:
			return
		}
	},
}

var arg_adminPassword string
var arg_adminUser string
var arg_domain string
var arg_apiToken string
var arg_port string
var arg_portSecure string
var clusterType string
var arg_config string
var arg_generate_config bool
var arg_component string
var install_olm bool
var ingressControllerVersion = "v1.7.0" // https://github.com/kubernetes/ingress-nginx/tags -> controller-v1.5.1

// var clusterTypeSelection = "[scaleway,linode,gke,digitalocean,kind]"
var clusterTypeList = []string{"kind", "linode", "scaleway", "gke", "digitalocean"}

func init() {
	installCmd.Flags().StringVar(&arg_config, "config", "", "config file yaml to propmtless install kubero")
	installCmd.Flags().BoolVar(&arg_generate_config, "generate-config", false, "generate config file yaml to propmtless install kubero")
	installCmd.Flags().StringVarP(&arg_component, "component", "c", "", "install component (kubernetes,olm,ingress,metrics,certmanager,kubero-operator,kubero-ui)")
	installCmd.Flags().StringVarP(&arg_adminUser, "user", "u", "", "Admin username for the kubero UI")
	installCmd.Flags().StringVarP(&arg_adminPassword, "user-password", "U", "", "Password for the admin user")
	installCmd.Flags().StringVarP(&arg_apiToken, "apitoken", "a", "", "API token for the admin user")
	installCmd.Flags().StringVarP(&arg_port, "port", "p", "", "Kubero UI HTTP port")
	installCmd.Flags().StringVarP(&arg_portSecure, "secureport", "P", "", "Kubero UI HTTPS port")
	installCmd.Flags().StringVarP(&arg_domain, "domain", "d", "", "Domain name for the kubero UI")
	rootCmd.AddCommand(installCmd)

	install_olm = false
}

func writeInstallCLIConfig() {
	if arg_generate_config {
		// Write config file
		kuberoConfig.KuberoUI.WebhookSecret = generateRandomString(20, "")
		kuberoConfig.KuberoUI.SessionKey = generateRandomString(20, "")
		kuberoConfig.KuberoUI.Registry.Password = generateRandomString(20, "")
		kuberoConfig.KuberoUI.AdminPassword = generateRandomString(12, "")
		kuberoConfig.KuberoUI.ApiToken = generateRandomString(20, "")
		kuberoConfigYaml, _ := yaml.Marshal(kuberoConfig)
		kuberoConfigErr := os.WriteFile("kuberoInstallConfigSample.yaml", kuberoConfigYaml, 0644)
		if kuberoConfigErr != nil {
			fmt.Println(kuberoConfigErr)
			return
		}
		cfmt.Println("{{✓ Kubero install config file generated to kuberoInstallConfigSample.yaml - make a copy to avoid being overwritten}}::lightGreen")
	}
}

func checkAllBinaries() {
	cfmt.Println("\n  Check for required binaries")
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

	if !checkBinary("gcloud") {
		cfmt.Println("{{⚠ gcloud is not installed}}::yellow (only required if you want to install a GKE cluster)")
	} else {
		cfmt.Println("{{✓ gcloud is installed}}::lightGreen")
	}
}

func printInstallSteps() {

	cfmt.Print(`
  Steps to install kubero:
    1. Create a kubernetes cluster {{(optional)}}::gray
    2. Install the OLM {{(optional)}}::gray
    3. Install the ingress controller {{(required)}}::gray
    4. Install the metrics server {{(optional, but recommended)}}::gray
    5. Install the cert-manager {{(optional)}}::gray
    6. Install the kubero operator {{(required)}}::gray
    7. Install the kubero UI {{(optional, but highly recommended)}}::gray
    8. Write the kubero CLI config
`)
}

func checkBinary(binary string) bool {
	_, err := exec.LookPath(binary)
	return err == nil
}

func installKubernetes() {
	if kuberoConfig.configLoaded {
		if !kuberoConfig.Install.Kubernetes {
			return
		}
	} else {
		kubernetesInstall := promptLine("1) Create a kubernetes cluster", "[y,n]", "y")
		if kubernetesInstall != "y" {
			return
		}

		clusterType = selectFromList("Select a Kubernetes provider", clusterTypeList, "")
	}

	switch clusterType {
	case "scaleway":
		installScaleway()
	case "linode":
		installLinode()
	case "gke":
		installGKE()
	case "digitalocean":
		installDigitalOcean()
	case "kind":
		installKind()
	default:
		cfmt.Println("{{✗ Unknown cluster type}}::red")
		os.Exit(1)
	}

}

func tellAChucknorrisJoke() {

	jokesapi := resty.New().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "kubero-cli/"+kuberoCliVersion).
		SetBaseURL("https://api.chucknorris.io/jokes/random")

	joke, _ := jokesapi.R().Get("?category=dev")
	var jokeResponse JokeResponse
	json.Unmarshal(joke.Body(), &jokeResponse)
	cfmt.Println("\r{{  " + jokeResponse.Value + "       }}::gray")
}

func mergeKubeconfig(kubeconfig []byte) error {

	new := clientcmd.NewDefaultPathOptions()
	config1, _ := new.GetStartingConfig()
	config2, err := clientcmd.Load(kubeconfig)
	if err != nil {
		return err
	}
	// append the second config to the first
	for k, v := range config2.Clusters {
		config1.Clusters[k] = v
	}
	for k, v := range config2.AuthInfos {
		config1.AuthInfos[k] = v
	}
	for k, v := range config2.Contexts {
		config1.Contexts[k] = v
	}

	config1.CurrentContext = config2.CurrentContext

	clientcmd.ModifyConfig(clientcmd.DefaultClientConfig.ConfigAccess(), *config1, true)
	return nil
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

	if kuberoConfig.configLoaded {
		return
	}

	clusterselect := promptLine("Is the CURRENT cluster the one you wish to install Kubero?", "[y,n]", "y")
	if clusterselect == "n" {
		os.Exit(0)
	}
}

func installOLM() {

	openshiftInstalled, _ := exec.Command("kubectl", "get", "deployment", "olm-operator", "-n", "openshift-operator-lifecycle-manager").Output()
	if len(openshiftInstalled) > 0 {
		cfmt.Println("{{✓ OLM is already installed}}::lightGreen")
		return
	}

	//namespace := promptLine("Install OLM in which namespace?", "[openshift-operator-lifecycle-manager,olm]", "olm")
	namespace := "olm"
	olmInstalled, _ := exec.Command("kubectl", "get", "deployment", "olm-operator", "-n", namespace).Output()
	if len(olmInstalled) > 0 {
		cfmt.Println("{{✓ OLM is already installed}}::lightGreen")
		return
	}

	olmRelease := "0.23.1"

	if kuberoConfig.configLoaded {
		install_olm = kuberoConfig.Install.OLM
		if !install_olm {
			return
		}

		if kuberoConfig.OLMRelease != "" {
			olmRelease = kuberoConfig.OLMRelease
		}
	} else {
		olmInstall := promptLine("2) Install OLM", "[y,n]", "n")
		if olmInstall != "y" {
			install_olm = false
			return
		} else {
			install_olm = true
		}

		olmRelease = promptLine("Install OLM from which release?", "[0.20.0,0.21.0,0.22.0,0.23.1]", "0.23.1")
	}

	olmURL := "https://github.com/operator-framework/operator-lifecycle-manager/releases/download/v" + olmRelease

	olmSpinner := spinner.New("Install OLM")

	olmCRDInstalled, _ := exec.Command("kubectl", "get", "crd", "subscriptions.operators.coreos.com").Output()
	if len(olmCRDInstalled) > 0 {
		cfmt.Println("{{✓ OLM CRD's already installed}}::lightGreen")
	} else {
		olmSpinner.Start("run command : kubectl create -f " + olmURL + "/crds.yaml")
		_, olmCRDErr := exec.Command("kubectl", "create", "-f", olmURL+"/crds.yaml").Output()
		if olmCRDErr != nil {
			fmt.Println("")
			olmSpinner.Error("OLM CRD installation failed. Try runnig this command manually: kubectl create -f " + olmURL + "/crds.yaml")
			log.Fatal(olmCRDErr)
		} else {
			olmSpinner.Success("OLM CRDs installed sucessfully")
		}
	}

	olmSpinner.Start("run command : kubectl create -f " + olmURL + "/olm.yaml")

	_, olmOLMErr := exec.Command("kubectl", "create", "-f", olmURL+"/olm.yaml").Output()
	if olmOLMErr != nil {
		fmt.Println("")
		olmSpinner.Error("Failed to run command. Try runnig this command manually: kubectl create -f " + olmURL + "/olm.yaml")
		log.Fatal(olmOLMErr)
	}
	olmSpinner.Success("OLM installed sucessfully")

	olmWaitSpinner := spinner.New("Wait for OLM to be ready")
	olmWaitSpinner.Start("run command : kubectl wait --for=condition=available deployment/olm-operator -n " + namespace + " --timeout=180s")
	_, olmWaitErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/olm-operator", "-n", namespace, "--timeout=180s").Output()
	if olmWaitErr != nil {
		olmWaitSpinner.Error("Failed to run command. Try runnig this command manually: kubectl wait --for=condition=available deployment/olm-operator -n " + namespace + " --timeout=180s")
		log.Fatal(olmWaitErr)
	}
	olmWaitSpinner.Success("OLM is ready")

	olmWaitCatalogSpinner := spinner.New("Wait for OLM Catalog to be ready")
	olmWaitCatalogSpinner.Start("run command : kubectl wait --for=condition=available deployment/catalog-operator -n " + namespace + " --timeout=180s")
	_, olmWaitCatalogErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/catalog-operator", "-n", namespace, "--timeout=180s").Output()
	if olmWaitCatalogErr != nil {
		olmWaitCatalogSpinner.Error("Failed to run command. Try runnig this command manually: kubectl wait --for=condition=available deployment/catalog-operator -n " + namespace + " --timeout=180s")
		log.Fatal(olmWaitCatalogErr)
	}
	olmWaitCatalogSpinner.Success("OLM Catalog is ready")
}

func installMetrics() {

	installed, _ := exec.Command("kubectl", "get", "deployments.apps", "metrics-server", "-n", "kube-system").Output()
	if len(installed) > 0 {
		cfmt.Println("{{✓ Metrics is already enabled}}::lightGreen")
		return
	}
	if kuberoConfig.configLoaded {
		if !kuberoConfig.Install.Metrics {
			return
		}
	} else {
		install := promptLine("4) Install Kubernetes internal metrics service (required for HPA and stats)", "[y,n]", "y")
		if install != "y" {
			return
		}
	}

	//components := "https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml"
	components := "https://raw.githubusercontent.com/kubero-dev/kubero-operator/main/deploy/metrics-server.yaml"
	_, installErr := exec.Command("kubectl", "apply", "-f", components).Output()

	if installErr != nil {
		fmt.Println("failed to install metrics server")
		log.Fatal(installErr)
	}
	cfmt.Println("{{✓ Metrics server installed}}::lightGreen")
}

func installIngress() {
	// TODO: is this check correct? The ingress-nginx yaml contains this selector by default which blocks the install when no nodes have this label
	ingressReadyLabel, _ := exec.Command("kubectl", "get", "nodes", "--selector=ingress-ready=true", "-o", "jsonpath='{.items[*].metadata.name}'").Output()
	if len(ingressReadyLabel) == 0 {
		cfmt.Println("{{✗ Ingress: no nodes with label ingress-ready=true found}}::red")
	}

	ingressInstalled, _ := exec.Command("kubectl", "get", "ns", "ingress-nginx").Output()
	if len(ingressInstalled) > 0 {
		cfmt.Println("{{✓ Ingress is already installed}}::lightGreen")
		return
	}

	ingressInstall := true
	if kuberoConfig.configLoaded {
		ingressInstall = kuberoConfig.Install.Ingress
	} else {
		ingressInstall = promptLine("3) Install Ingress", "[y,n]", "y") == "y"
	}
	if !ingressInstall || len(ingressReadyLabel) == 0 {
		return
	} else {

		if clusterType == "" {
			clusterType = selectFromList("Which cluster type have you insalled?", clusterTypeList, "")
		}

		prefill := "baremetal"
		switch clusterType {
		case "kind":
			prefill = "kind"
		case "linode":
			prefill = "cloud"
		case "gke":
			prefill = "cloud"
		case "scaleway":
			prefill = "scw"
		case "digitalocean":
			prefill = "do"
		}
		ingressProvider := prefill
		if kuberoConfig.configLoaded {
			ingressProvider = kuberoConfig.IngressProvider
		} else {
			ingressProviderList := []string{"kind", "aws", "baremetal", "cloud", "do", "exoscale", "scw"}
			ingressProvider = selectFromList("Provider [kind, aws, baremetal, cloud(Azure,Google,Oracle,Linode), do(digital ocean), exoscale, scw(scaleway)]", ingressProviderList, prefill)
		}

		ingressSpinner := spinner.New("Install Ingress")
		URL := "https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-" + ingressControllerVersion + "/deploy/static/provider/" + ingressProvider + "/deploy.yaml"
		ingressSpinner.Start("run command : kubectl apply -f " + URL)
		_, ingressErr := exec.Command("kubectl", "apply", "-f", URL).Output()
		if ingressErr != nil {
			ingressSpinner.Error("Failed to run command. Try runnig this command manually: kubectl apply -f " + URL)
			log.Fatal(ingressErr)
		}
		ingressSpinner.Success("Ingress installed sucessfully")
	}

}

func installKuberoOperator() {

	cfmt.Println("\n  {{6) Install Kubero Operator}}::bold")

	kuberoInstalled, _ := exec.Command("kubectl", "get", "operator", "kubero-operator.operators").Output()
	if len(kuberoInstalled) > 0 {
		cfmt.Println("{{✓ Kubero Operator is already installed}}::lightGreen")
		return
	}

	if install_olm {
		installKuberoOLMOperator()
	} else {
		installKuberoOperatorSlim()
	}
}

func installKuberoOLMOperator() {

	kuberoSpinner := spinner.New("Install Kubero Operator")
	kuberoSpinner.Start("run command : kubectl apply -f https://operatorhub.io/install/kubero-operator.yaml")
	_, kuberoErr := exec.Command("kubectl", "apply", "-f", "https://operatorhub.io/install/kubero-operator.yaml").Output()
	if kuberoErr != nil {
		fmt.Println("")
		kuberoSpinner.Error("Failed to run command to install the Operator. Try runnig this command manually: kubectl apply -f https://operatorhub.io/install/kubero-operator.yaml")
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

func installKuberoOperatorSlim() {

	kuberoSpinner := spinner.New("Install Kubero Operator")
	kuberoSpinner.Start("run command : kubectl apply -f https://raw.githubusercontent.com/kubero-dev/kubero-operator/main/deploy/operator.yaml")
	_, kuberoErr := exec.Command("kubectl", "apply", "-f", "https://raw.githubusercontent.com/kubero-dev/kubero-operator/main/deploy/operator.yaml").Output()
	if kuberoErr != nil {
		fmt.Println("")
		kuberoSpinner.Error("Failed to run command to install the Operator. Try runnig this command manually: kubectl apply -f https://raw.githubusercontent.com/kubero-dev/kubero-operator/main/deploy/operator.yaml")
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
	kuberoUIInstall := true
	if kuberoConfig.configLoaded {
		kuberoUIInstall = kuberoConfig.Install.KuberoUI
	} else {
		kuberoUIInstall = promptLine("7) Install Kubero UI", "[y,n]", "y") == "y"
	}
	if !kuberoUIInstall {
		return
	}

	kuberoNSinstalled, _ := exec.Command("kubectl", "get", "ns", "kubero").Output()
	if len(kuberoNSinstalled) > 0 {
		cfmt.Println("{{✓ Kubero Namespace exists}}::lightGreen")
	} else {
		_, kuberoNSErr := exec.Command("kubectl", "create", "namespace", "kubero").Output()
		if kuberoNSErr != nil {
			fmt.Println("Failed to run command to create the namespace. Try runnig this command manually: kubectl create namespace kubero")
			log.Fatal(kuberoNSErr)
		} else {
			cfmt.Println("{{✓ Kubero Namespace created}}::lightGreen")
		}
	}

	kuberoSecretInstalled, _ := exec.Command("kubectl", "get", "secret", "kubero-secrets", "-n", "kubero").Output()
	if len(kuberoSecretInstalled) > 0 {
		cfmt.Println("{{✓ Kubero Secret exists}}::lightGreen")
	} else {
		webhookSecret := generateRandomString(20, "")
		sessionKey := generateRandomString(20, "")

		if kuberoConfig.configLoaded {
			webhookSecret = kuberoConfig.KuberoUI.WebhookSecret
			sessionKey = kuberoConfig.KuberoUI.SessionKey
			arg_adminUser = kuberoConfig.KuberoUI.AdminUser
			arg_adminPassword = kuberoConfig.KuberoUI.AdminPassword
			arg_apiToken = kuberoConfig.KuberoUI.ApiToken
		} else {
			webhookSecret = promptLine("Random string for your webhook secret", "", webhookSecret)
			sessionKey = promptLine("Random string for your session key", "", sessionKey)

			if arg_adminUser == "" {
				arg_adminUser = promptLine("Admin User", "", "admin")
			}

			if arg_adminPassword == "" {
				arg_adminPassword = promptLine("Admin Password", "", generateRandomString(12, ""))
			}

			if arg_apiToken == "" {
				arg_apiToken = promptLine("Random string for admin API token", "", generateRandomString(20, ""))
			}
		}

		var userDB []User
		userDB = append(userDB, User{Username: arg_adminUser, Password: arg_adminPassword, Insecure: true, Apitoken: arg_apiToken})
		userDBjson, _ := json.Marshal(userDB)
		userDBencoded := base64.StdEncoding.EncodeToString(userDBjson)

		createSecretCommand := exec.Command("kubectl", "create", "secret", "generic", "kubero-secrets",
			"--from-literal=KUBERO_WEBHOOK_SECRET="+webhookSecret,
			"--from-literal=KUBERO_SESSION_KEY="+sessionKey,
			"--from-literal=KUBERO_USERS="+userDBencoded,
		)

		if kuberoConfig.configLoaded {
			if kuberoConfig.KuberoUI.Github.Enabled {
				createSecretCommand.Args = append(createSecretCommand.Args, "--from-literal=GITHUB_PERSONAL_ACCESS_TOKEN="+kuberoConfig.KuberoUI.Github.PersonalAccessToken)
			}
			if kuberoConfig.KuberoUI.Gitea.Enabled {
				createSecretCommand.Args = append(createSecretCommand.Args, "--from-literal=GITEA_PERSONAL_ACCESS_TOKEN="+kuberoConfig.KuberoUI.Gitea.PersonalAccessToken)
				createSecretCommand.Args = append(createSecretCommand.Args, "--from-literal=GITEA_BASEURL="+kuberoConfig.KuberoUI.Gitea.BaseURL)
			}
			if kuberoConfig.KuberoUI.Gogs.Enabled {
				createSecretCommand.Args = append(createSecretCommand.Args, "--from-literal=GOGS_PERSONAL_ACCESS_TOKEN="+kuberoConfig.KuberoUI.Gogs.PersonalAccessToken)
				createSecretCommand.Args = append(createSecretCommand.Args, "--from-literal=GOGS_BASEURL="+kuberoConfig.KuberoUI.Gogs.BaseURL)
			}
			if kuberoConfig.KuberoUI.Gitlab.Enabled {
				createSecretCommand.Args = append(createSecretCommand.Args, "--from-literal=GITLAB_PERSONAL_ACCESS_TOKEN="+kuberoConfig.KuberoUI.Gitlab.PersonalAccessToken)
				createSecretCommand.Args = append(createSecretCommand.Args, "--from-literal=GITLAB_BASEURL="+kuberoConfig.KuberoUI.Gitlab.BaseURL)
			}
			if kuberoConfig.KuberoUI.Bitbucket.Enabled {
				createSecretCommand.Args = append(createSecretCommand.Args, "--from-literal=BITBUCKET_USERNAME="+kuberoConfig.KuberoUI.Bitbucket.Username)
				createSecretCommand.Args = append(createSecretCommand.Args, "--from-literal=BITBUCKET_APP_PASSWORD="+kuberoConfig.KuberoUI.Bitbucket.AppPassword)
			}
		} else {
			githubConfigure := promptLine("Configure Github", "[y,n]", "y")
			githubPersonalAccessToken := ""
			if githubConfigure == "y" {
				githubPersonalAccessToken = promptLine("Github personal access token", "", "")
				createSecretCommand.Args = append(createSecretCommand.Args, "--from-literal=GITHUB_PERSONAL_ACCESS_TOKEN="+githubPersonalAccessToken)
			}

			giteaConfigure := promptLine("Configure Gitea", "[y,n]", "n")
			giteaPersonalAccessToken := ""
			giteaBaseUrl := ""
			if giteaConfigure == "y" {
				giteaPersonalAccessToken = promptLine("Gitea personal access token", "", "")
				giteaBaseUrl = promptLine("Gitea URL", "http://localhost:3000", "")
				createSecretCommand.Args = append(createSecretCommand.Args, "--from-literal=GITEA_PERSONAL_ACCESS_TOKEN="+giteaPersonalAccessToken)
				createSecretCommand.Args = append(createSecretCommand.Args, "--from-literal=GITEA_BASEURL="+giteaBaseUrl)
			}

			gogsConfigure := promptLine("Configure Gogs", "[y,n]", "n")
			gogsPersonalAccessToken := ""
			gogsBaseUrl := ""
			if gogsConfigure == "y" {
				gogsPersonalAccessToken = promptLine("Gogs personal access token", "", "")
				gogsBaseUrl = promptLine("Gogs URL", "http://localhost:3000", "")
				createSecretCommand.Args = append(createSecretCommand.Args, "--from-literal=GOGS_PERSONAL_ACCESS_TOKEN="+gogsPersonalAccessToken)
				createSecretCommand.Args = append(createSecretCommand.Args, "--from-literal=GOGS_BASEURL="+gogsBaseUrl)
			}

			gitlabConfigure := promptLine("Configure Gitlab", "[y,n]", "n")
			gitlabPersonalAccessToken := ""
			gitlabBaseUrl := ""
			if gitlabConfigure == "y" {
				gitlabPersonalAccessToken = promptLine("Gitlab personal access token", "", "")
				gitlabBaseUrl = promptLine("Gitlab URL", "http://localhost:3080", "")
				createSecretCommand.Args = append(createSecretCommand.Args, "--from-literal=GITLAB_PERSONAL_ACCESS_TOKEN="+gitlabPersonalAccessToken)
				createSecretCommand.Args = append(createSecretCommand.Args, "--from-literal=GITLAB_BASEURL="+gitlabBaseUrl)
			}

			bitbucketConfigure := promptLine("Configure Bitbucket", "[y,n]", "n")
			bitbucketUsername := ""
			bitbucketAppPassword := ""
			if bitbucketConfigure == "y" {
				bitbucketUsername = promptLine("Bitbucket Username", "", "")
				bitbucketAppPassword = promptLine("Bitbucket App Password", "", "")
				createSecretCommand.Args = append(createSecretCommand.Args, "--from-literal=BITBUCKET_USERNAME="+bitbucketUsername)
				createSecretCommand.Args = append(createSecretCommand.Args, "--from-literal=BITBUCKET_APP_PASSWORD="+bitbucketAppPassword)
			}
		}

		createSecretCommand.Args = append(createSecretCommand.Args, "-n", "kubero")

		_, kuberoErr := createSecretCommand.Output()

		if kuberoErr != nil {
			cfmt.Println("{{✗ Failed to run command to create the secrets.}}::red")
			log.Fatal(kuberoErr)
		} else {
			cfmt.Println("{{✓ Kubero Secret created}}::lightGreen")
		}
	}

	kuberoUIInstalled, _ := exec.Command("kubectl", "get", "kuberoes.application.kubero.dev", "-n", "kubero").Output()
	if len(kuberoUIInstalled) > 0 {
		cfmt.Println("{{✓ Kubero UI allready installed}}::lightGreen")
	} else {
		installer := resty.New()

		installer.SetBaseURL("https://raw.githubusercontent.com")
		kf, _ := installer.R().Get("kubero-dev/kubero-operator/main/config/samples/application_v1alpha1_kubero.yaml")

		var kuberoUIConfig KuberoUIConfig
		yaml.Unmarshal(kf.Body(), &kuberoUIConfig)

		webhookURL := "https://" + arg_domain + "/api/repo/webhooks"
		kuberoUIssl := true

		if kuberoConfig.configLoaded {
			arg_domain = kuberoConfig.KuberoUI.UI.Host
			webhookURL = kuberoConfig.KuberoUI.UI.WebhookUrl
			kuberoUIssl = kuberoConfig.KuberoUI.UI.SSL
		} else {
			if arg_domain == "" {
				arg_domain = promptLine("Kuberi UI Domain", "", "kubero.localhost.com")
			}

			webhookURL = promptLine("URL to which the webhooks should be sent", "", "https://"+arg_domain+"/api/repo/webhooks")
			kuberoUIssl = promptLine("Enable SSL for the Kubero UI", "[y/n]", "y") == "y"
		}

		kuberoUIConfig.Spec.Ingress.Hosts[0].Host = arg_domain
		kuberoUIConfig.Spec.Kubero.WebhookURL = webhookURL

		if kuberoUIssl {
			kuberoUIConfig.Spec.Ingress.Annotations.KubernetesIoIngressClass = "letsencrypt-prod"
			kuberoUIConfig.Spec.Ingress.Annotations.KubernetesIoTLSacme = "true"

			kuberoUIConfig.Spec.Ingress.TLS = []KuberoUItls{
				{
					Hosts:      []string{arg_domain},
					SecretName: "kubero-tls",
				},
			}
		}

		kuberoUIRegistry := false
		if kuberoConfig.configLoaded {
			kuberoUIRegistry = kuberoConfig.KuberoUI.Registry.Enabled
		} else {
			kuberoUIRegistry = promptLine("Enable Buildpipeline for Kubero (BETA)", "[y/n]", "n") == "y"
		}
		if kuberoUIRegistry {
			kuberoUIConfig.Spec.Registry.Enabled = true

			kuberoUICreateRegistry := false
			kuberoUIRegistryPort := "443"
			kuberoUIRegistryHost := "registry." + arg_domain
			kuberoUIRegistryUsername := "admin"
			kuberoUIRegistryPassword := generateRandomString(12, "")
			kuberoUIRegistryStorage := "10Gi"

			storageClassList := getAvailableStorageClasses()
			kuberoUIRegistryStorageClassName := ""

			if kuberoConfig.configLoaded {
				kuberoUICreateRegistry = kuberoConfig.KuberoUI.Registry.Local
				kuberoUIRegistryPort = kuberoConfig.KuberoUI.Registry.Port
				if kuberoConfig.KuberoUI.Registry.Host != "" {
					kuberoUIRegistryHost = kuberoConfig.KuberoUI.Registry.Host
				}
				kuberoUIRegistryUsername = kuberoConfig.KuberoUI.Registry.Username
				kuberoUIRegistryPassword = kuberoConfig.KuberoUI.Registry.Password
				kuberoUIRegistryStorage = kuberoConfig.KuberoUI.Registry.StorageSize
				kuberoUIRegistryStorageClassName = kuberoConfig.KuberoUI.Registry.StorageClassName

				if !slices.Contains(storageClassList, kuberoUIRegistryStorageClassName) {
					cfmt.Println("{{✗ The storage class " + kuberoUIRegistryStorageClassName + " is not available.}}::red")
					return
				}
			} else {
				kuberoUICreateRegistry = promptLine("Create a local Registry for Kubero", "[y/n]", "n") == "y"
				kuberoUIRegistryPort = promptLine("Registry port", "", "443")
				kuberoUIRegistryHost = promptLine("Registry domain", "", "registry."+arg_domain)
				kuberoUIRegistryUsername = promptLine("Registry username", "", "admin")
				kuberoUIRegistryPassword = promptLine("Registry password", "", generateRandomString(12, ""))
				kuberoUIRegistryStorage = promptLine("Registry storage size", "", "10Gi")
				kuberoUIRegistryStorageClassName = selectFromList("Registry storage class", storageClassList, "")
			}

			if kuberoUICreateRegistry {
				kuberoUIConfig.Spec.Registry.Create = true
			}

			kuberoUIConfig.Spec.Registry.Port, _ = strconv.Atoi(kuberoUIRegistryPort)
			kuberoUIConfig.Spec.Registry.Host = kuberoUIRegistryHost
			kuberoUIConfig.Spec.Registry.Account.Username = kuberoUIRegistryUsername
			kuberoUIConfig.Spec.Registry.Account.Password = kuberoUIRegistryPassword

			kuberoUIRegistryPasswordBytes, _ := bcrypt.GenerateFromPassword([]byte(kuberoUIRegistryPassword), 14)
			kuberoUIConfig.Spec.Registry.Account.Hash = string(kuberoUIRegistryPasswordBytes)

			kuberoUIConfig.Spec.Registry.Storage = kuberoUIRegistryStorage
			kuberoUIConfig.Spec.Registry.StorageClassName = kuberoUIRegistryStorageClassName
		}

		if clusterType == "" && !kuberoConfig.configLoaded {
			clusterType = selectFromList("Which cluster type have you insalled?", clusterTypeList, "")
		}

		if clusterType == "linode" ||
			clusterType == "digitalocean" ||
			clusterType == "scaleway" ||
			clusterType == "gke" {
			kuberoUIConfig.Spec.Ingress.ClassName = "nginx"
		}

		kuberiUIYaml, _ := yaml.Marshal(kuberoUIConfig)
		kuberiUIErr := os.WriteFile("kuberoUI.yaml", kuberiUIYaml, 0644)
		if kuberiUIErr != nil {
			fmt.Println(kuberiUIErr)
			return
		}

		_, olminstallErr := exec.Command("kubectl", "apply", "-f", "kuberoUI.yaml", "-n", "kubero").Output()
		if olminstallErr != nil {
			cfmt.Println("{{✗ Failed to run command to install Kubero UI. Try runnig this command manually: kubectl apply -f kuberoUI.yaml -n kubero}}::red")
			return
		} else {
			e := os.Remove("kuberoUI.yaml")
			if e != nil {
				log.Fatal(e)
			}
			cfmt.Println("{{✓ Kubero UI installed}}::lightGreen")
		}

		kuberoUISpinner := spinner.New("Wait for Kubero UI to be created")
		kuberoUISpinner.Start("Wait for Kubero UI to be created")

		var kuberoWait []byte
		for len(kuberoWait) == 0 {
			// kubectl get --ignore-not-found deployment kubero
			kuberoWait, _ = exec.Command("kubectl", "get", "--ignore-not-found", "deployment", "kubero", "-n", "kubero").Output()
			kuberoUISpinner.UpdateMessage("Waiting for Kubero UI to be created")
			time.Sleep(1 * time.Second)
		}

		kuberoUISpinner.UpdateMessage("Waiting for Kubero UI to be ready")

		time.Sleep(5 * time.Second)
		// kubectl wait --for=condition=available deployment/kubero -n kubero --timeout=180s
		_, olmWaitErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/kubero", "-n", "kubero", "--timeout=180s").Output()
		if olmWaitErr != nil {
			kuberoUISpinner.Error("Failed to wait for Kubero UI to become ready")
			log.Fatal(olmWaitErr)
		}
		kuberoUISpinner.Success("Kubero UI is ready")
	}

}

func installCertManager() {
	install := true
	if kuberoConfig.configLoaded {
		install = kuberoConfig.Install.CertManager
	} else {
		install = promptLine("5) Install SSL Certmanager", "[y,n]", "y") == "y"
	}
	if !install {
		return
	}

	if install_olm {
		installOLMCertManager()
	} else {
		installCertManagerSlim()
	}
}

func installCertManagerSlim() {

	kuberoUIInstalled, _ := exec.Command("kubectl", "get", "crd", "certificates.cert-manager.io").Output()
	if len(kuberoUIInstalled) > 0 {
		cfmt.Println("{{✓ Certmanager already installed}}::lightGreen")
		return
	}

	certManagerSpinner := spinner.New("Install Cert Manager")
	//certmanagerUrl := "https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml"
	certmanagerUrl := "https://github.com/cert-manager/cert-manager/releases/latest/download/cert-manager.yaml"
	certManagerSpinner.Start("run command : kubectl create -f " + certmanagerUrl)
	_, certManagerErr := exec.Command("kubectl", "create", "-f", certmanagerUrl).Output()
	if certManagerErr != nil {
		certManagerSpinner.Error("Failed to run command. Try runnig this command manually: kubectl create -f " + certmanagerUrl)
		log.Fatal(certManagerErr)
	}

	certManagerSpinner.UpdateMessage("Waiting for Cert Manager to be ready")
	time.Sleep(5 * time.Second)
	_, certManagerWaitErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/cert-manager-webhook", "-n", "cert-manager", "--timeout=180s", "-n", "cert-manager").Output()
	if certManagerWaitErr != nil {
		certManagerSpinner.Error("Failed to run command. Try runnig it manually: kubectl wait --for=condition=available deployment/cert-manager-webhook -n cert-manager --timeout=180s -n cert-manager")
		log.Fatal(certManagerWaitErr)
	}
	certManagerSpinner.Success("Cert Manager installed")

	installCertManagerClusterissuer()

}

func installCertManagerClusterissuer() {

	installer := resty.New()

	installer.SetBaseURL("https://raw.githubusercontent.com")
	kf, _ := installer.R().Get("kubero-dev/kubero-cli/main/templates/certmanagerClusterIssuer.prod.yaml")

	var certmanagerClusterIssuer CertmanagerClusterIssuer
	yaml.Unmarshal(kf.Body(), &certmanagerClusterIssuer)

	if kuberoConfig.configLoaded {
		certmanagerClusterIssuer.Spec.Acme.Email = kuberoConfig.CertManager.AcmeEmail
	} else {
		certmanagerClusterIssuer.Spec.Acme.Email = promptLine("Letsencrypt ACME contact email", "", "noreply@yourdomain.com")
	}

	certmanagerClusterIssuerYaml, _ := yaml.Marshal(certmanagerClusterIssuer)
	certmanagerClusterIssuerYamlErr := os.WriteFile("kuberoCertmanagerClusterIssuer.yaml", certmanagerClusterIssuerYaml, 0644)
	if certmanagerClusterIssuerYamlErr != nil {
		fmt.Println(certmanagerClusterIssuerYamlErr)
		return
	}

	_, certmanagerClusterIssuerErr := exec.Command("kubectl", "apply", "-f", "kuberoCertmanagerClusterIssuer.yaml", "-n", "cert-manager").Output()
	if certmanagerClusterIssuerErr != nil {
		cfmt.Println("{{✗ Failed to create Certmanager Clusterissuer. Try runnig this command manually: kubectl apply -f kuberoCertmanagerClusterIssuer.yaml -n cert-manager}}::red")
		return
	} else {
		e := os.Remove("kuberoCertmanagerClusterIssuer.yaml")
		if e != nil {
			log.Fatal(e)
		}
		cfmt.Println("{{✓ Cert Manager Cluster Issuer created}}::lightGreen")
	}
}

func installOLMCertManager() {

	certManagerInstalled, _ := exec.Command("kubectl", "get", "deployment", "cert-manager-webhook", "-n", "operators").Output()
	if len(certManagerInstalled) > 0 {
		cfmt.Println("{{✓ Cert Manager already installed}}::lightGreen")
		return
	}

	certManagerSpinner := spinner.New("Install Cert Manager")
	certManagerSpinner.Start("run command : kubectl create -f https://operatorhub.io/install/cert-manager.yaml")
	_, certManagerErr := exec.Command("kubectl", "create", "-f", "https://operatorhub.io/install/cert-manager.yaml").Output()
	if certManagerErr != nil {
		certManagerSpinner.Error("Failed to run command. Try runnig this command manually: kubectl create -f https://operatorhub.io/install/cert-manager.yaml")
		log.Fatal(certManagerErr)
	}
	certManagerSpinner.Success("Cert Manager installed")

	time.Sleep(2 * time.Second)
	certManagerSpinner = spinner.New("Wait for Cert Manager to be ready")
	certManagerSpinner.Start("run command : kubectl wait --for=condition=available deployment/cert-manager-webhook -n cert-manager --timeout=180s -n operators")
	_, certManagerWaitErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/cert-manager-webhook", "-n", "cert-manager", "--timeout=180s", "-n", "operators").Output()
	if certManagerWaitErr != nil {
		certManagerSpinner.Error("Failed to run command. Try runnig it manually: kubectl wait --for=condition=available deployment/cert-manager-webhook -n cert-manager --timeout=180s -n operators")
		log.Fatal(certManagerWaitErr)
	}
	certManagerSpinner.Success("Cert Manager is ready")
}

func writeCLIconfig() {

	writeCli := true
	if kuberoConfig.configLoaded {
		writeCli = kuberoConfig.Install.WriteCliConfig
	} else {
		writeCli = promptLine("8) Write the Kubero CLI config", "[y,n]", "y") == "y"
	}

	if !writeCli {
		return
	}

	//TODO consider using SSL here.
	url := "http://" + arg_domain + ":" + arg_port
	token := arg_apiToken

	if kuberoConfig.configLoaded {
		// TODO: arg_port is not used ?
	} else {
		url = promptLine("Kubero Host address", "", "http://"+arg_domain+":"+arg_port)
		token = promptLine("Kubero Token", "", arg_apiToken)
	}
	viper.Set("api.url", url)
	viper.Set("api.token", token)

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", config)

	viper.WriteConfig()
}

func printDNSinfo() {

	ingressInstalled, err := exec.Command("kubectl", "get", "ingress", "-n", "kubero", "-o", "json").Output()
	if err != nil {
		cfmt.Println("{{✗ Failed to fetch DNS informations}}::red")
		return
	}
	var kuberoIngress KuberoIngress
	json.Unmarshal(ingressInstalled, &kuberoIngress)

	cfmt.Println("{{⚠ make sure your DNS is pointing to your Kubernetes cluster}}::yellow")

	//TODO this should be replaces by the default reviewapp domain
	if len(kuberoIngress.Items) > 0 &&
		len(kuberoIngress.Items[0].Spec.Rules[0].Host) > 0 &&
		len(kuberoIngress.Items[0].Status.LoadBalancer.Ingress) > 0 &&
		len(kuberoIngress.Items[0].Status.LoadBalancer.Ingress[0].IP) > 0 {
		cfmt.Printf("{{  %s.		IN		A		%s}}::lightBlue\n", kuberoIngress.Items[0].Spec.Rules[0].Host, kuberoIngress.Items[0].Status.LoadBalancer.Ingress[0].IP)
		cfmt.Printf("{{  *.review.example.com.			IN		A		%s}}::lightBlue", kuberoIngress.Items[0].Status.LoadBalancer.Ingress[0].IP)
	}

}

func finalMessage() {
	cfmt.Println(`

    ,--. ,--.        ,--.
    |  .'   /,--.,--.|  |-.  ,---. ,--.--. ,---.
    |  .   ' |  ||  || .-. '| .-. :|  .--'| .-. |
    |  |\   \'  ''  '| '-' |\   --.|  |   ' '-' '
    '--' '--' '----'  '---'  '----''--'    '---'

    Documentation:
    https://docs.kubero.dev
    `)

	if arg_domain != "" && arg_port != "" && arg_apiToken != "" && arg_adminPassword != "" {
		cfmt.Println(`
    Your Kubero UI :{{
    URL : ` + arg_domain + `:` + arg_port + `
    User: ` + arg_adminUser + `
    Pass: ` + arg_adminPassword + `}}::lightBlue
	`)
	} else {
		cfmt.Println("\n\n    {{Done - you can now login to your Kubero UI}}::lightGreen\n\n ")
	}
}

func generateRandomString(length int, chars string) string {
	if chars == "" {
		chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!+?._-%"
	}
	var letterRunes = []rune(chars)

	b := make([]rune, length)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func getAvailableStorageClasses() []string {
	var storageClasses []string
	storageClassesRaw, _ := exec.Command("kubectl", "get", "storageclasses", "-o", "json").Output()
	var storageClassesList StorageClassesList
	json.Unmarshal(storageClassesRaw, &storageClassesList)
	for _, storageClass := range storageClassesList.Items {
		storageClasses = append(storageClasses, storageClass.Metadata.Name)
	}
	return storageClasses
}

type StorageClassesList struct {
	APIVersion string `json:"apiVersion"`
	Items      []struct {
		APIVersion string `json:"apiVersion"`
		Kind       string `json:"kind"`
		Metadata   struct {
			Annotations struct {
				KubectlKubernetesIoLastAppliedConfiguration string `json:"kubectl.kubernetes.io/last-applied-configuration"`
				StorageclassKubernetesIoIsDefaultClass      string `json:"storageclass.kubernetes.io/is-default-class"`
			} `json:"annotations"`
			CreationTimestamp time.Time `json:"creationTimestamp"`
			Name              string    `json:"name"`
			ResourceVersion   string    `json:"resourceVersion"`
			UID               string    `json:"uid"`
		} `json:"metadata"`
		Provisioner       string `json:"provisioner"`
		ReclaimPolicy     string `json:"reclaimPolicy"`
		VolumeBindingMode string `json:"volumeBindingMode"`
	} `json:"items"`
	Kind     string `json:"kind"`
	Metadata struct {
		ResourceVersion string `json:"resourceVersion"`
		SelfLink        string `json:"selfLink"`
	} `json:"metadata"`
}
