/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Api struct {
		Url   string `json:"url" yaml:"url"`
		Token string `json:"token" yaml:"token"`
	} `json:"api" yaml:"api"`
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Configure your kubero-cli",
	/*
			Long: `A longer description that spans multiple lines and likely contains examples
		and usage of using your command. For example:

		Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.`,
	*/
	Run: func(cmd *cobra.Command, args []string) {

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

		viper.WriteConfig()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
