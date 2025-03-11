package utils

import (
	"fmt"
	"math/rand"
	"os/exec"
	"time"
)

func GenerateRandomString(length int, chars string) string {
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

func CheckBinary(binary string) bool {
	_, err := exec.LookPath(binary)
	return err == nil
}

func CheckAllBinaries(binaries ...string) error {
	for _, binary := range binaries {
		if !CheckBinary(binary) {
			return fmt.Errorf("binary %s not found", binary)
		}
	}
	return nil
}

func CheckClusters() error {
	return CheckCluster("kind")
}

func CheckCluster(clusterType string) error {
	if clusterType == "kind" {
		return nil
	} else {
		return checkKubeConfig()
	}
}

func checkKubeConfig() error {
	if !CheckBinary("kubectl") {
		return fmt.Errorf("kubectl not found in PATH")
	}

	cmd := exec.Command("kubectl", "cluster-info")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("kubectl cluster-info failed: %s", err)
	}
	return nil
}

func CreateNamespace(namespace string) error {
	_, err := exec.Command("kubectl", "create", "namespace", namespace).Output()
	return err
}

func DeleteNamespace(namespace string) error {
	_, err := exec.Command("kubectl", "delete", "namespace", namespace).Output()
	return err
}

func CheckNamespace(namespace string) error {
	_, err := exec.Command("kubectl", "get", "namespace", namespace).Output()
	return err
}
