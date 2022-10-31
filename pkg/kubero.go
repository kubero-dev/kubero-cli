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
	viper.SetDefault("ContentDir", "content")
	viper.SetDefault("LayoutDir", "layouts")
	viper.SetConfigName("kubero")        // name of config file (without extension)
	viper.SetConfigType("yaml")          // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/kubero/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.kubero") // call multiple times to add many search paths
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found; using defaults")
		} else {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}
	cmd.Execute()
}
