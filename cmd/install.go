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
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/tools/clientcmd"
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
		installSwitch()
		checkCluster()
		installOLM()
		installIngress()
		installCertManager()
		installKuberoOperator()
		installKuberoUi()
		writeCLIconfig()
		printDNSinfo()
		finalMessage()
	},
}

var arg_adminPassword string
var arg_adminUser string
var arg_domain string
var arg_apiToken string
var arg_port string
var arg_portSecure string
var clusterType string
var ingressControllerVersion = "v1.5.1" // https://github.com/kubernetes/ingress-nginx/tags -> controller-v1.5.1

var clusterTypeSelection = "[scaleway,linode,gke,digitalocean,kind]"

func init() {
	installCmd.Flags().StringVarP(&arg_adminUser, "user", "u", "", "Admin username for the kubero UI")
	installCmd.Flags().StringVarP(&arg_adminPassword, "user-password", "U", "", "Password for the admin user")
	installCmd.Flags().StringVarP(&arg_apiToken, "apitoken", "a", "", "API token for the admin user")
	installCmd.Flags().StringVarP(&arg_port, "port", "p", "", "Kubero UI HTTP port")
	installCmd.Flags().StringVarP(&arg_portSecure, "secureport", "P", "", "Kubero UI HTTPS port")
	installCmd.Flags().StringVarP(&arg_domain, "domain", "d", "", "Domain name for the kubero UI")
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

	if !checkBinary("gcloud") {
		cfmt.Println("{{⚠ gcloud is not installed}}::yellow (only required if you want to install a GKE cluster)")
	} else {
		cfmt.Println("{{✓ gcloud is installed}}::lightGreen")
	}
}

func checkBinary(binary string) bool {
	_, err := exec.LookPath(binary)
	return err == nil
}

func installSwitch() {
	kubernetesInstall := promptLine("Start a kubernetes Cluster", "[y,n]", "y")
	if kubernetesInstall != "y" {
		return
	}

	clusterType = promptLine("Select a cluster type", clusterTypeSelection, "linode")

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

func installScaleway() {

	// create the cluster
	// https://api.scaleway.com/k8s/v1/regions/{region}/clusters/{cluster_id}/available-versions

	// check state of cluster
	// https://api.scaleway.com/k8s/v1/regions/{region}/clusters

	// get the kubeconfig
	// https://api.scaleway.com/k8s/v1/regions/{region}/clusters/{cluster_id}/kubeconfig

	cfmt.Println("{{⚠ Installing Kubernetes on Scaleway is currently beta state in kubero-cli}}::yellow")
	cfmt.Println("{{  Please report if you run into errors}}::yellow")

	var cluster ScalewayCreate

	cluster.ProjectID = os.Getenv("SCALEWAY_PROJECTID")
	if cluster.ProjectID == "" {
		cfmt.Println("{{✗ SCALEWAY_PROJECTID is not set}}::red")
		log.Fatal("missing SCALEWAY_PROJECTID")
	}
	cluster.OrganizationID = os.Getenv("SCALEWAY_PROJECTID")

	token := os.Getenv("SCALEWAY_ACCESS_TOKEN")
	if token == "" {
		cfmt.Println("{{✗ SCALEWAY_ACCESS_TOKEN is not set}}::red")
		log.Fatal("missing SCALEWAY_ACCESS_TOKEN")
	}

	api := resty.New().
		SetHeader("X-Auth-Token", token).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "kubero-cli/0.0.1").
		SetBaseURL("https://api.scaleway.com/k8s/v1/regions")

	cluster.Name = promptLine("Kubernetes Cluster Name", "", "kubero-"+strconv.Itoa(rand.Intn(1000)))
	region := promptLine("Cluster Region", "[fr-par,nl-ams,pl-waw]", "nl-ams")
	cluster.Version = promptLine("Kubernetes Version", "[1.23.13,1.22.15,1.21.14]", "1.24.7")

	// TODO lets make this configurable if needed in the future
	cluster.Cni = "unknown_cni"
	cluster.Ingress = "unknown_ingress" // is marked as deprecated in the api but required for now
	/*
		TODO : not implemented yet, but prepare for it in the future
		cluster.AutoscalerConfig.Estimator = "unknown_estimator"
		cluster.AutoscalerConfig.Expander = "unknown_expander"
		cluster.AutoscalerConfig.ScaleDownUtilizationThreshold = 0.5
		cluster.AutoscalerConfig.MaxGracefulTerminationSec = 60
	*/
	cluster.AutoUpgrade.Enable = false
	cluster.AutoUpgrade.MaintenanceWindow.StartHour = 3
	cluster.AutoUpgrade.MaintenanceWindow.Day = "any"

	// TODO load the options from the api
	nodeType := promptLine("Node Types", "[DEV1-M,DEV1-XL,GP1-M]", "DEV1-M")

	cluster.Pools = append(cluster.Pools, ScalewayNodePool{
		Name:             "default",
		NodeType:         nodeType,
		Autoscaling:      false,
		Size:             3,
		ContainerRuntime: "unknown_runtime",
		RootVolumeType:   "default_volume_type",
		//RootVolumeSize:   50,
	})

	//fmt.Printf("%+v\n", cluster)
	newCluster, _ := api.R().SetBody(cluster).Post(region + "/clusters")

	var clusterResponse ScalewayCreateResponse
	if newCluster.StatusCode() < 299 {
		json.Unmarshal(newCluster.Body(), &clusterResponse)
		cfmt.Println("{{✓ Scaleway Kubernetes cluster created}}::lightGreen")
	} else {
		cfmt.Println("{{✗ Scaleway Kubernetes Cluster creation failed}}::red")
		log.Fatal(string(newCluster.Body()))
	}

	spinner := spinner.New()
	spinner.Start("Waiting for cluster to be ready")
	for {
		clusterStatus, _ := api.R().Get(region + "/clusters/" + clusterResponse.ID)
		var clusterStatusResponse ScalewayCreateResponse
		json.Unmarshal(clusterStatus.Body(), &clusterStatusResponse)
		if clusterStatusResponse.Status == "ready" {
			spinner.Success("Scaleway Kubernetes Cluster is ready")
			break
		}
		time.Sleep(2 * time.Second)
	}

	k, _ := api.R().Get(region + "/clusters/" + clusterResponse.ID + "/kubeconfig")

	var scalewayKubeconfigResponse ScalewayKubeconfigResponse
	json.Unmarshal(k.Body(), &scalewayKubeconfigResponse)
	kubeconfig, _ := base64.StdEncoding.DecodeString(scalewayKubeconfigResponse.Content)

	err := mergeKubeconfig([]byte(kubeconfig))
	if err != nil {
		cfmt.Println("{{✗ Failed to download kubeconfig}}::red")
		log.Fatal(err)
	} else {
		cfmt.Println("{{✓ Kubeconfig downloaded}}::lightGreen")
	}

}

func installGKE() {
	// implememted with gcloud, since it is required for the download of the kubeconfig anyway

	// gcloud config list
	// gcloud config get project
	// gcloud container clusters create kubero-cluster-4 --region=us-central1-c
	// gcloud container clusters get-credentials kubero-cluster-4 --region=us-central1-c

	// https://cloud.google.com/kubernetes-engine/docs/reference/libraries#client-libraries-install-go
	// https://github.com/googleapis/google-cloud-go

	gcloudName := promptLine("Kubernetes Cluster Name", "", "kubero-"+strconv.Itoa(rand.Intn(1000)))
	gcloudRegion := promptLine("Region", "[https://cloud.google.com/compute/docs/regions-zones]", "us-central1-c")
	gcloudClusterVersion := promptLine("Cluster Version", "[https://cloud.google.com/kubernetes-engine/docs/release-notes-regular]", "1.23.8-gke.1900")

	spinner := spinner.New("Spin up a GKE cluster")
	spinner.Start("run command : gcloud container clusters create " + gcloudName + " --region=" + gcloudRegion + " --cluster-version=" + gcloudClusterVersion)
	_, err := exec.Command("gcloud", "container", "clusters", "create", gcloudName,
		"--region="+gcloudRegion,
		"--cluster-version="+gcloudClusterVersion).Output()
	if err != nil {
		fmt.Println()
		spinner.Error("Failed to run command. Try runnig it manually and skip this step")
		log.Fatal(err)
	}
	spinner.Success("GKE cluster started sucessfully")

	spinner.Start("Get credentials for the GKE cluster")
	_, err = exec.Command("gcloud", "container", "clusters", "get-credentials", gcloudName, "--region="+gcloudRegion).Output()
	if err != nil {
		fmt.Println()
		spinner.Error("Failed to run command. Try runnig it manually and skip this step")
		log.Fatal(err)
	} else {
		spinner.Success("GKE cluster credentials set")
	}

}

func installLinode() {
	// https://www.linode.com/docs/api/linode-kubernetes-engine-lke/#kubernetes-cluster-create
	// https://www.linode.com/docs/api/linode-kubernetes-engine-lke/#kubernetes-cluster-view
	// https://www.linode.com/docs/api/linode-kubernetes-engine-lke/#kubeconfig-view

	cfmt.Println("{{⚠ Installing Kubernetes on Linode is currently beta state in kubero-cli}}::yellow")
	cfmt.Println("{{  Please report if you run into errors}}::yellow")

	token := os.Getenv("LINODE_ACCESS_TOKEN")
	if token == "" {
		cfmt.Println("{{✗ LINODE_ACCESS_TOKEN is not set}}::red")
		log.Fatal("missing LINODE_ACCESS_TOKEN")
	}

	api := resty.New().
		SetAuthScheme("Bearer").
		SetAuthToken(token).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "kubero-cli/0.0.1").
		SetBaseURL("https://api.linode.com/v4/lke/clusters")

	var clusterConfig LinodeCreateClusterRequest
	clusterConfig.Label = promptLine("Cluster Name", "", "kubero-"+strconv.Itoa(rand.Intn(1000)))
	clusterConfig.Region = promptLine("Region", "[https://www.linode.com/global-infrastructure/]", "us-central") // TODO load the list of regions or point to e better document

	workerNodesCount, _ := strconv.Atoi(promptLine("Worker Nodes Count", "", "3"))
	workerNodesType := promptLine("Worker Nodes Type", "[https://www.linode.com/pricing/]", "g6-standard-2") // TODO load the list of types or point to e better document

	clusterConfig.K8SVersion = promptLine("Kubernetes Version", "[1.23]", "1.23")
	clusterConfig.Tags = []string{"kubero"}
	clusterConfig.NodePools = []LinodeNodepool{
		{
			Type:  workerNodesType,
			Count: workerNodesCount,
		},
	}

	spinner := spinner.New("Spin up a Linode Kubernetes Cluster")

	spinner.Start("Create Linode Kubernetes Cluster")
	clusterResponse, _ := api.R().SetBody(clusterConfig).Post("")
	if clusterResponse.StatusCode() > 299 {
		fmt.Println()
		spinner.Error("Failed to create Linode Kubernetes Cluster")
		log.Fatal(clusterResponse.String())
	}
	spinner.Success("Linode Kubernetes Cluster created")

	var cluster LinodeCreateClusterResponse
	json.Unmarshal(clusterResponse.Body(), &cluster)

	// According to the docs, the cluster is ready after 2-5 minutes.
	cfmt.Println("{{  Wait for Linode Kubernetes Cluster to be ready}}::lightBlue")
	cfmt.Println("{{  According to the docs this may take up to 7 minutes}}::lightBlue")
	cfmt.Println("{{  Time for a coffee break and some Chuck Norris jokes.}}::lightBlue")
	spinner.Start("Wait for Linode Kubernetes Cluster to be ready")

	var LinodeKubeconfig struct {
		Kubeconfig string `json:"kubeconfig"`
	}

	for i := 0; true; i++ {
		time.Sleep(15 * time.Second)
		r, _ := api.R().SetResult(&LinodeKubeconfig).Get("/" + strconv.Itoa(cluster.ID) + "/kubeconfig")
		if r.StatusCode() > 299 {
			tellAChucknorrisJoke()
		}
		if LinodeKubeconfig.Kubeconfig != "" {
			spinner.Success("Linode Kubernetes Cluster is ready")
			break
		}
	}
	kubeconfig, err := base64.StdEncoding.DecodeString(LinodeKubeconfig.Kubeconfig)

	if err != nil {
		fmt.Println()
		spinner.Error("Failed to decode kubeconfig")
		log.Fatal(err)
	}

	err = mergeKubeconfig(kubeconfig)
	if err != nil {
		fmt.Println()
		spinner.Error("Failed to merge kubeconfig")
		log.Fatal(err)
	}

	spinner.Success("Linode Kubernetes Cluster credentials set")

}

func tellAChucknorrisJoke() {

	jokesapi := resty.New().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "kubero-cli/0.0.1").
		SetBaseURL("https://api.chucknorris.io/jokes/random")

	joke, _ := jokesapi.R().Get("?category=dev")
	var jokeResponse JokeResponse
	json.Unmarshal(joke.Body(), &jokeResponse)
	cfmt.Println("\r{{  " + jokeResponse.Value + "       }}::gray")
}

func installDigitalOcean() {
	// https://docs.digitalocean.com/reference/api/api-reference/#operation/kubernetes_create_cluster

	cfmt.Println("{{⚠ Installing Kubernetes on Digital Ocean is currently beta state in kubero-cli}}::yellow")
	cfmt.Println("{{  Please report if you run into errors}}::yellow")

	token := os.Getenv("DIGITALOCEAN_ACCESS_TOKEN")
	if token == "" {
		cfmt.Println("{{✗ DIGITALOCEAN_ACCESS_TOKEN is not set}}::red")
		log.Fatal("missing DIGITALOCEAN_ACCESS_TOKEN")
	}

	doApi := resty.New().
		SetAuthScheme("Bearer").
		SetAuthToken(token).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "kubero-cli/0.0.1"). //TODO dynamic version
		SetBaseURL("https://api.digitalocean.com")

	var doConfig DigitalOceanKubernetesConfig
	doConfig.NodePools = []struct {
		Size  string `json:"size"`
		Count int    `json:"count"`
		Name  string `json:"name"`
	}{
		{
			Size:  "s-1vcpu-2gb",
			Count: 1,
			Name:  "worker-pool",
		},
	}

	doConfig.Name = promptLine("Kubernetes Cluster Name", "", "kubero-"+strconv.Itoa(rand.Intn(1000)))
	doConfig.Region = promptLine("Cluster Region", "[nyc1,sgp1,lon1,ams3,fra1,...]", "nyc1")
	doConfig.Version = promptLine("Cluster Version", "[1.24.4-do.0,1.17.11-do.0,1.16.14-do.0]", "1.24.4-do.0")

	doConfig.NodePools[0].Size = promptLine("Cluster Node Size", "[s-1vcpu-2gb,s-2vcpu-4gb,s-4vcpu-8gb,s-8vcpu-16gb,s-16vcpu-32gb,s-32vcpu-64gb,s-48vcpu-96gb,s-64vcpu-128gb]", "s-1vcpu-2gb")
	doConfig.NodePools[0].Count, _ = strconv.Atoi(promptLine("Cluster Node Count", "", "1"))

	kf, _ := doApi.R().
		SetBody(doConfig).
		Post("/v2/kubernetes/clusters")

	if kf.StatusCode() > 299 {
		fmt.Println(kf.String())
		cfmt.Println("{{✗ failed to create digital ocean cluster}}::red")
		os.Exit(1)
	} else {
		cfmt.Println("{{✓ digital ocean cluster created}}::lightGreen")
	}

	var doCluster DigitalOcean
	json.Unmarshal(kf.Body(), &doCluster)

	doSpinner := spinner.New("Starting a kubernetes cluster on digital ocean")
	doSpinner.Start("Waiting for digital ocean cluster to be ready. This may take a few minutes. Time enough to get a coffee ☕")
	clusterID := doCluster.KubernetesCluster.ID

	for {
		time.Sleep(2 * time.Second)
		doWait, _ := doApi.R().
			Get("/v2/kubernetes/clusters/" + clusterID)

		if doWait.StatusCode() > 299 {
			fmt.Println(doWait.String())
			doSpinner.Error("Failed to create digital ocean cluster")
			continue
		} else {
			var doCluster DigitalOcean
			json.Unmarshal(doWait.Body(), &doCluster)
			if doCluster.KubernetesCluster.Status.State == "running" {
				doSpinner.Success("digital ocean cluster created")
				break
			}
		}
	}

	kubectl, _ := doApi.R().
		Get("v2/kubernetes/clusters/" + clusterID + "/kubeconfig")
	mergeKubeconfig(kubectl.Body())

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

func installKind() {

	if !checkBinary("kind") {
		log.Fatal("kind binary is not installed")
	}

	installer := resty.New()

	installer.SetBaseURL("https://raw.githubusercontent.com")
	kf, _ := installer.R().Get("/kubero-dev/kubero/main/kind.yaml")

	var kindConfig KindConfig
	yaml.Unmarshal(kf.Body(), &kindConfig)

	kindConfig.Name = promptLine("Kind Cluster Name", "", "kubero-"+strconv.Itoa(rand.Intn(1000)))
	kindConfig.Nodes[0].Image = "kindest/node:v1.25.3" //TODO make configurable version

	if arg_port == "" {
		arg_port = promptLine("Local HTTP Port", "", "80")
	}
	kindConfig.Nodes[0].ExtraPortMappings[0].HostPort, _ = strconv.Atoi(arg_port)

	if arg_portSecure == "" {
		arg_portSecure = promptLine("Local HTTPS Port", "", "443")
	}
	kindConfig.Nodes[0].ExtraPortMappings[1].HostPort, _ = strconv.Atoi(arg_portSecure)

	kindConfigYaml, _ := yaml.Marshal(kindConfig)
	//fmt.Println("-------------- kind.yaml ---------------")
	//fmt.Println(string(kindConfigYaml))
	//fmt.Println("----------------------------------------")

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

	olmSpinner := spinner.New("Install OLM")

	olmCRDInstalled, _ := exec.Command("kubectl", "get", "crd", "subscriptions.operators.coreos.com").Output()
	if len(olmCRDInstalled) > 0 {
		cfmt.Println("{{✓ OLM CRD's allredy installed}}::lightGreen")
	} else {
		olmSpinner.Start("run command : kubectl create -f " + olmURL + "/olm.yaml")
		_, olmCRDErr := exec.Command("kubectl", "create", "-f", olmURL+"/crds.yaml").Output()
		if olmCRDErr != nil {
			fmt.Println("")
			olmSpinner.Error("OLM CRD installation failed. Try runnig it manually")
			log.Fatal(olmCRDErr)
		} else {
			olmSpinner.Success("OLM CRDs installed sucessfully")
		}
	}

	olmSpinner.Start("run command : kubectl create -f " + olmURL + "/olm.yaml")

	_, olmOLMErr := exec.Command("kubectl", "create", "-f", olmURL+"/olm.yaml").Output()
	if olmOLMErr != nil {
		fmt.Println("")
		olmSpinner.Error("Failed to run command. Try runnig it manually")
		log.Fatal(olmOLMErr)
	}
	olmSpinner.Success("OLM installed sucessfully")

	olmWaitSpinner := spinner.New("Wait for OLM to be ready")
	olmWaitSpinner.Start("run command : kubectl wait --for=condition=available deployment/olm-operator -n " + namespace + " --timeout=180s")
	_, olmWaitErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/olm-operator", "-n", namespace, "--timeout=180s").Output()
	if olmWaitErr != nil {
		olmWaitSpinner.Error("Failed to run command. Try runnig it manually")
		log.Fatal(olmWaitErr)
	}
	olmWaitSpinner.Success("OLM is ready")

	olmWaitCatalogSpinner := spinner.New("Wait for OLM Catalog to be ready")
	olmWaitCatalogSpinner.Start("run command : kubectl wait --for=condition=available deployment/catalog-operator -n " + namespace + " --timeout=180s")
	_, olmWaitCatalogErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/catalog-operator", "-n", namespace, "--timeout=180s").Output()
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
		return
	} else {

		if clusterType == "" {
			clusterType = promptLine("Which cluster type have you insalled?", clusterTypeSelection, "")
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

		ingressProvider := promptLine("Provider", "[kind,aws,baremetal,cloud(Azure,Google,Oracle,Linode),do(digital ocean),exoscale,scw(scaleway)]", prefill)
		ingressSpinner := spinner.New("Install Ingress")
		URL := "https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-" + ingressControllerVersion + "/deploy/static/provider/" + ingressProvider + "/deploy.yaml"
		ingressSpinner.Start("run command : kubectl apply -f " + URL)
		_, ingressErr := exec.Command("kubectl", "apply", "-f", URL).Output()
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

		if arg_adminUser == "" {
			arg_adminUser = promptLine("Admin User", "", "admin")
		}

		if arg_adminPassword == "" {
			arg_adminPassword = promptLine("Admin Password", "", generatePassword(12))
		}

		if arg_apiToken == "" {
			arg_apiToken = promptLine("Random string for admin API token", "", generatePassword(20))
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
		installer := resty.New()

		installer.SetBaseURL("https://raw.githubusercontent.com")
		kf, _ := installer.R().Get("kubero-dev/kubero-operator/main/config/samples/application_v1alpha1_kubero.yaml")

		var kuberiUIConfig KuberoUIConfig
		yaml.Unmarshal(kf.Body(), &kuberiUIConfig)

		if arg_domain == "" {
			arg_domain = promptLine("Kuberi UI Domain", "", "kubero.lacolhost.com")
		}
		kuberiUIConfig.Spec.Ingress.Hosts[0].Host = arg_domain

		webhookURL := promptLine("URL to which the webhooks should be sent", "", arg_domain+"/api/repo/webhooks")
		kuberiUIConfig.Spec.Kubero.WebhookURL = webhookURL

		if clusterType == "" {
			clusterType = promptLine("Which cluster type have you insalled?", clusterTypeSelection, "")
		}

		if clusterType == "linode" ||
			clusterType == "digitalocean" ||
			clusterType == "gke" {
			kuberiUIConfig.Spec.Ingress.ClassName = "nginx"
		}

		kuberiUIYaml, _ := yaml.Marshal(kuberiUIConfig)
		kuberiUIErr := os.WriteFile("kuberoUI.yaml", kuberiUIYaml, 0644)
		if kuberiUIErr != nil {
			fmt.Println(kuberiUIErr)
			return
		}

		_, olminstallErr := exec.Command("kubectl", "apply", "-f", "kuberoUI.yaml", "-n", "kubero").Output()
		if olminstallErr != nil {
			fmt.Println(olminstallErr)
			cfmt.Println("{{✗ Failed to run command to install Kubero UI. Rerun installer to finish installation}}::red")
			return
		} else {
			e := os.Remove("kuberoUI.yaml")
			if e != nil {
				log.Fatal(e)
			}
			cfmt.Println("{{✓ Kubero UI installed}}::lightGreen")
		}

		kuberoUISpinner := spinner.New("Wait for Kubero UI to be ready")
		time.Sleep(8 * time.Second) //linide needs a bit more time to get the ingress up
		kuberoUISpinner.Start("run command : kubectl wait --for=condition=available deployment/kubero-sample -n kubero --timeout=180s")
		_, olmWaitErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/kubero-sample", "-n", "kubero", "--timeout=180s").Output()
		if olmWaitErr != nil {
			fmt.Println("") // keeps the spinner from overwriting the last line
			kuberoUISpinner.Error("Failed to run command. Rerun installer to finish installation")
			log.Fatal(olmWaitErr)
		}
		kuberoUISpinner.Success("Kubero UI is ready")
	}

}

func installCertManager() {
	certManagerInstalled, _ := exec.Command("kubectl", "get", "deployment", "cert-manager-webhook", "-n", "operators").Output()
	if len(certManagerInstalled) > 0 {
		cfmt.Println("{{✓ Cert Manager allready installed}}::lightGreen")
	} else {

		install := promptLine("Install SSL Certmanager", "[y,n]", "y")
		if install != "y" {
			return
		}

		certManagerSpinner := spinner.New("Install Cert Manager")
		certManagerSpinner.Start("run command : kubectl create -f https://operatorhub.io/install/cert-manager.yaml")
		_, certManagerErr := exec.Command("kubectl", "create", "-f", "https://operatorhub.io/install/cert-manager.yaml").Output()
		if certManagerErr != nil {
			fmt.Println("") // keeps the spinner from overwriting the last line
			certManagerSpinner.Error("Failed to run command. Try runnig it manually")
			log.Fatal(certManagerErr)
		}
		certManagerSpinner.Success("Cert Manager installed")

		time.Sleep(2 * time.Second)
		certManagerSpinner = spinner.New("Wait for Cert Manager to be ready")
		certManagerSpinner.Start("run command : kubectl wait --for=condition=available deployment/cert-manager-webhook -n cert-manager --timeout=180s -n operators")
		_, certManagerWaitErr := exec.Command("kubectl", "wait", "--for=condition=available", "deployment/cert-manager-webhook", "-n", "cert-manager", "--timeout=180s", "-n", "operators").Output()
		if certManagerWaitErr != nil {
			fmt.Println("") // keeps the spinner from overwriting the last line
			certManagerSpinner.Error("Failed to run command. Try runnig it manually")
			log.Fatal(certManagerWaitErr)
		}
		certManagerSpinner.Success("Cert Manager is ready")
	}
}

func writeCLIconfig() {

	ingressInstall := promptLine("Generate CLI config", "[y,n]", "y")
	if ingressInstall != "y" {
		return
	}

	//TODO consider using SSL here.
	url := promptLine("Kubero Host adress", "", "http://"+arg_domain+":"+arg_port)
	viper.Set("api.url", url)

	token := promptLine("Kubero Token", "", arg_apiToken)
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
		len(kuberoIngress.Items[0].Status.LoadBalancer.Ingress) > 0 {
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
	https://github.com/kubero-dev/kubero/wiki
	`)

	if arg_domain != "" && arg_port != "" && arg_apiToken != "" && arg_adminPassword != "" {
		cfmt.Println(`
	Your Kubero UI :{{
	URL : ` + arg_domain + `:` + arg_port + `
	User: ` + arg_adminUser + `
	Pass: ` + arg_adminPassword + `}}::lightBlue
	`)
	}
}

func generatePassword(length int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!+?._-%")
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
