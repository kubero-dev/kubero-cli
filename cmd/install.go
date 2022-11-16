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
	"github.com/leaanthony/spinner"
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

		checkAllBinaries()

		kind := promptLine("Start a local kubernetes kind cluster", "[y,n]", "n")
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

func checkAllBinaries() {
	if !checkBinary("kubectl") {
		cfmt.Println("{{✗ kubectl is not installed}}::red")
	} else {
		cfmt.Println("{{✓ kubectl is installed}}::green")
	}

	if !checkBinary("operator-sdk") {
		cfmt.Println("{{✗ operator-sdk is not installed}}::red")
	} else {
		cfmt.Println("{{✓ operator-sdk is installed}}::green")
	}

	if !checkBinary("kind") {
		cfmt.Println("{{⚠ kind is not installed}}::yellow (only required if you want to install a local kind cluster)")
	} else {
		cfmt.Println("{{✓ kind is installed}}::green")
	}
}

func checkBinary(binary string) bool {
	_, err := exec.LookPath(binary)
	return err == nil
}

func installKind() {

	if !checkBinary("kindooo") {
		log.Fatal("kind is not installed")
	}

	installer := resty.New()
	// TODO : installing the binaries needs to respect the OS and architecture
	//	installer.SetBaseURL("https://kind.sigs.k8s.io")
	//	installer.R().Get("/dl/v0.17.0/kind-linux-amd64")

	installer.SetBaseURL("https://raw.githubusercontent.com")
	kindConfig, _ := installer.R().Get("/kubero-dev/kubero/main/kind.yaml")

	fmt.Println("-------------- kind.yaml ---------------")
	fmt.Println(kindConfig.String())
	fmt.Println("----------------------------------------")
	/*
		kindConfigErr := os.WriteFile("kind.yaml", kindConfig.Body(), 0644)
		if kindConfigErr != nil {
			fmt.Println(kindConfigErr)
			return
		}
	*/

	kindSpinner := spinner.New("Spin up a local Kind cluster")
	kindSpinner.Start("run command : kind create cluster --config kind.yaml")
	out, err := exec.Command("kind", "create", "cluster", "--config", "kind.yaml").Output()
	if err != nil {
		kindSpinner.Error("Failed to run command. Try runnig it manually and skip this step")
		log.Fatal(err)
	}
	kindSpinner.Success("Kind cluster started sucessfully")

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
