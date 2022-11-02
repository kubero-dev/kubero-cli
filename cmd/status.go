/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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

		resp, err := client.Get("/api/cli/config/podsize")

		fmt.Println(client.Header)

		// Explore response object
		fmt.Println("Response Info:")
		fmt.Println("  Error      :", err)
		fmt.Println("  Status Code:", resp.StatusCode())
		fmt.Println("  Status     :", resp.Status())
		fmt.Println("  Proto      :", resp.Proto())
		fmt.Println("  Time       :", resp.Time())
		fmt.Println("  Received At:", resp.ReceivedAt())
		fmt.Println("  Body       :\n", resp)
		fmt.Println()

		// Explore trace info
		fmt.Println("Request Trace Info:")
		ti := resp.Request.TraceInfo()
		fmt.Println("  DNSLookup    :", ti.DNSLookup)
		fmt.Println("  ConnTime     :", ti.ConnTime)
		fmt.Println("  TCPConnTime  :", ti.TCPConnTime)
		fmt.Println("  TLSHandshake :", ti.TLSHandshake)
		fmt.Println("  ServerTime   :", ti.ServerTime)
		fmt.Println("  ResponseTime :", ti.ResponseTime)
		fmt.Println("  TotalTime    :", ti.TotalTime)
		fmt.Println("  IsConnReused :", ti.IsConnReused)
		fmt.Println("  IsConnWasIdle:", ti.IsConnWasIdle)
		fmt.Println("  ConnIdleTime :", ti.ConnIdleTime)
		fmt.Println()
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
