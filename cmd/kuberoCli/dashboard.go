package kuberoCli

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"log"
)

// dashboardCmd represents the dashboard command
var dashboardCmd = &cobra.Command{
	Use:     "dashboard",
	Aliases: []string{"db"},
	Short:   "Opens the Kubero dashboard in your browser",
	Long:    `Use the dashboard subcommand to open the Kubero dashboard in your browser.`,
	Run: func(cmd *cobra.Command, args []string) {

		url := currentInstance.ApiUrl
		openURLErr := browser.OpenURL(url)
		if openURLErr != nil {
			log.Fatal("Failed to open the browser:", openURLErr)
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}
