package kuberoCli

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/leaanthony/spinner"
	"gopkg.in/yaml.v3"
)

func installKind() {

	if !checkBinary("kind") {
		log.Fatal("kind binary is not installed")
	}

	installer := resty.New()

	// get kind.yaml
	installer.SetBaseURL("https://raw.githubusercontent.com")
	kf, _ := installer.R().Get("/kubero-dev/kubero/main/kind.yaml")

	var kindConfig KindConfig
	yamlUnmarshalErr := yaml.Unmarshal(kf.Body(), &kindConfig)
	if yamlUnmarshalErr != nil {
		fmt.Println(yamlUnmarshalErr)
		return
	}

	// set cluster name
	kindConfig.Name = promptLine("Kind Cluster Name", "", "kubero-"+strconv.Itoa(rand.Intn(1000)))

	// select Kubernetes version
	kv, _ := installer.R().Get("/kubero-dev/kubero-cli/main/templates/kindVersions.yaml")
	var kindDefaults struct {
		AvailableKubernetesVersions []string `yaml:"availableKubernetesVersions"`
	}
	yamlUnmarshalBErr := yaml.Unmarshal(kv.Body(), &kindDefaults)
	if yamlUnmarshalBErr != nil {
		fmt.Println(yamlUnmarshalBErr)
		return
	}
	version := selectFromList("Kubernetes Version", kindDefaults.AvailableKubernetesVersions, "")

	kindConfig.Nodes[0].Image = "kindest/node:" + version

	if argPort == "" {
		argPort = promptLine("Local HTTP Port", "", "80")
	}
	kindConfig.Nodes[0].ExtraPortMappings[0].HostPort, _ = strconv.Atoi(argPort)

	if argPortSecure == "" {
		argPortSecure = promptLine("Local HTTPS Port", "", "443")
	}
	kindConfig.Nodes[0].ExtraPortMappings[1].HostPort, _ = strconv.Atoi(argPortSecure)

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
	_, _ = cfmt.Println("run command : kind create cluster --config kind.yaml")
	kindSpinner.Start("Creating Kind cluster")
	out, err := exec.Command("kind", "create", "cluster", "--config", "kind.yaml").Output()
	if err != nil {
		kindSpinner.Error("Failed to run command. Try running this command manually and skip this step : 'kind create cluster --config kind.yaml'")
		log.Fatal(err)
	}
	kindSpinner.Success("Kind cluster started successfully")

	fmt.Println(string(out))
}
