package kuberoCli

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/leaanthony/spinner"
)

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
		SetHeader("User-Agent", "kubero-cli/"+kuberoCliVersion).
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

	kubernetesVersions, regions, sizes := getDigitaloceanOptions(doApi)

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

func getDigitaloceanOptions(api *resty.Client) ([]string, []string, []string) {
	token := os.Getenv("DIGITALOCEAN_ACCESS_TOKEN")
	if token == "" {
		cfmt.Println("{{✗ DIGITALOCEAN_ACCESS_TOKEN is not set}}::red")
		log.Fatal("missing DIGITALOCEAN_ACCESS_TOKEN")
	}

	optionsResponse, err := api.R().Get("/v2/kubernetes/options")
	if err != nil {
		cfmt.Println("{{✗ failed to get digitalocean options}}::red")
		log.Fatal(err)
	}

	var versionsResponse DigitaloceanOptions
	json.Unmarshal(optionsResponse.Body(), &versionsResponse)

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
