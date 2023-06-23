package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/leaanthony/spinner"
)

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
	/*
		cluster.ProjectID = os.Getenv("SCALEWAY_PROJECTID")
		if cluster.ProjectID == "" {
			cfmt.Println("{{✗ SCALEWAY_PROJECTID is not set}}::red")
			log.Fatal("missing SCALEWAY_PROJECTID")
		}
	*/
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

	versions := getScalewayVersions(api, region)
	cluster.Version = promptLine("Kubernetes Version", "["+strings.Join(versions, ",")+"]", versions[0])

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

	fmt.Printf("%+v\n", cluster)
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

func getScalewayVersions(api *resty.Client, region string) []string {
	token := os.Getenv("SCALEWAY_ACCESS_TOKEN")
	if token == "" {
		cfmt.Println("{{✗ SCALEWAY_ACCESS_TOKEN is not set}}::red")
		log.Fatal("missing SCALEWAY_ACCESS_TOKEN")
	}

	versions_r, _ := api.R().Get(region + "/versions")

	var versionsResponse ScalewayVersionsResponse
	json.Unmarshal(versions_r.Body(), &versionsResponse)

	var versions []string
	for _, v := range versionsResponse.Versions {
		versions = append(versions, v.Name)
	}

	return versions
}
