/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"kubero/cmd"
	"os"

	"github.com/spf13/viper"
)

func main() {

	loadConfig()
	cmd.InitClient()

	cmd.Execute()
}

// var configPath = os.Getenv("HOME") + "/.config/kubero-cli"

func loadConfig() {

	gitdir := cmd.GetGitdir() + "/../.kubero"

	apiToken := os.Getenv("KUBERO_API_TOKEN")
	apiURL := os.Getenv("KUBERO_API_URL")
	viper.SetDefault("api.url", apiURL)
	viper.SetDefault("api.token", apiToken)
	viper.SetConfigName("credentials") // name of config file (without extension)
	viper.SetConfigType("yaml")        // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(gitdir)
	viper.ReadInConfig()

	// Load personal config
	personal := viper.New()
	personal.SetConfigName("kubero")        // name of config file (without extension)
	personal.SetConfigType("yaml")          // REQUIRED if the config file does not have the extension in the name
	personal.AddConfigPath("/etc/kubero/")  // path to look for the config file in
	personal.AddConfigPath("$HOME/.kubero") // call multiple times to add many search paths
	personal.ReadInConfig()

	/*
		if err != nil && errCred != nil {

			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				fmt.Println("No config file found. Run 'kubero login' to create one.")
				os.Exit(1)
			} else {
				fmt.Printf("Error while loading config files: %v \n\n\n%v", err, errCred)
			}
		}
	*/

	// Merge configs
	viper.MergeConfigMap(personal.AllSettings())
}
