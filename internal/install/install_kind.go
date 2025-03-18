package install

import (
	"github.com/go-resty/resty/v2"
	l "github.com/kubero-dev/kubero-cli/internal/log"
	"github.com/leaanthony/spinner"
	"gopkg.in/yaml.v3"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
)

func (m *ManagerInstall) installKind() error {
	if !utils.CheckBinary("kind") {
		l.Error("kind binary not found. Please install kind and try again")
		return os.ErrNotExist
	}

	installer := resty.New()

	// get kind.yaml
	installer.SetBaseURL("https://raw.githubusercontent.com")
	kf, _ := installer.R().Get("/kubero-dev/kubero/main/kind.yaml")

	var kindConfig KindConfig
	yamlUnmarshalErr := yaml.Unmarshal(kf.Body(), &kindConfig)
	if yamlUnmarshalErr != nil {
		l.Error("Failed to unmarshal kind.yaml")
		return yamlUnmarshalErr
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
		l.Error("Failed to unmarshal kindVersions.yaml")
		return yamlUnmarshalBErr
	}
	version := selectFromList("Kubernetes Version", kindDefaults.AvailableKubernetesVersions, "")

	kindConfig.Nodes[0].Image = "kindest/node:" + version

	if m.argPort == "" {
		m.argPort = promptLine("Local HTTP Port", "", "80")
	}
	kindConfig.Nodes[0].ExtraPortMappings[0].HostPort, _ = strconv.Atoi(m.argPort)

	if m.argPortSecure == "" {
		m.argPortSecure = promptLine("Local HTTPS Port", "", "443")
	}
	kindConfig.Nodes[0].ExtraPortMappings[1].HostPort, _ = strconv.Atoi(m.argPortSecure)

	kindConfigYaml, _ := yaml.Marshal(kindConfig)

	kindConfigErr := os.WriteFile("kind.yaml", kindConfigYaml, 0644)
	if kindConfigErr != nil {
		l.Error("Failed to write kind.yaml")
		return kindConfigErr
	}

	kindSpinner := spinner.New("Spin up a local Kind cluster")
	l.Info("run command : kind create cluster --config kind.yaml")
	kindSpinner.Start("Creating Kind cluster")
	out, err := exec.Command("kind", "create", "cluster", "--config", "kind.yaml").Output()
	if err != nil {
		kindSpinner.Error("Failed to run command. Try running this command manually and skip this step : 'kind create cluster --config kind.yaml'")
		return err
	}
	kindSpinner.Success("Kind cluster started successfully")

	l.Info(string(out))

	return nil
}
