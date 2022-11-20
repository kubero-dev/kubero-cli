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
	"github.com/spf13/viper"
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
		writeCLIconffig()
		finalMessage()
	},
}

var arg_adminPassword string
var arg_adminUser string
var arg_domain string
var arg_apiToken string
var arg_port string
var arg_portSecure string

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
}

func checkBinary(binary string) bool {
	_, err := exec.LookPath(binary)
	return err == nil
}

func installGKE() {
	// TODO
	// variables:
	// - cluster name
	// - cluster region
	// gcloud config list
	// gcloud config get project
	// gcloud container clusters get-credentials kubero-cluster-4 --region=us-central1-c
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

	installer.SetBaseURL("https://raw.githubusercontent.com")
	kf, _ := installer.R().Get("/kubero-dev/kubero/main/kind.yaml")

	var kindConfig KindConfig
	yaml.Unmarshal(kf.Body(), &kindConfig)

	kindConfig.Name = promptLine("Kind Cluster Name", "", "kubero-"+strconv.Itoa(rand.Intn(1000)))
	kindConfig.Nodes[0].Image = "kindest/node:v1.25.3"

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

	olmCRDInstalled, _ := exec.Command("kubectl", "get", "crd", "subscriptions.operators.coreos.com").Output()
	if len(olmCRDInstalled) > 0 {
		cfmt.Println("{{✓ OLM CRD's allredy installed}}::lightGreen")
	} else {
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

		if githubPersonalAccessToken != "" {
			createSecretCommand.Args = append(createSecretCommand.Args, "--from-literal=KUBERO_GITHUB_PERSONAL_ACCESS_TOKEN="+githubPersonalAccessToken)
		}
		if giteaPersonalAccessToken != "" {
			createSecretCommand.Args = append(createSecretCommand.Args, "--from-literal=KUBERO_GITEA_PERSONAL_ACCESS_TOKEN="+giteaPersonalAccessToken)
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
		var outb, errb bytes.Buffer

		installer := resty.New()

		installer.SetBaseURL("https://raw.githubusercontent.com")
		kf, _ := installer.R().Get("kubero-dev/kubero-operator/main/config/samples/application_v1alpha1_kubero.yaml")

		var kuberiUIConfig KuberoUIConfig
		yaml.Unmarshal(kf.Body(), &kuberiUIConfig)

		if arg_domain == "" {
			arg_domain = promptLine("Kuberi UI Domain", "", "kubero.lacolhost.com")
		}
		kuberiUIConfig.Spec.Ingress.Hosts[0].Host = arg_domain

		kuberiUIYaml, _ := yaml.Marshal(kuberiUIConfig)
		kuberiUIErr := os.WriteFile("kuberoUI.yaml", kuberiUIYaml, 0644)
		if kuberiUIErr != nil {
			fmt.Println(kuberiUIErr)
			return
		}

		kuberoUI := exec.Command("kubectl", "apply", "-f", "kuberoUI.yaml", "-n", "kubero")
		kuberoUI.Stdout = &outb
		kuberoUI.Stderr = &errb
		err := kuberoUI.Run()
		if err != nil {
			fmt.Println(errb.String())
			fmt.Println(outb.String())
			cfmt.Println("{{✗ Failed to run command to install Kubero UI. Try runnig it manually}}::red")
			log.Fatal()
		} else {
			e := os.Remove("kuberoUI.yaml")
			if e != nil {
				log.Fatal(e)
			}
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

func writeCLIconffig() {

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

func finalMessage() {
	cfmt.Println(`

{{⚠ make sure your DNS is pointing to your Kubernetes cluster}}::yellow

	,--. ,--.        ,--.
	|  .'   /,--.,--.|  |-.  ,---. ,--.--. ,---.
	|  .   ' |  ||  || .-. '| .-. :|  .--'| .-. |
	|  |\   \'  ''  '| '-' |\   --.|  |   ' '-' '
	'--' '--' '----'  '---'  '----''--'    '---'

Your Kubero UI :{{
  URL : ` + arg_domain + `:` + arg_port + `
  User: ` + arg_adminUser + `
  Pass: ` + arg_adminPassword + `}}::lightBlue

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

type KuberoUIConfig struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Affinity struct {
		} `yaml:"affinity"`
		FullnameOverride string `yaml:"fullnameOverride"`
		Image            struct {
			PullPolicy string `yaml:"pullPolicy"`
			Repository string `yaml:"repository"`
			Tag        string `yaml:"tag"`
		} `yaml:"image"`
		ImagePullSecrets []interface{} `yaml:"imagePullSecrets"`
		Ingress          struct {
			Annotations struct {
			} `yaml:"annotations"`
			ClassName string `yaml:"className"`
			Enabled   bool   `yaml:"enabled"`
			Hosts     []struct {
				Host  string `yaml:"host"`
				Paths []struct {
					Path     string `yaml:"path"`
					PathType string `yaml:"pathType"`
				} `yaml:"paths"`
			} `yaml:"hosts"`
			TLS []interface{} `yaml:"tls"`
		} `yaml:"ingress"`
		NameOverride string `yaml:"nameOverride"`
		NodeSelector struct {
		} `yaml:"nodeSelector"`
		PodAnnotations struct {
		} `yaml:"podAnnotations"`
		PodSecurityContext struct {
		} `yaml:"podSecurityContext"`
		ReplicaCount int `yaml:"replicaCount"`
		Resources    struct {
		} `yaml:"resources"`
		SecurityContext struct {
		} `yaml:"securityContext"`
		Service struct {
			Port int    `yaml:"port"`
			Type string `yaml:"type"`
		} `yaml:"service"`
		ServiceAccount struct {
			Annotations struct {
			} `yaml:"annotations"`
			Create bool   `yaml:"create"`
			Name   string `yaml:"name"`
		} `yaml:"serviceAccount"`
		Tolerations []interface{} `yaml:"tolerations"`
		Kubero      struct {
			Debug      string `yaml:"debug"`
			Namespace  string `yaml:"namespace"`
			Context    string `yaml:"context"`
			WebhookURL string `yaml:"webhook_url"`
			SessionKey string `yaml:"sessionKey"`
			Auth       struct {
				Github struct {
					Enabled     bool   `yaml:"enabled"`
					ID          string `yaml:"id"`
					Secret      string `yaml:"secret"`
					CallbackURL string `yaml:"callbackUrl"`
					Org         string `yaml:"org"`
				} `yaml:"github"`
				Oauth2 struct {
					Enabled     bool   `yaml:"enabled"`
					Name        string `yaml:"name"`
					ID          string `yaml:"id"`
					AuthURL     string `yaml:"authUrl"`
					TokenURL    string `yaml:"tokenUrl"`
					Secret      string `yaml:"secret"`
					CallbackURL string `yaml:"callbackUrl"`
				} `yaml:"oauth2"`
				Config string `yaml:"config"`
				Kubero struct {
					Context   string `yaml:"context"`
					Namespace string `yaml:"namespace"`
					Port      int    `yaml:"port"`
				} `yaml:"kubero"`
				Buildpacks []struct {
					Name     string `yaml:"name"`
					Language string `yaml:"language"`
					Fetch    struct {
						Repository string `yaml:"repository"`
						Tag        string `yaml:"tag"`
					} `yaml:"fetch"`
					Build struct {
						Repository string `yaml:"repository"`
						Tag        string `yaml:"tag"`
						Command    string `yaml:"command"`
					} `yaml:"build"`
					Run struct {
						Repository         string `yaml:"repository"`
						Tag                string `yaml:"tag"`
						ReadOnlyAppStorage bool   `yaml:"readOnlyAppStorage"`
						SecurityContext    struct {
							AllowPrivilegeEscalation bool `yaml:"allowPrivilegeEscalation"`
							ReadOnlyRootFilesystem   bool `yaml:"readOnlyRootFilesystem"`
						} `yaml:"securityContext"`
						Command string `yaml:"command"`
					} `yaml:"run,omitempty"`
				} `yaml:"buildpacks"`
				PodSizeList []struct {
					Name        string `yaml:"name"`
					Description string `yaml:"description"`
					Default     bool   `yaml:"default,omitempty"`
					Resources   struct {
						Requests struct {
							Memory string `yaml:"memory"`
							CPU    string `yaml:"cpu"`
						} `yaml:"requests"`
						Limits struct {
							Memory string `yaml:"memory"`
							CPU    string `yaml:"cpu"`
						} `yaml:"limits"`
					} `yaml:"resources,omitempty"`
					Active bool `yaml:"active,omitempty"`
				} `yaml:"podSizeList"`
			} `yaml:"auth"`
		} `yaml:"kubero"`
	} `yaml:"spec"`
}
