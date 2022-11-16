/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/go-resty/resty/v2"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install all required components for kubero",
	Long: `This command will try to install all required components for kubero.

required binaries:
 - kubectl
 - operator-sdk
 - kind (optional)`,
	Run: func(cmd *cobra.Command, args []string) {
		kind := promptLine("Install a local kubernetes kind cluster", "[y,n]", "n")
		if kind == "y" {
			installKind()
		}
		/*
			olm := promptLine("Install OLM", "[y,n]", "y")
			if olm == "y" {
				installOLM()
			}

			ingress := promptLine("Install ingress-nginx", "[y,n]", "y")
			if ingress == "y" {
				installIngressNginx()
			}

			kuberoOperator := promptLine("Install kubero-operator", "[y,n]", "y")
			if kuberoOperator == "y" {
				installKuberoOperator()
			}

			kuberoUi := promptLine("Install kubero-ui", "[y,n]", "y")
			if kuberoUi == "y" {
				installKuberoUi()
			}
		*/
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func installKind() {
	fmt.Println("install kind")

	installer := resty.New()
	//	installer.SetBaseURL("https://kind.sigs.k8s.io")
	//	installer.R().Get("/dl/v0.17.0/kind-linux-amd64")

	installer.SetBaseURL("https://raw.githubusercontent.com")
	kindConfig, _ := installer.R().Get("/kubero-dev/kubero/main/kind.yaml")

	fmt.Println("----------------------------------------")
	fmt.Println("kindConfig:", kindConfig.String())
	fmt.Println("----------------------------------------")
	/*
		kindConfigErr := os.WriteFile("kind.yaml", kindConfig.Body(), 0644)
		if kindConfigErr != nil {
			fmt.Println(kindConfigErr)
			return
		}
	*/
	out, err := exec.Command("kind", "create", "cluster", "--config", "kind.yaml").Output()
	cfmt.Println("{{  run command : }}::lightWhite kind create cluster --config kind.yaml")
	if err != nil {
		cfmt.Println("{{  error : }}::red failed to run command. Try runnig it manually and skip this step")
		log.Fatal(err)
	}

	fmt.Println(string(out))

}

/*
func installOLM() {
	fmt.Println("install OLM")
}

func installIngressNginx() {
	fmt.Println("install ingress-nginx")
}

func installKuberoOperator() {
	fmt.Println("install kubero-operator")
}

func installKuberoUi() {
	fmt.Println("install kubero-ui")
}
*/
