package kuberoCli

import (
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show your configuration",
	/*
			Long: `A longer description that spans multiple lines and likely contains examples
		and usage of using your command. For example:

		Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.`,
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("config called")
			},
	*/
}

func init() {
	rootCmd.AddCommand(configCmd)
}
