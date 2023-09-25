package kuberoCli

import (
	"github.com/spf13/cobra"
)

type Config struct {
	Api struct {
		Url   string `json:"url" yaml:"url"`
		Token string `json:"token" yaml:"token"`
	} `json:"api" yaml:"api"`
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to your Kubero instance",
	/*
			Long: `A longer description that spans multiple lines and likely contains examples
		and usage of using your command. For example:

		Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.`,
	*/
	Run: func(cmd *cobra.Command, args []string) {
		/*
			fmt.Println("Initializing kubero-cli")
			url := promptLine("Kubero Host adress", viper.GetString("api.url"), viper.GetString("api.url"))
			viper.Set("api.url", url)

			token := promptLine("Kubero Token", viper.GetString("api.token"), viper.GetString("api.token"))
			viper.Set("api.token", token)

			var config Config
			if err := viper.Unmarshal(&config); err != nil {
				fmt.Println(err)
				return
			}

			repoAuth := promptLine("Create authentication file in this repository", "[y,n]", "n")
			if repoAuth == "y" {
				viper.WriteConfigAs(".kubero/kubero.yaml") //TODO: make .kubero path configurable
				if err := viper.WriteConfig(); err != nil {
					fmt.Println(err)
					return
				}
			}
		*/
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
