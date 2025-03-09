package kuberoCli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var remoteCreateCmd = &cobra.Command{
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
	remoteCmd.AddCommand(remoteCreateCmd)
}
