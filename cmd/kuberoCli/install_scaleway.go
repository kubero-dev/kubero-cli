package kuberoCli

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strconv"
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

	// _, _ = cfmt.Println("{{⚠ Installing Kubernetes on Scaleway is currently beta state in kubero-cli}}::yellow")
	// _, _ = cfmt.Println("{{  Please report if you run into errors}}::yellow")

	var cluster ScalewayCreate
	/*
		cluster.ProjectID = os.Getenv("SCALEWAY_PROJECTID")
		if cluster.ProjectID == "" {
			_, _ = cfmt.Println("{{✗ SCALEWAY_PROJECTID is not set}}::red")
			log.Fatal("missing SCALEWAY_PROJECTID")
		}
	*/
	cluster.OrganizationID = os.Getenv("SCALEWAY_PROJECTID")

	token := os.Getenv("SCALEWAY_ACCESS_TOKEN")
	if token == "" {
		_, _ = cfmt.Println("{{✗ SCALEWAY_ACCESS_TOKEN is not set}}::red")
		log.Fatal("missing SCALEWAY_ACCESS_TOKEN")
	}

	api := resty.New().
		SetHeader("X-Auth-Token", token).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "kubero-cli/"+kuberoCliVersion).
		SetBaseURL("https://api.scaleway.com/k8s/v1/regions")

	cluster.Name = promptLine("Kubernetes Cluster Name", "", "kubero-"+strconv.Itoa(rand.Intn(1000)))

	regionsList := []string{"fr-par", "nl-ams", "pl-waw"}
	region := selectFromList("Cluster Region", regionsList, "")

	versions := getScalewayVersions(api, region)
	cluster.Version = selectFromList("Kubernetes Version", versions, "")

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
	nodeType := promptLine("Node Types", "[DEV1-M,DEV1-XL,START1-M]", "DEV1-M")

	clusterSize, _ := strconv.Atoi(promptLine("Cluster Size", "[at least 3]", "3"))

	cluster.Pools = append(cluster.Pools, ScalewayNodePool{
		Name:             "default",
		NodeType:         nodeType,
		Autoscaling:      false,
		Size:             clusterSize,
		ContainerRuntime: "unknown_runtime",
		RootVolumeType:   "default_volume_type",
		//RootVolumeSize:   50,
	})

	//fmt.Printf("%+v\n", cluster)
	newCluster, _ := api.R().SetBody(cluster).Post(region + "/clusters")

	var clusterResponse ScalewayCreateResponse
	if newCluster.StatusCode() < 299 {
		jsonUnmarshalErr := json.Unmarshal(newCluster.Body(), &clusterResponse)
		if jsonUnmarshalErr != nil {
			_, _ = cfmt.Println("{{✗ Failed to create Scaleway Kubernetes cluster}}::red")
			return
		}
		_, _ = cfmt.Println("{{✓ Scaleway Kubernetes cluster created}}::lightGreen")
	} else {
		_, _ = cfmt.Println("{{✗ Scaleway Kubernetes Cluster creation failed}}::red")
		log.Fatal(string(newCluster.Body()))
	}

	spinnerObj := spinner.New()
	spinnerObj.Start("Waiting for cluster to be ready")
	for {
		clusterStatus, _ := api.R().Get(region + "/clusters/" + clusterResponse.ID)
		var clusterStatusResponse ScalewayCreateResponse
		jsonUnmarshalErr := json.Unmarshal(clusterStatus.Body(), &clusterStatusResponse)
		if jsonUnmarshalErr != nil {
			_, _ = cfmt.Println("{{✗ Failed to get Scaleway Kubernetes cluster status}}::red")
			return
		}
		if clusterStatusResponse.Status == "ready" {
			spinnerObj.Success("Scaleway Kubernetes Cluster is ready")
			break
		}
		time.Sleep(2 * time.Second)
	}

	k, _ := api.R().Get(region + "/clusters/" + clusterResponse.ID + "/kubeconfig")

	var scalewayKubeconfigResponse ScalewayKubeconfigResponse
	jsonUnmarshalErr := json.Unmarshal(k.Body(), &scalewayKubeconfigResponse)
	if jsonUnmarshalErr != nil {
		_, _ = cfmt.Println("{{✗ Failed to download kubeconfig}}::red")
		return
	}
	kubeconfig, _ := base64.StdEncoding.DecodeString(scalewayKubeconfigResponse.Content)

	err := mergeKubeconfig(kubeconfig)
	if err != nil {
		_, _ = cfmt.Println("{{✗ Failed to download kubeconfig}}::red")
		log.Fatal(err)
	} else {
		_, _ = cfmt.Println("{{✓ Kubeconfig downloaded}}::lightGreen")
	}

}

func getScalewayVersions(api *resty.Client, region string) []string {
	token := os.Getenv("SCALEWAY_ACCESS_TOKEN")
	if token == "" {
		_, _ = cfmt.Println("{{✗ SCALEWAY_ACCESS_TOKEN is not set}}::red")
		log.Fatal("missing SCALEWAY_ACCESS_TOKEN")
	}

	versionsR, _ := api.R().Get(region + "/versions")

	var versionsResponse ScalewayVersionsResponse
	jsonUnmarshalErr := json.Unmarshal(versionsR.Body(), &versionsResponse)
	if jsonUnmarshalErr != nil {
		_, _ = cfmt.Println("{{✗ Failed to get Scaleway Kubernetes versions}}::red")
		return nil
	}

	var versions []string
	for _, v := range versionsResponse.Versions {
		versions = append(versions, v.Name)
	}

	return versions
}
