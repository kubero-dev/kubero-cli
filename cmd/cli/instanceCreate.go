package cli

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"fmt"

	"github.com/spf13/cobra"
)

// instanceCreateCmd represents the instanceCreate command
var instanceCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "cr"},
	Short:   "Create an new instance",
	Long:    `Create an new instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("instanceCreate called")
		createInstanceForm()
	},
}

func init() {
	instanceCmd.AddCommand(instanceCreateCmd)
}
