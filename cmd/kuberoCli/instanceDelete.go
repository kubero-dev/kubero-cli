/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package kuberoCli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// instanceDeleteCmd represents the instanceDelete command
var instanceDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"del", "rm"},
	Short:   "Delete an instance",
	Long:    `Delete an instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("instanceDelete called")
		deleteInstanceForm()
	},
}

func init() {
	instanceCmd.AddCommand(instanceDeleteCmd)
}
