/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

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
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func installKind() {
	fmt.Println("install kind")
}

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
