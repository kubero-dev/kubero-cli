package kuberoCli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// appsCmd represents the apps command
var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "**DEPRECATED** Manage your apps",
	Long: `**DEPRECATED** Manage your apps

An App runs allways in a Pipeline.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("apps called")

		fmt.Println(getGitRemote())
	},
}

func init() {
	//rootCmd.AddCommand(appsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	appsCmd.Flags().StringP("pipeline", "p", "", "Name of the pipeline")
	appsCmd.MarkFlagRequired("pipeline")
}
