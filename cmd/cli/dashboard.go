package cli

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"github.com/faelmori/kubero-cli/cmd/common"
	"github.com/faelmori/kubero-cli/internal/config"
	"github.com/faelmori/kubero-cli/internal/log"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"path/filepath"
)

func DashboardCmds() []*cobra.Command {
	return []*cobra.Command{
		cmdDashboard(),
	}
}

func cmdDashboard() *cobra.Command {
	var configPath string

	dashboardCmd := &cobra.Command{
		Use:     "dashboard",
		Aliases: []string{"db"},
		Short:   "Opens the Kubero dashboard in your browser",
		Long:    `Use the dashboard subcommand to open the Kubero dashboard in your browser.`,
		Annotations: common.GetDescriptions([]string{
			"Opens the Kubero dashboard in your browser",
			`Use the dashboard subcommand to open the Kubero dashboard in your browser.`,
		}, false),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := filepath.Dir(configPath)
			file := filepath.Base(configPath)
			cfgMgr := config.NewViperConfig(path, file)
			if cfgMgrErr := cfgMgr.LoadConfig(); cfgMgrErr != nil {
				return cfgMgrErr
			}
			url := cfgMgr.GetConfig().Api.Url
			openURLErr := browser.OpenURL(url)
			if openURLErr != nil {
				return openURLErr
			}

			log.Info("Opening the Kubero dashboard in your browser...")

			return nil
		},
	}

	dashboardCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to the Kubero dashboard configuration file")

	return dashboardCmd
}
