package kuberoCli

import (
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to your Kubero instance",
	Long:  `Use the login subcommand to login to your Kubero instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		ensureIntanceOrCreate()
		setKuberoCredentials("")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func ensureIntanceOrCreate() {

	instanceNameList = append(instanceNameList, "<create new>")

	instanceName := selectFromList("Select an instance", instanceNameList, currentInstanceName)
	if instanceName == "<create new>" {
		createInstanceForm()
	} else {
		setCurrentInstance(instanceName)
	}

}

func setKuberoCredentials(token string) {

	if token == "" {
		token = promptLine("Kubero Token", "", "")
	}

	credentialsConfig.Set(currentInstanceName, token)
	credentialsConfig.WriteConfig()
}
