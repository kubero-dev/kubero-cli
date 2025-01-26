package kuberoCli

import (
	"encoding/base64"
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

func installLinode() {
	// https://www.linode.com/docs/api/linode-kubernetes-engine-lke/#kubernetes-cluster-create
	// https://www.linode.com/docs/api/linode-kubernetes-engine-lke/#kubernetes-cluster-view
	// https://www.linode.com/docs/api/linode-kubernetes-engine-lke/#kubeconfig-view

	_, _ = cfmt.Println("{{⚠ Installing Kubernetes on Linode is currently beta state in kubero-cli}}::yellow")
	_, _ = cfmt.Println("{{  Please report if you run into errors}}::yellow")

	token := os.Getenv("LINODE_ACCESS_TOKEN")
	if token == "" {
		_, _ = cfmt.Println("{{✗ LINODE_ACCESS_TOKEN is not set}}::red")
		log.Fatal("missing LINODE_ACCESS_TOKEN")
	}

	api := resty.New().
		SetAuthScheme("Bearer").
		SetAuthToken(token).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "kubero-cli/"+kuberoCliVersion).
		SetBaseURL("https://api.linode.com/v4/lke/clusters")

	var clusterConfig LinodeCreateClusterRequest
	clusterConfig.Label = promptLine("Cluster Name", "", "kubero-"+strconv.Itoa(rand.Intn(1000)))
	clusterConfig.Region = promptLine("Region", "[https://www.linode.com/global-infrastructure/]", "us-central") // TODO load the list of regions or point to e better document

	workerNodesCount, _ := strconv.Atoi(promptLine("Worker Nodes Count", "", "3"))
	workerNodesType := promptLine("Worker Nodes Type", "[https://www.linode.com/pricing/]", "g6-standard-2") // TODO load the list of types or point to e better document

	clusterConfig.K8SVersion = promptLine("Kubernetes Version", "[1.25]", "1.25")
	clusterConfig.Tags = []string{"kubero"}
	clusterConfig.NodePools = []LinodeNodepool{
		{
			Type:  workerNodesType,
			Count: workerNodesCount,
		},
	}

	spinnerObj := spinner.New("Spin up a Linode Kubernetes Cluster")

	spinnerObj.Start("Create Linode Kubernetes Cluster")
	clusterResponse, _ := api.R().SetBody(clusterConfig).Post("")
	if clusterResponse.StatusCode() > 299 {
		fmt.Println()
		spinnerObj.Error("Failed to create Linode Kubernetes Cluster")
		log.Fatal(clusterResponse.String())
	}
	spinnerObj.Success("Linode Kubernetes Cluster created")

	var cluster LinodeCreateClusterResponse
	jsonUnmarshalErr := json.Unmarshal(clusterResponse.Body(), &cluster)
	if jsonUnmarshalErr != nil {
		log.Fatal(jsonUnmarshalErr)
		return
	}

	// According to the docs, the cluster is ready after 2-5 minutes.
	_, _ = cfmt.Println("{{  Wait for Linode Kubernetes Cluster to be ready}}::lightBlue")
	_, _ = cfmt.Println("{{  According to the docs this may take up to 7 minutes}}::lightBlue")
	_, _ = cfmt.Println("{{  Time for a coffee break and some Chuck Norris jokes.}}::lightBlue")
	spinnerObj.Start("Wait for Linode Kubernetes Cluster to be ready")

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
			spinnerObj.Success("Linode Kubernetes Cluster is ready")
			break
		}
	}
	kubeconfig, err := base64.StdEncoding.DecodeString(LinodeKubeconfig.Kubeconfig)

	if err != nil {
		fmt.Println()
		spinnerObj.Error("Failed to decode kubeconfig")
		log.Fatal(err)
	}

	err = mergeKubeconfig(kubeconfig)
	if err != nil {
		fmt.Println()
		spinnerObj.Error("Failed to merge kubeconfig")
		log.Fatal(err)
	}

	spinnerObj.Success("Linode Kubernetes Cluster credentials set")

}
