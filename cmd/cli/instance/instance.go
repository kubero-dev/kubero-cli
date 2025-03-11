package instance

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"github.com/faelmori/kubero-cli/cmd/common"
	c "github.com/faelmori/kubero-cli/internal/config"
	u "github.com/faelmori/kubero-cli/internal/utils"
	"github.com/spf13/cobra"
)

func InstanceCmds() []*cobra.Command {
	instanceRootCmd := cmdInstance()
	instanceRootCmd.AddCommand(cmdInstanceCreate())
	instanceRootCmd.AddCommand(cmdInstanceDelete())
	instanceRootCmd.AddCommand(cmdInstanceSelect())
	return []*cobra.Command{
		instanceRootCmd,
	}
}

func cmdInstance() *cobra.Command {
	var path string

	var instanceCmd = &cobra.Command{
		Use:     "instance",
		Aliases: []string{"i"},
		Short:   "List available instances",
		Long:    `Print a list of available instances.`,
		Annotations: common.GetDescriptions([]string{
			"List available instances",
			`Print a list of available instances.`,
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			cfg := c.NewViperConfig("", "")
			cfg.GetInstanceManager().PrintInstanceList()
		},
	}

	instanceCmd.Flags().StringVarP(&path, "path", "p", "", "Path to the instance file")

	return instanceCmd
}

func cmdInstanceCreate() *cobra.Command {
	var path string
	var instanceCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new", "cr"},
		Short:   "Create an new instance",
		Long:    `Create an new instance.`,
		Annotations: common.GetDescriptions([]string{
			"Create an new instance",
			`Create an new instance. This command will create a new instance.`,
		}, false),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := c.NewViperConfig("", "")
			return cfg.GetInstanceManager().CreateInstanceForm()
		},
	}

	instanceCreateCmd.Flags().StringVarP(&path, "path", "p", "", "Path to the instance file")

	return instanceCreateCmd
}

func cmdInstanceDelete() *cobra.Command {
	var path string
	var instanceDeleteCmd = &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del", "rm"},
		Short:   "Delete an instance",
		Long:    `Delete an instance.`,
		Annotations: common.GetDescriptions([]string{
			"Delete an instance",
			`Delete an instance.`,
		}, false),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := c.NewViperConfig("", "")
			return cfg.GetInstanceManager().DeleteInstanceForm()
		},
	}

	instanceDeleteCmd.Flags().StringVarP(&path, "path", "p", "", "Path to the instance file")

	return instanceDeleteCmd
}

func cmdInstanceSelect() *cobra.Command {
	var path string
	var instanceSelectCmd = &cobra.Command{
		Use:     "select",
		Aliases: []string{"use"},
		Short:   "Select an instance",
		Long:    `Select an instance to use.`,
		Annotations: common.GetDescriptions([]string{
			"Select an instance",
			`Select an instance to use.`,
		}, false),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := c.NewViperConfig("", "")
			utilsPrompt := u.NewConsolePrompt()
			newInstanceName := utilsPrompt.SelectFromList("Select an instance", cfg.GetInstanceManager().GetInstanceNameList(), "")
			return cfg.GetInstanceManager().SetCurrentInstance(newInstanceName)
		},
	}

	instanceSelectCmd.Flags().StringVarP(&path, "path", "p", "", "Path to the instance file")

	return instanceSelectCmd
}
