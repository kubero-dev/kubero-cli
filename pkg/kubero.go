/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"kubero/cmd"
	"os"

	"github.com/spf13/viper"
)

func main() {
	viper.SetDefault("ContentDir", "content")
	viper.SetDefault("LayoutDir", "layouts")
	viper.SetConfigName("kubero") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	credentials := viper.New()
	credentials.SetConfigName("credentials")   // name of config file (without extension)
	credentials.SetConfigType("yaml")          // REQUIRED if the config file does not have the extension in the name
	credentials.AddConfigPath("/etc/kubero/")  // path to look for the config file in
	credentials.AddConfigPath("$HOME/.kubero") // call multiple times to add many search paths
	viper.AddConfigPath(".")
	credentials.ReadInConfig()

	viper.MergeConfigMap(credentials.AllSettings())

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found; using defaults")
		} else {
			fmt.Println("No config file found")
			os.Exit(1)
			//panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}
	cmd.Execute()
}
