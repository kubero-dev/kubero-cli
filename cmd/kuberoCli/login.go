package kuberoCli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

type Access struct {
	AccessToken string `json:"access_token"`
	//ExpiresIn   int    `json:"expires_in"`
	//TokenType   string `json:"token_type"`
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to your Kubero instance",
	Long:  `Use the login subcommand to login to your Kubero instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		ensureInstanceOrCreate()
		login("", "")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func ensureInstanceOrCreate() {

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
		token = promptLine("Token", "", "")
	}

	credentialsConfig.Set(currentInstanceName, token)
	writeConfigErr := credentialsConfig.WriteConfig()
	if writeConfigErr != nil {
		fmt.Println("Error writing config file: ", writeConfigErr)
		return
	}
}

func login(user string, pass string) {

	if user == "" {
		user = promptLine("Username", "", "")
	}

	if pass == "" {
		pass = promptLine("Password", "", "")
	}

	res, err := api.Login(user, pass)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	if res.StatusCode() >= 200 && res.StatusCode() < 300 {

		var a Access
		json.Unmarshal(res.Body(), &a)

		fmt.Println("Login successful")
		fmt.Println("Access token: ", a.AccessToken)

		setKuberoCredentials(a.AccessToken)
	} else {
		fmt.Println(res.StatusCode(), "Login failed")
	}

}
