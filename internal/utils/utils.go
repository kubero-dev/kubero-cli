package utils

import (
	"bytes"
	"fmt"
	"github.com/faelmori/kubero-cli/internal/log"
	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
	"k8s.io/client-go/tools/clientcmd"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

type Utils struct{}

func NewUtils() *Utils { return &Utils{} }

func (u *Utils) GenerateRandomString(length int, chars string) string {
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

func (u *Utils) CheckBinary(binary string) bool {
	envs := os.Environ()
	if ep, ok := os.LookupEnv("PATH"); ok {
		envs = append(envs, "PATH=$PATH:"+ep)
	}
	_, err := exec.LookPath(binary)
	return err == nil
}

func (u *Utils) CheckAllBinaries(binaries ...string) error {
	for _, binary := range binaries {
		if !u.CheckBinary(binary) {
			return fmt.Errorf("binary %s not found", binary)
		}
	}
	return nil
}

func (u *Utils) CheckClusters() error {
	var outb, errb bytes.Buffer
	clusterInfo := exec.Command("kubectl", "cluster-info")
	clusterInfo.Stdout = &outb
	clusterInfo.Stderr = &errb
	err := clusterInfo.Run()
	if err != nil {
		log.Error("command failed : kubectl cluster-info")
		log.Error(err.Error() + "\n" + outb.String())
		return err
	}
	log.Info("Cluster info: ", outb.String())

	out, _ := exec.Command("kubectl", "config", "get-contexts").Output()
	log.Info("Current cluster: ", string(out))

	clusterSelect := NewConsolePrompt().PromptLine("Is the CURRENT cluster the one you wish to install Kubero?", "[y,n]", "y")
	if clusterSelect == "n" {
		log.Error("Please select the correct cluster and try again")
		return fmt.Errorf("cluster not selected")
	}

	return nil
}

func (u *Utils) CheckKubeConfig() error {
	if !u.CheckBinary("kubectl") {
		return fmt.Errorf("kubectl not found in PATH")
	}

	cmd := exec.Command("kubectl", "cluster-info")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("kubectl cluster-info failed: %s", err)
	}
	return nil
}

func (u *Utils) CreateNamespace(namespace string) error {
	_, err := exec.Command("kubectl", "create", "namespace", namespace).Output()
	return err
}

func (u *Utils) DeleteNamespace(namespace string) error {
	_, err := exec.Command("kubectl", "delete", "namespace", namespace).Output()
	return err
}

func (u *Utils) CheckNamespace(namespace string) error {
	_, err := exec.Command("kubectl", "get", "namespace", namespace).Output()
	return err
}

func (u *Utils) MergeKubeconfig(kubeconfig []byte) error {
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

func (u *Utils) PrintCLI(table *tablewriter.Table, r *resty.Response, outputFormat string) {
	if outputFormat == "json" {
		fmt.Println(r)
	} else {
		table.Render()
	}
}

func (u *Utils) BoolToEmoji(b bool) string {
	if b {
		return "✅"
	}
	return "❌"
}
