/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show a status of your kubero instance",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("status called")

		url := viper.GetString("api.url")
		fmt.Println(url)

		user := viper.Get("api.user")
		fmt.Println(user)

		res, err := http.Get(url + "/api/config")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res.Body)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
