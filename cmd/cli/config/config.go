package config

import (
	"fmt"
	"github.com/faelmori/kubero-cli/internal/config"
	"github.com/spf13/cobra"
	"path/filepath"
)

func ConfigCmds() []*cobra.Command {
	return []*cobra.Command{
		cmdConfigCli(),
		cmdConfigSet(),
		cmdConfigGet(),
	}
}

func cmdConfigCli() *cobra.Command {
	var path string

	cfgCmd := &cobra.Command{
		Use:   "config",
		Short: "Show your configuration",
		Long: `Show your configuration. This command will show your current configuration.
You can use the 'config set' command to set a new configuration.`,
		Run: func(cmd *cobra.Command, args []string) {
			pt := filepath.Dir(path)
			nm := filepath.Base(path)

			cfgManager := config.NewViperConfig(pt, nm)
			cfgMap := cfgManager.GetConfig()

			fmt.Println(fmt.Sprintf("Configuration file: %s", path))
			fmt.Println("Configuration:")
			fmt.Println(cfgMap)
		},
	}

	cfgCmd.Flags().StringVarP(&path, "path", "p", "", "Path to the configuration file")

	return cfgCmd
}

func cmdConfigSet() *cobra.Command {
	var path string

	cfgCmd := &cobra.Command{
		Use:   "set",
		Short: "Set a new configuration",
		Long: `Set a new configuration. This command will set a new configuration.
You can use the 'config' command to show your current configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			pt := filepath.Dir(path)
			nm := filepath.Base(path)
			cfgManager := config.NewViperConfig(pt, nm)

			if ptErr := cfgManager.SetPath(pt); ptErr != nil {
				return ptErr
			}
			if nmErr := cfgManager.SetName(nm); nmErr != nil {
				return nmErr
			}

			return nil
		},
	}

	cfgCmd.Flags().StringVarP(&path, "path", "p", "", "Path to the configuration file")

	return cfgCmd
}

func cmdConfigGet() *cobra.Command {
	var path string

	cfgCmd := &cobra.Command{
		Use:   "get",
		Short: "Get a configuration",
		Long: `Get a configuration. This command will get a configuration.
You can use the 'config' command to show your current configuration.`,
		Run: func(cmd *cobra.Command, args []string) {
			pt := filepath.Dir(path)
			nm := filepath.Base(path)

			cfgManager := config.NewViperConfig(pt, nm)
			cfgMap := cfgManager.GetConfig()

			fmt.Println(fmt.Sprintf("Configuration file: %s", path))
			fmt.Println("Configuration:")

			fmt.Println(cfgMap)
		},
	}

	cfgCmd.Flags().StringVarP(&path, "path", "p", "", "Path to the configuration file")

	return cfgCmd
}
