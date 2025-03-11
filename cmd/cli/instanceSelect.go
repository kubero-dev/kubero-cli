package cli

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"github.com/spf13/cobra"
)

// instanceSelectCmd represents the instanceSelect command
var instanceSelectCmd = &cobra.Command{
	Use:     "select",
	Aliases: []string{"use"},
	Short:   "Select an instance",
	Long:    `Select an instance to use.`,
	Run: func(cmd *cobra.Command, args []string) {
		newInstanceName := selectFromList("Select an instance", instanceNameList, "")
		setCurrentInstance(newInstanceName)
	},
}

func init() {
	instanceCmd.AddCommand(instanceSelectCmd)
}
