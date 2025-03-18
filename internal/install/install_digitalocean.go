package install

import (
	"encoding/json"
	"fmt"
	l "github.com/kubero-dev/kubero-cli/internal/log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/i582/cfmt/cmd/cfmt"
	v "github.com/kubero-dev/kubero-cli/version"
	"github.com/leaanthony/spinner"
)

var KuberoCliVersion = v.Version()

func (m *ManagerInstall) installDigitalOcean() error {
	// https://docs.digitalocean.com/reference/api/api-reference/#operation/kubernetes_create_cluster

	l.Warn("Installing Kubernetes on Digital Ocean is currently beta state in kubero-cli")
	l.Warn("Please report if you run into errors")

	token := os.Getenv("DIGITALOCEAN_ACCESS_TOKEN")
	if token == "" {
		l.Error("missing DIGITALOCEAN_ACCESS_TOKEN")
		return os.ErrNotExist
	}

	doApi := resty.New().
		SetAuthScheme("Bearer").
		SetAuthToken(token).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "kubero-cli/"+KuberoCliVersion).
		SetBaseURL("https://api.digitalocean.com")

	var doConfig DigitalOceanKubernetesConfig
	doConfig.NodePools = []struct {
		Size  string `json:"size" gorm:"column:size"`
		Count int    `json:"count" gorm:"column:count"`
		Name  string `json:"name" gorm:"column:name"`
	}([]struct {
		Size  string `json:"size"`
		Count int    `json:"count"`
		Name  string `json:"name"`
	}{
		{
			Size:  "s-1vcpu-2gb",
			Count: 1,
			Name:  "worker-pool",
		},
	})

	kubernetesVersions, regions, sizes := m.getDigitaloceanOptions(doApi)

	doConfig.Name = promptLine("Kubernetes Cluster Name", "", "kubero-"+strconv.Itoa(rand.Intn(1000)))
	doConfig.Region = selectFromList("Cluster Region", regions, "")
	doConfig.Version = selectFromList("Cluster Version", kubernetesVersions, "")

	doConfig.NodePools[0].Size = selectFromList("Cluster Node Size", sizes, "")
	doConfig.NodePools[0].Count, _ = strconv.Atoi(promptLine("Cluster Node Count", "", "3"))

	kf, _ := doApi.R().
		SetBody(doConfig).
		Post("/v2/kubernetes/clusters")

	if kf.StatusCode() > 299 {
		fmt.Println(kf.String())
		l.Error("failed to create digital ocean cluster")
		return fmt.Errorf("failed to create digital ocean cluster")
	} else {
		l.Info("digital ocean cluster created")
	}

	var doCluster DigitalOcean
	jsonUnmarshalErr := json.Unmarshal(kf.Body(), &doCluster)
	if jsonUnmarshalErr != nil {
		l.Error("failed to unmarshal digital ocean cluster response")
		return jsonUnmarshalErr
	}

	doSpinner := spinner.New("Starting a kubernetes cluster on digital ocean")
	doSpinner.Start("Waiting for digital ocean cluster to be ready. This may take a few minutes. Time enough to get a coffee ☕")
	clusterID := doCluster.KubernetesCluster.ID

	for {
		time.Sleep(2 * time.Second)
		doWait, _ := doApi.R().
			Get("/v2/kubernetes/clusters/" + clusterID)

		if doWait.StatusCode() > 299 {
			l.Warn(doWait.String())
			doSpinner.Error("Failed to create digital ocean cluster")
			continue
		} else {
			var doCluster DigitalOcean
			jsonUnmarshalBErr := json.Unmarshal(doWait.Body(), &doCluster)
			if jsonUnmarshalBErr != nil {
				l.Error("failed to unmarshal digital ocean cluster response")
				return jsonUnmarshalBErr
			}
			if doCluster.KubernetesCluster.Status.State == "running" {
				doSpinner.Success("digital ocean cluster created")
				break
			}
		}
	}

	kubectl, _ := doApi.R().
		Get("v2/kubernetes/clusters/" + clusterID + "/kubeconfig")
	mergeKubeconfigErr := utils.MergeKubeconfig(kubectl.Body())
	if mergeKubeconfigErr != nil {
		l.Error("failed to merge kubeconfig")
		return mergeKubeconfigErr
	}

	return nil
}

func (m *ManagerInstall) getDigitaloceanOptions(api *resty.Client) ([]string, []string, []string) {
	token := os.Getenv("DIGITALOCEAN_ACCESS_TOKEN")
	if token == "" {
		_, _ = cfmt.Println("{{✗ DIGITALOCEAN_ACCESS_TOKEN is not set}}::red")
		l.Fatal("missing DIGITALOCEAN_ACCESS_TOKEN")
	}

	optionsResponse, err := api.R().Get("/v2/kubernetes/options")
	if err != nil {
		_, _ = cfmt.Println("{{✗ failed to get digitalocean options}}::red")
		l.Fatal(err)
	}

	var versionsResponse DigitaloceanOptions
	jsonUnmarshalErr := json.Unmarshal(optionsResponse.Body(), &versionsResponse)
	if jsonUnmarshalErr != nil {
		fmt.Println(jsonUnmarshalErr)
		return nil, nil, nil
	}

	var versions []string
	for _, v := range versionsResponse.Options.Versions {
		versions = append(versions, v.Slug)
	}

	var regions []string
	for _, v := range versionsResponse.Options.Regions {
		regions = append(regions, v.Slug)
	}

	var sizes []string
	for _, v := range versionsResponse.Options.Sizes {
		sizes = append(sizes, v.Slug)
	}

	return versions, regions, sizes
}
