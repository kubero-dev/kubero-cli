package kuberoCli

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
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
		rand.New(rand.NewSource(time.Now().UnixNano()))

		checkAllBinaries()

		switch argComponent {
		case "metrics":
			installMetrics()
			return
		case "certManager":
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
		case "monitoring":
			installMonitoring()
			return
		case "kubernetes":
			installKubernetes()
			checkCluster()
			return
		case "":
			printInstallSteps()
			installKubernetes()     // 1
			checkCluster()          //
			installOLM()            // 2
			installKuberoOperator() // 3
			installIngress()        // 4
			installMetrics()        // 5
			installCertManager()    // 6
			installMonitoring()     // 7
			installKuberoUi()       // 8
			writeCLIConfig()        // 9
			printDNSinfo()
			finalMessage()
			return
		default:
			return
		}
	},
}

var argAdminPassword string
var argAdminUser string
var argDomain string
var argApiToken string
var argPort string
var argPortSecure string
var clusterType string
var argComponent string
var installOlm bool
var monitoringInstalled bool
var ingressControllerVersion = "v1.10.0" // https://github.com/kubernetes/ingress-nginx/tags -> controller-v1.5.1

// var clusterTypeSelection = "[scaleway,linode,gke,digitalocean,kind]"
var clusterTypeList = []string{"kind", "linode", "scaleway", "gke", "digitalocean"}

func init() {
	installCmd.Flags().StringVarP(&argComponent, "component", "c", "", "install component (kubernetes,olm,ingress,metrics,certManager,kubero-operator,monitoring,kubero-ui)")
	installCmd.Flags().StringVarP(&argAdminUser, "user", "u", "", "Admin username for the kubero UI")
	installCmd.Flags().StringVarP(&argAdminPassword, "user-password", "U", "", "Password for the admin user")
	installCmd.Flags().StringVarP(&argApiToken, "apiToken", "a", "", "API token for the admin user")
	installCmd.Flags().StringVarP(&argPort, "port", "p", "", "Kubero UI HTTP port")
	installCmd.Flags().StringVarP(&argPortSecure, "securePort", "P", "", "Kubero UI HTTPS port")
	installCmd.Flags().StringVarP(&argDomain, "domain", "d", "", "Domain name for the kubero UI")
	rootCmd.AddCommand(installCmd)

	installOlm = false
	monitoringInstalled = false
}

func checkAllBinaries() {
	_, _ = cfmt.Println("\n  Check for required binaries")
	if !checkBinary("kubectl") {
		_, _ = cfmt.Println("{{✗ kubectl is not installed}}::red")
	} else {
		_, _ = cfmt.Println("{{✓ kubectl is installed}}::lightGreen")
	}

	if !checkBinary("kind") {
		_, _ = cfmt.Println("{{⚠ kind is not installed}}::yellow (only required if you want to install a local kind cluster)")
	} else {
		_, _ = cfmt.Println("{{✓ kind is installed}}::lightGreen")
	}

	if !checkBinary("gcloud") {
		_, _ = cfmt.Println("{{⚠ gcloud is not installed}}::yellow (only required if you want to install a GKE cluster)")
	} else {
		_, _ = cfmt.Println("{{✓ gcloud is installed}}::lightGreen")
	}
}

func printInstallSteps() {

	_, _ = cfmt.Print(`
  Steps to install kubero:
    1. Create a kubernetes cluster {{(optional)}}::gray
    2. Install the OLM {{(optional)}}::gray
    3. Install the kubero operator {{(required)}}::gray
    4. Install the ingress controller {{(required)}}::gray
    5. Install the metrics server {{(optional, but recommended)}}::gray
    6. Install the cert-manager {{(optional)}}::gray
    7. Install the monitoring stack {{(optional, but recommended)}}::gray
    8. Install the kubero UI {{(optional, but highly recommended)}}::gray
    9. Write the kubero CLI config
`)
}

func checkBinary(binary string) bool {
	_, err := exec.LookPath(binary)
	return err == nil
}

func installKubernetes() {
	kubernetesInstall := promptLine("1) Create a kubernetes cluster", "[y,n]", "y")
	if kubernetesInstall != "y" {
		return
	}

	clusterType = selectFromList("Select a Kubernetes provider", clusterTypeList, "")

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
		_, _ = cfmt.Println("{{✗ Unknown cluster type}}::red")
		os.Exit(1)
	}

}

func tellAChucknorrisJoke() {

	jokesApi := resty.New().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "kubero-cli/"+kuberoCliVersion).
		SetBaseURL("https://api.chucknorris.io/jokes/random")

	joke, _ := jokesApi.R().Get("?category=dev")
	var jokeResponse JokeResponse
	_ = json.Unmarshal(joke.Body(), &jokeResponse)
	_, _ = cfmt.Println("\r{{  " + jokeResponse.Value + "       }}::gray")
}

func mergeKubeconfig(kubeconfig []byte) error {

	newDefaultPathOptions := clientcmd.NewDefaultPathOptions()
	config1, _ := newDefaultPathOptions.GetStartingConfig()
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

	_ = clientcmd.ModifyConfig(clientcmd.DefaultClientConfig.ConfigAccess(), *config1, true)
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

	clusterSelect := promptLine("Is the CURRENT cluster the one you wish to install Kubero?", "[y,n]", "y")
	if clusterSelect == "n" {
		os.Exit(0)
	}
}

func installOLM() {

	openshiftInstalled, _ := exec.Command("kubectl", "get", "deployment", "olm-operator", "-n", "openshift-operator-lifecycle-manager").Output()
	if len(openshiftInstalled) > 0 {
		_, _ = cfmt.Println("{{✓ OLM is already installed}}::lightGreen")
		return
	}

	//namespace := promptLine("Install OLM in which namespace?", "[openshift-operator-lifecycle-manager,olm]", "olm")
	namespace := "olm"
	olmInstalled, _ := exec.Command("kubectl", "get", "deployment", "olm-operator", "-n", namespace).Output()
	if len(olmInstalled) > 0 {
		_, _ = cfmt.Println("{{✓ OLM is already installed}}::lightGreen")
		return
	}

	olmInstall := promptLine("2) Install OLM", "[y,n]", "n")
	if olmInstall != "y" {
		installOlm = false
		return
	} else {
		installOlm = true
	}

	olmVersionList := getGithubVersionList("operator-framework/operator-lifecycle-manager", 10)
	olmRelease := selectFromList("Select OLM version", olmVersionList, "")
	olmURL := "https://github.com/operator-framework/operator-lifecycle-manager/releases/download/" + olmRelease

	olmSpinner := spinner.New("Install OLM")

	olmCRDInstalled, _ := exec.Command("kubectl", "get", "crd", "subscriptions.operators.coreos.com").Output()
	if len(olmCRDInstalled) > 0 {
		_, _ = cfmt.Println("{{✓ OLM CRD's already installed}}::lightGreen")
	} else {
		_, _ = cfmt.Println("  run command : kubectl create -f " + olmURL + "/olm.yaml")
		olmSpinner.Start("Installing OLM CRDs")
		_, olmCRDErr := exec.Command("kubectl", "create", "-f", olmURL+"/crds.yaml").Output()
		if olmCRDErr != nil {
			fmt.Println("")
			olmSpinner.Error("OLM CRD installation failed. Try running this command manually: kubectl create -f " + olmURL + "/crds.yaml")
			log.Fatal(olmCRDErr)
		} else {
			olmSpinner.Success("OLM CRDs installed successfully")
		}
	}

	_, _ = cfmt.Println("  run command : kubectl create -f " + olmURL + "/olm.yaml")
	olmSpinner.Start("Install OLM")

	_, olmOLMErr := exec.Command("kubectl", "create", "-f", olmURL+"/olm.yaml").Output()
	if olmOLMErr != nil {
		fmt.Println("")
		olmSpinner.Error("Failed to run command. Try running this command manually: kubectl create -f " + olmURL + "/olm.yaml")
		log.Fatal(olmOLMErr)
	}
	olmSpinner.Success("OLM installed successfully")

	olmWaitSpinner := spinner.New("Wait for OLM to be ready")
	_, _ = cfmt.Println("  run command : kubectl wait --for=condition=available deployment/olm-operator -n " + namespace + " --timeout=180s")
	olmWaitSpinner.Start("Wait for OLM to be ready")
	_, olmWaitErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/olm-operator", "-n", namespace, "--timeout=180s").Output()
	if olmWaitErr != nil {
		olmWaitSpinner.Error("Failed to run command. Try running this command manually: kubectl wait --for=condition=available deployment/olm-operator -n " + namespace + " --timeout=180s")
		log.Fatal(olmWaitErr)
	}
	olmWaitSpinner.Success("OLM is ready")

	olmWaitCatalogSpinner := spinner.New("Wait for OLM Catalog to be ready")
	_, _ = cfmt.Println("  run command : kubectl wait --for=condition=available deployment/catalog-operator -n " + namespace + " --timeout=180s")
	olmWaitCatalogSpinner.Start("Wait for OLM Catalog to be ready")
	_, olmWaitCatalogErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/catalog-operator", "-n", namespace, "--timeout=180s").Output()
	if olmWaitCatalogErr != nil {
		olmWaitCatalogSpinner.Error("Failed to run command. Try running this command manually: kubectl wait --for=condition=available deployment/catalog-operator -n " + namespace + " --timeout=180s")
		log.Fatal(olmWaitCatalogErr)
	}
	olmWaitCatalogSpinner.Success("OLM Catalog is ready")
}

func installMetrics() {

	installed, _ := exec.Command("kubectl", "get", "deployments.apps", "metrics-server", "-n", "kube-system").Output()
	if len(installed) > 0 {
		_, _ = cfmt.Println("{{✓ Metrics is already enabled}}::lightGreen")
		return
	}
	install := promptLine("5) Install Kubernetes internal metrics service (required for HPA, Horizontal Pod Autoscaling)", "[y,n]", "y")
	if install != "y" {
		return
	}

	//components := "https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml"
	components := "https://raw.githubusercontent.com/kubero-dev/kubero-operator/main/deploy/metrics-server.yaml"
	_, installErr := exec.Command("kubectl", "apply", "-f", components).Output()

	if installErr != nil {
		fmt.Println("failed to install metrics server")
		log.Fatal(installErr)
	}
	_, _ = cfmt.Println("{{✓ Metrics server installed}}::lightGreen")
}

func installIngress() {

	ingressInstalled, _ := exec.Command("kubectl", "get", "ns", "ingress-nginx").Output()
	if len(ingressInstalled) > 0 {
		_, _ = cfmt.Println("{{✓ Ingress is already installed}}::lightGreen")
		return
	}

	ingressInstall := promptLine("4) Install Ingress", "[y,n]", "y")
	if ingressInstall != "y" {
		return
	} else {

		if clusterType == "" {
			clusterType = selectFromList("Which cluster type have you installed?", clusterTypeList, "")
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

		ingressProviderList := []string{"kind", "aws", "baremetal", "cloud", "do", "exoscale", "scw"}
		ingressProvider := selectFromList("Provider [kind, aws, baremetal, cloud(Azure,Google,Oracle,Linode), do(digital ocean), exoscale, scw(scaleway)]", ingressProviderList, prefill)

		// ingressController version can bot be loaded from GitHub api, since the return is alphabetic
		ingressSpinner := spinner.New("Install Ingress")
		URL := "https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-" + ingressControllerVersion + "/deploy/static/provider/" + ingressProvider + "/deploy.yaml"
		_, _ = cfmt.Println("  run command : kubectl apply -f " + URL)
		ingressSpinner.Start("Install Ingress")
		_, ingressErr := exec.Command("kubectl", "apply", "-f", URL).Output()
		if ingressErr != nil {
			ingressSpinner.Error("Failed to run command. Try running this command manually: kubectl apply -f " + URL)
			log.Fatal(ingressErr)
		}

		ingressSpinner.Success("Ingress installed successfully")
	}

}

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
		// kubectl api-resources --api-group=application.kubero.dev --no-headers=true
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
		// kubectl api-resources --api-group=application.kubero.dev --no-headers=true
		kuberoWait, _ = exec.Command("kubectl", "api-resources", "--api-group=application.kubero.dev", "--no-headers=true").Output()
		time.Sleep(1 * time.Second)
	}
	kuberoSpinner.UpdateMessage("Kubero Operator CRD's installed")

	time.Sleep(5 * time.Second)
	// kubectl wait --for=condition=available deployment/kubero -n kubero --timeout=180s
	kuberoSpinner.UpdateMessage("Wait for Kubero Operator to be ready")
	_, olmWaitErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/kubero-operator-controller-manager", "-n", "kubero-operator-system", "--timeout=300s").Output()
	if olmWaitErr != nil {
		kuberoSpinner.Error("Failed to wait for Kubero UI to become ready")
		log.Fatal(olmWaitErr)
	}
	kuberoSpinner.Success("Kubero Operator installed successfully")

}

func createNamespace(namespace string) {

	kuberoNSInstalled, _ := exec.Command("kubectl", "get", "ns", namespace).Output()
	if len(kuberoNSInstalled) > 0 {
		_, _ = cfmt.Printf("{{✓ Namespace %s exists}}::lightGreen\n", namespace)
	} else {
		_, kuberoNSErr := exec.Command("kubectl", "create", "namespace", namespace).Output()
		if kuberoNSErr != nil {
			fmt.Println("Failed to run command to create the namespace. Try running this command manually: kubectl create namespace " + namespace)
			log.Fatal(kuberoNSErr)
		} else {
			_, _ = cfmt.Printf("{{✓ Namespace %s created}}::lightGreen\n", namespace)
		}
	}
}

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

		var userDB []User
		userDB = append(userDB, User{Username: argAdminUser, Password: argAdminPassword, Insecure: true, ApiToken: argApiToken})
		userDBjson, _ := json.Marshal(userDB)
		userDBencoded := base64.StdEncoding.EncodeToString(userDBjson)

		createSecretCommand := exec.Command("kubectl", "create", "secret", "generic", "kubero-secrets",
			"--from-literal=KUBERO_WEBHOOK_SECRET="+webhookSecret,
			"--from-literal=KUBERO_SESSION_KEY="+sessionKey,
			"--from-literal=KUBERO_USERS="+userDBencoded,
		)

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

		createSecretCommand.Args = append(createSecretCommand.Args, "-n", "kubero")

		_, kuberoErr := createSecretCommand.Output()

		if kuberoErr != nil {
			_, _ = cfmt.Println("{{✗ Failed to run command to create the secrets.}}::red")
			log.Fatal(kuberoErr)
		} else {
			_, _ = cfmt.Println("{{✓ Kubero Secret created}}::lightGreen")
		}
	}

	kuberoUIInstalled, _ := exec.Command("kubectl", "get", "kuberoes.application.kubero.dev", "-n", "kubero").Output()
	if len(kuberoUIInstalled) > 0 {
		_, _ = cfmt.Println("{{✓ Kubero UI already installed}}::lightGreen")
	} else {
		installer := resty.New()

		installer.SetBaseURL("https://raw.githubusercontent.com")
		kf, _ := installer.R().Get("kubero-dev/kubero-operator/main/config/samples/application_v1alpha1_kubero.yaml")

		var kuberoUIConfig KuberoUIConfig
		_ = yaml.Unmarshal(kf.Body(), &kuberoUIConfig)

		if argDomain == "" {
			argDomain = promptLine("Kubero UI Domain", "", "kubero.localhost")
		}
		kuberoUIConfig.Spec.Ingress.Hosts[0].Host = argDomain

		// Warn if domain contains kubero.net

		webhookDomain := argDomain
		if strings.Contains(argDomain, "kubero.localhost") {
			_, _ = cfmt.Println("{{⚠ kubero.localhost might not be reachable won't get any Webhooks. GitHub will fail to connect your pipeline. }}::yellow")
			webhookDomain = "webhook.local.kubero.net"
		}

		webhookURL := promptLine("URL to which the webhooks should be sent (localhost fails with GitHub)", "", "https://"+webhookDomain+"/api/repo/webhooks")
		kuberoUIConfig.Spec.Kubero.WebhookURL = webhookURL

		kuberoUISsl := promptLine("Enable SSL for the Kubero UI", "[y/n]", "y")
		if kuberoUISsl == "y" {

			clusterIssuer := promptLine("Kubero UI ClusterIssuer", "", "letsencrypt-prod")
			kuberoUIConfig.Spec.Ingress.Annotations.KubernetesIoIngressClass = clusterIssuer
			kuberoUIConfig.Spec.Ingress.Annotations.KubernetesIoTlsAcme = "true"

			kuberoUIConfig.Spec.Ingress.TLS = []KuberoUITls{
				{
					Hosts:      []string{argDomain},
					SecretName: "kubero-tls",
				},
			}
		}

		kuberoUIRegistry := promptLine("Enable BuildPipeline for Kubero (BETA)", "[y/n]", "n")
		if kuberoUIRegistry == "y" {
			kuberoUIConfig.Spec.Registry.Enabled = true

			kuberoUICreateRegistry := promptLine("Create a local Registry for Kubero", "[y/n]", "n")
			if kuberoUICreateRegistry == "y" {
				kuberoUIConfig.Spec.Registry.Create = true

				kuberoUIRegistryStorage := promptLine("Registry storage size", "", "10Gi")
				kuberoUIConfig.Spec.Registry.Storage = kuberoUIRegistryStorage

				storageClassList := getAvailableStorageClasses()

				kuberoUIRegistryStorageClassName := selectFromList("Registry storage class", storageClassList, "")
				kuberoUIConfig.Spec.Registry.StorageClassName = kuberoUIRegistryStorageClassName
			}

			kuberoUIRegistryHost := promptLine("Registry", "[registry.kubero.mydomain.com]", "")
			kuberoUIConfig.Spec.Registry.Host = kuberoUIRegistryHost

			kuberoUIRegistrySubPath := promptLine("SubPath (optional) ", "[example/foo/bar]", "")
			kuberoUIConfig.Spec.Registry.SubPath = kuberoUIRegistrySubPath

			kuberoUIConfig.Spec.Registry.Port = 443

			kuberoUIRegistryUsername := promptLine("Registry username", "", "admin")
			kuberoUIConfig.Spec.Registry.Account.Username = kuberoUIRegistryUsername

			kuberoUIRegistryPassword := promptLine("Registry password", "", generateRandomString(12, ""))
			kuberoUIConfig.Spec.Registry.Account.Password = kuberoUIRegistryPassword

			kuberoUIRegistryPasswordBytes, _ := bcrypt.GenerateFromPassword([]byte(kuberoUIRegistryPassword), 14)
			kuberoUIConfig.Spec.Registry.Account.Hash = string(kuberoUIRegistryPasswordBytes)
		}

		kuberoUIAudit := promptLine("Enable Audit Logging", "[y/n]", "n")
		if kuberoUIAudit == "y" {
			kuberoUIConfig.Spec.Kubero.AuditLogs.Enabled = true

			storageClassList := getAvailableStorageClasses()

			kuberoUIRegistryStorageClassName := selectFromList("AuditLogs storage class", storageClassList, "")
			kuberoUIConfig.Spec.Kubero.AuditLogs.StorageClassName = kuberoUIRegistryStorageClassName

		}

		if monitoringInstalled {
			kuberoUIConfig.Spec.Prometheus.Enabled = true
			kuberoUIConfig.Spec.Prometheus.Endpoint = promptLine("Prometheus URL", "", "http://kubero-prometheus-server")
		} else {
			kuberoUIConfig.Spec.Prometheus.Enabled = false
		}

		kuberoUIConsole := promptLine("Enable Console Access to running containers", "[y/n]", "y")

		if kuberoUIConsole == "y" {
			kuberoUIConfig.Spec.Kubero.Config.Kubero.Console.Enabled = true
		}

		//kuberoUIConfig.Spec.Image.Tag = "v2.0.0-rc.8"

		if clusterType == "" {
			clusterType = selectFromList("Which cluster type have you installed?", clusterTypeList, "")
		}

		if clusterType == "linode" ||
			clusterType == "digitalocean" ||
			clusterType == "scaleway" ||
			clusterType == "gke" {
			kuberoUIConfig.Spec.Ingress.ClassName = "nginx"
		}

		kuberoUIYaml, _ := yaml.Marshal(kuberoUIConfig)
		kuberoUIErr := os.WriteFile("kuberoUI.yaml", kuberoUIYaml, 0644)

		if kuberoUIErr != nil {
			fmt.Println(kuberoUIErr)
			return
		}

		_, olmInstallErr := exec.Command("kubectl", "apply", "-f", "kuberoUI.yaml", "-n", "kubero").Output()
		if olmInstallErr != nil {
			_, _ = cfmt.Println("{{✗ Failed to run command to install Kubero UI. Try running this command manually: kubectl apply -f kuberoUI.yaml -n kubero}}::red")
			return
		} else {
			e := os.Remove("kuberoUI.yaml")
			if e != nil {
				log.Fatal(e)
			}
			_, _ = cfmt.Println("{{✓ Kubero UI installed}}::lightGreen")
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
		_, olmWaitErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/kubero", "-n", "kubero", "--timeout=300s").Output()
		if olmWaitErr != nil {
			kuberoUISpinner.Error("Failed to wait for Kubero UI to become ready")
			log.Fatal(olmWaitErr)
		}
		kuberoUISpinner.Success("Kubero UI is ready")
	}

}

func installMonitoring() {

	if promptLine("7) Enable long-term metrics", "[y/n]", "y") == "y" {
		monitoringInstalled = true
	} else {
		monitoringInstalled = false
		return
	}

	createNamespace("kubero")

	spinnerObj := spinner.New("enable metrics")
	if promptLine("7.1) Create local Prometheus instance", "[y/n]", "y") == "y" {
		URL := "https://raw.githubusercontent.com/kubero-dev/kubero-operator/main/config/samples/application_v1alpha1_kuberoprometheus.yaml"
		_, _ = cfmt.Println("  run command : kubectl apply -n kubero -f " + URL)
		spinnerObj.Start("Installing Prometheus")
		_, ingressErr := exec.Command("kubectl", "apply", "-n", "kubero", "-f", URL).Output()
		if ingressErr != nil {
			spinnerObj.Error("Failed to run command. Try running this command manually: kubectl apply -f " + URL)
			log.Fatal(ingressErr)
		}
		/*
			spinner.UpdateMessage("Waiting for Prometheus to be ready")

			time.Sleep(5 * time.Second)
			// kubectl wait --for=condition=available deployment/kubero -n kubero --timeout=180s
			_, olmWaitErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/kubero-prometheus-server", "-n", "kubero", "--timeout=300s").Output()
			if olmWaitErr != nil {
				spinner.Error("Failed to wait for Prometheus to become ready")
				log.Fatal(olmWaitErr)
			}
		*/
		spinnerObj.Success("Prometheus installed successfully")
	}

	if promptLine("7.2) Enable KubeMetrics", "[y/n]", "y") == "y" {
		_, _ = cfmt.Println("  run command : kubectl patch kuberoes kubero -n kubero --type=merge")
		spinnerObj.Start("Enabling Metrics")

		patch := `{
			"spec": {
				"prometheus": {
					"kube-state-metrics": {
						"enabled": true
					}
				}
			}
		}`

		_, patchResult := exec.Command("kubectl", "patch", "kuberoprometheuses", "kubero-prometheus", "-n", "kubero", "--type=merge", "-p", patch).Output()
		if patchResult != nil {
			spinnerObj.Error("Failed to patch the kubero prometheus CRD to enable kube metrics", patchResult.Error(), patch)
		}
		spinnerObj.Success("Metrics enabled successfully")

	}

	patch := `{
		"spec": {
		  "template": {
			"metadata": {
			  "annotations": {
				"prometheus.io/port": "10254",
				"prometheus.io/scrape": "true"
			  }
			},
			"spec": {
			  "containers": [
				{
				  "name": "controller",
				  "ports": [
					{
					  "containerPort": 10254,
					  "name": "prometheus",
					  "protocol": "TCP"
					}
				  ],
				  "args": [
					"/nginx-ingress-controller",
					"--election-id=ingress-nginx-leader",
					"--controller-class=k8s.io/ingress-nginx",
					"--ingress-class=nginx",
					"--configmap=$(POD_NAMESPACE)/ingress-nginx-controller",
					"--validating-webhook=:8443",
					"--validating-webhook-certificate=/usr/local/certificates/cert",
					"--validating-webhook-key=/usr/local/certificates/key",
					"--watch-ingress-without-class=true",
					"--enable-metrics=true",
					"--publish-status-address=localhost"
				  ]
				}
			  ]
			}
		  }
		}
	  }`
	_, ingressPatch := exec.Command("kubectl", "patch", "deployments.apps", "ingress-nginx-controller", "-n", "ingress-nginx", "-p", patch).Output()
	if ingressPatch != nil {
		_, _ = cfmt.Println("{{✗ Failed to patch the ingress controller. }}::red\nHere is a detailed information how to do it manually: https://github.com/kubernetes/ingress-nginx/blob/main/docs/user-guide/monitoring.md")
		//log.Fatal(ingressPatch)
	}

	patch = `{
		"spec": {
			"ports": [
				{
					"name": "prometheus",
					"nodePort": 31280,
					"port": 10254,
					"protocol": "TCP",
					"targetPort": "prometheus"
				}
			]
		}
	}`

	_, ingressPatch = exec.Command("kubectl", "patch", "svc", "ingress-nginx-controller", "-n", "ingress-nginx", "-p", patch).Output()
	if ingressPatch != nil {
		_, _ = cfmt.Println("{{✗ Failed to patch the ingress controller service. }}::red\nHere is a detailed information how to do it manually: https://github.com/kubernetes/ingress-nginx/blob/main/docs/user-guide/monitoring.md")
		//log.Fatal(ingressPatch)
	}

}

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
	//certManagerUrl := "https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml"
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
	for i := 0; i < 4; i++ {
		tellAChucknorrisJoke()
		time.Sleep(15 * time.Second)
	}
	_, certManagerWaitErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/cert-manager-webhook", "-n", "cert-manager", "--timeout=180s", "-n", "operators").Output()
	if certManagerWaitErr != nil {
		certManagerSpinner.Error("Failed to run command. Try running it manually: kubectl wait --for=condition=available deployment/cert-manager-webhook -n cert-manager --timeout=180s -n operators")
		log.Fatal(certManagerWaitErr)
	}
	certManagerSpinner.Success("Cert Manager is ready")

	installCertManagerClusterIssuer("default")
}

func writeCLIConfig() {

	ingressInstall := promptLine("10) Write the Kubero CLI config", "[y,n]", "n")
	if ingressInstall != "y" {
		return
	}

	//TODO consider using SSL here.
	url := promptLine("Kubero Host address", "", "http://"+argDomain+":"+argPort)
	viper.Set("api.url", url)

	token := promptLine("Kubero Token", "", argApiToken)
	viper.Set("api.token", token)

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", config)

	_ = viper.WriteConfig()
}

func printDNSinfo() {

	ingressInstalled, err := exec.Command("kubectl", "get", "ingress", "-n", "kubero", "-o", "json").Output()
	if err != nil {
		_, _ = cfmt.Println("{{✗ Failed to fetch DNS information}}::red")
		return
	}
	var kuberoIngress KuberoIngress
	_ = json.Unmarshal(ingressInstalled, &kuberoIngress)

	_, _ = cfmt.Println("{{⚠ make sure your DNS is pointing to your Kubernetes cluster}}::yellow")

	//TODO this should be replaces by the default review app domain
	if len(kuberoIngress.Items) > 0 &&
		len(kuberoIngress.Items[0].Spec.Rules[0].Host) > 0 &&
		len(kuberoIngress.Items[0].Status.LoadBalancer.Ingress) > 0 &&
		len(kuberoIngress.Items[0].Status.LoadBalancer.Ingress[0].IP) > 0 {
		_, _ = cfmt.Printf("{{  %s.		IN		A		%s}}::lightBlue\n", kuberoIngress.Items[0].Spec.Rules[0].Host, kuberoIngress.Items[0].Status.LoadBalancer.Ingress[0].IP)
		_, _ = cfmt.Printf("{{  *.review.example.com.			IN		A		%s}}::lightBlue", kuberoIngress.Items[0].Status.LoadBalancer.Ingress[0].IP)
	}

}

func finalMessage() {
	_, _ = cfmt.Println(`

    ,--. ,--.        ,--.
    |  .'   /,--.,--.|  |-.  ,---. ,--.--. ,---.
    |  .   ' |  ||  || .-. '| .-. :|  .--'| .-. |
    |  |\   \'  ''  '| '-' |\   --.|  |   ' '-' '
    '--' '--' '----'  '---'  '----''--'    '---'

    Documentation:
    https://docs.kubero.dev
    `)

	protocol := "https"
	if argPort == "80" {
		protocol = "http"
	}
	_, _ = cfmt.Println(`
    Your Kubero UI :{{
    URL : ` + protocol + "://" + argDomain + ":" + argPort + `
    User: ` + argAdminUser + `
    Pass: ` + argAdminPassword + `}}::lightBlue
	`)

	_, _ = cfmt.Println("\n\n    {{Done - you can now login to your Kubero UI}}::lightGreen\n\n ")

}

func generateRandomString(length int, chars string) string {
	if chars == "" {
		chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!+?._-%"
	}
	var letterRunes = []rune(chars)

	b := make([]rune, length)
	rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func getAvailableStorageClasses() []string {
	var storageClasses []string
	storageClassesRaw, _ := exec.Command("kubectl", "get", "storageClasses", "-o", "json").Output()
	var storageClassesList StorageClassesList
	_ = json.Unmarshal(storageClassesRaw, &storageClassesList)
	for _, storageClass := range storageClassesList.Items {
		storageClasses = append(storageClasses, storageClass.Metadata.Name)
	}
	return storageClasses
}

func getGithubVersionList(repository string, limit int) []string {

	githubapi := resty.New().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "kubero-cli/"+kuberoCliVersion).
		SetBaseURL("https://api.github.com/repos/")

	tags, _ := githubapi.R().Get(repository + "/tags")
	var versions []GithubVersion
	_ = json.Unmarshal(tags.Body(), &versions)

	var versionList []string
	versionList = make([]string, 0)

	for _, version := range versions {
		if limit == 0 || len(versionList) < limit {
			versionList = append(versionList, version.Name)
		}
	}

	return versionList
}

type StorageClassesList struct {
	APIVersion string `json:"apiVersion"`
	Items      []struct {
		APIVersion string `json:"apiVersion"`
		Kind       string `json:"kind"`
		Metadata   struct {
			Annotations struct {
				KubectlKubernetesIoLastAppliedConfiguration string `json:"kubectl.kubernetes.io/last-applied-configuration"`
				StorageClassKubernetesIoIsDefaultClass      string `json:"storageClass.kubernetes.io/is-default-class"`
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
