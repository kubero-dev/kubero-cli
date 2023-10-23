/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package kuberoCli

import (
	_ "embed"
	"os"
	"os/exec"
	"runtime"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/cobra"
)

// debugCmd represents the debug command
var debugCmd = &cobra.Command{
	Use:     "debug",
	Aliases: []string{"dbg"},
	Short:   "Print debug informations",
	Long: `This command will print debug informations like:
	- Kubero CLI version
	- OS/Arch
	- Kubernetes version
	- Kuberop operator version
	- Kuberop operator namespace
	- Kubernetes metrics server version
	- Kubernetes cert-manager version`,
	Run: func(cmd *cobra.Command, args []string) {
		cfmt.Println("{{Kubero CLI}}::bold")
		printCLIVersion()
		printOsArch()
		cfmt.Println("\n{{Kubernetes}}::bold")
		printKubernetesVersion()
		cfmt.Println("{{Kubero Operator}}::bold")
		checkKuberoOperator()
		cfmt.Println("{{\nKubero UI}}::bold")
		checkKuberoUI()
		cfmt.Println("{{\nCert Manager}}::bold")
		checkCertManager()
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)
}

func printCLIVersion() {
	cfmt.Println("kuberoCLIVersion: ", kuberoCliVersion)
}

func printOsArch() {
	cfmt.Println("OS: ", runtime.GOOS)
	cfmt.Println("Arch: ", runtime.GOARCH)
	cfmt.Println("goVersion: ", runtime.Version())
}

func printKubernetesVersion() {
	hasKubectl := checkBinary("kubectl")
	if !hasKubectl {
		promptWarning("kubectl is not installed. Installer won't be able to install kubero. Please install kubectl and try again.")
		os.Exit(1)
	}
	version, _ := exec.Command("kubectl", "version", "-o", "yaml").Output()
	cfmt.Println(string(version))
}

func checkKuberoOperator() {
	cmdOut, _ := exec.Command("kubectl", "get", "deployments.apps", "-n", "kubero-operator-system").Output()
	cfmt.Print(string(cmdOut))

	cfmt.Println("{{\nKubero Operator Image}}::bold")
	cmdOut, _ = exec.Command("kubectl", "get", "deployment", "kubero-operator-controller-manager", "-o=jsonpath={$.spec.template.spec.containers[:1].image}", "-n", "kubero-operator-system").Output()
	cfmt.Print(string(cmdOut))
	cfmt.Println("")
}

func checkKuberoUI() {
	cmdOut, _ := exec.Command("kubectl", "get", "deployments.apps", "-n", "kubero").Output()
	cfmt.Print(string(cmdOut))

	cfmt.Println("{{\nKubero UI Ingress}}::bold")
	cmdOut, _ = exec.Command("kubectl", "get", "ingress", "-n", "kubero").Output()
	cfmt.Print(string(cmdOut))

	cfmt.Println("{{\nKubero UI Secrets}}::bold")
	cmdOut, _ = exec.Command("kubectl", "get", "secrets", "-n", "kubero").Output()
	cfmt.Print(string(cmdOut))

	cfmt.Println("{{\nKubero UI Image}}::bold")
	cmdOut, _ = exec.Command("kubectl", "get", "deployment", "kubero", "-o=jsonpath={$.spec.template.spec.containers[:1].image}", "-n", "kubero").Output()
	cfmt.Print(string(cmdOut))
	cfmt.Println("")
}

func checkCertManager() {
	cmdOut, _ := exec.Command("kubectl", "get", "deployments.apps", "-n", "cert-manager").Output()
	cfmt.Print(string(cmdOut))

	cfmt.Println("{{\nCert Manager Cluster Issuers}}::bold")
	cmdOut, _ = exec.Command("kubectl", "get", "clusterissuers.cert-manager.io").Output()
	cfmt.Print(string(cmdOut))
}
