package kuberoCli

import (
	"github.com/spf13/cobra"
)

var remoteSelectCmd = &cobra.Command{
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
	remoteCmd.AddCommand(remoteSelectCmd)
}
