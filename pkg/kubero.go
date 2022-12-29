/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"kubero/cmd"

	"github.com/spf13/viper"
)

func main() {

	loadConfig()
	cmd.InitClient()

	cmd.Execute()
}

// var configPath = os.Getenv("HOME") + "/.config/kubero-cli"

func loadConfig() {

	viper.SetDefault("api.url", "http://localhost:2000")
	viper.SetConfigName("kubero") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // TODO this should search for the git repo root
	err := viper.ReadInConfig()

	personal := viper.New()
	personal.SetConfigName("kubero")        // name of config file (without extension)
	personal.SetConfigType("yaml")          // REQUIRED if the config file does not have the extension in the name
	personal.AddConfigPath("/etc/kubero/")  // path to look for the config file in
	personal.AddConfigPath("$HOME/.kubero") // call multiple times to add many search paths
	errCred := personal.ReadInConfig()

	viper.MergeConfigMap(personal.AllSettings())
	if err != nil && errCred != nil {

		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found; using defaults")
		} else {
			fmt.Printf("Error while loading config files: %v \n\n\n%v", err, errCred)
		}
	}
}
