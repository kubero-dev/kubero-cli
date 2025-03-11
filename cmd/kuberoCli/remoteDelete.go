package kuberoCli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var remoteDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"del", "rm"},
	Short:   "Delete an instance from the local configuration",
	Long:    `Delete an instance from the local configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("instanceDelete called")
		deleteInstanceForm()
	},
}

func init() {
	remoteCmd.AddCommand(remoteDeleteCmd)
}
