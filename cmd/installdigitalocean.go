package cmd

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
	doConfig.Version = promptLine("Cluster Version", "[1.25.4-do.0,1.24.8-do.0]", "1.25.4-do.0")

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
