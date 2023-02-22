package cmd

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strconv"

	"github.com/leaanthony/spinner"
)

func installGKE() {
	// implememted with gcloud, since it is required for the download of the kubeconfig anyway

	// gcloud config list
	// gcloud config get project
	// gcloud container clusters create kubero-cluster-4 --region=us-central1-c
	// gcloud container clusters get-credentials kubero-cluster-4 --region=us-central1-c

	// https://cloud.google.com/kubernetes-engine/docs/reference/libraries#client-libraries-install-go
	// https://github.com/googleapis/google-cloud-go

	gcloudName := promptLine("Kubernetes Cluster Name", "", "kubero-"+strconv.Itoa(rand.Intn(1000)))
	gcloudRegion := promptLine("Region", "[https://cloud.google.com/compute/docs/regions-zones]", "us-central1-c")
	gcloudClusterVersion := promptLine("Cluster Version", "[https://cloud.google.com/kubernetes-engine/docs/release-notes-regular]", "1.23.8-gke.1900")

	spinner := spinner.New("Spin up a GKE cluster")
	spinner.Start("run command : gcloud container clusters create " + gcloudName + " --region=" + gcloudRegion + " --cluster-version=" + gcloudClusterVersion)
	_, err := exec.Command("gcloud", "container", "clusters", "create", gcloudName,
		"--region="+gcloudRegion,
		"--cluster-version="+gcloudClusterVersion).Output()
	if err != nil {
		fmt.Println()
		spinner.Error("Failed to run command. Try runnig this command manually and skip this step: 'gcloud container clusters create " + gcloudName + " --region=" + gcloudRegion + " --cluster-version=" + gcloudClusterVersion + "'")
		log.Fatal(err)
	}
	spinner.Success("GKE cluster started sucessfully")

	spinner.Start("Get credentials for the GKE cluster")
	_, err = exec.Command("gcloud", "container", "clusters", "get-credentials", gcloudName, "--region="+gcloudRegion).Output()
	if err != nil {
		fmt.Println()
		spinner.Error("Failed to run command. Try runnig this command manually and skip this step: 'gcloud container clusters get-credentials " + gcloudName + " --region=" + gcloudRegion + "'")
		log.Fatal(err)
	} else {
		spinner.Success("GKE cluster credentials set")
	}

}
