/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package kuberoCli

import (
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/cobra"
)

// dashboardCmd represents the dashboard command
var dashboardCmd = &cobra.Command{
	Use:     "dashboard",
	Aliases: []string{"db"},
	Short:   "Opens the Kubero dashboard in your browser",
	Long:    `Use the dashboard subcommand to open the Kubero dashboard in your browser.`,
	Run: func(cmd *cobra.Command, args []string) {

		ooo := getIACBaseDir()
		cfmt.Println(ooo)
		/*
			url := currentInstance.Apiurl
			browser.OpenURL(url)
		*/
	},
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}
