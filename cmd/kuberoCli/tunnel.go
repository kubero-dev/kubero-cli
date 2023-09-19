/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package kuberoCli

import (
	"time"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/jonasfj/go-localtunnel"

	"github.com/spf13/cobra"
)

var tunnelHost string
var tunnelPort int
var tunnelSubdomain string

// tunnelCmd represents the tunnel command
var tunnelCmd = &cobra.Command{
	Use:   "tunnel",
	Short: "**EXPERIMENTAL** Create a tunnel to the cluster",
	Long:  `Use the tunnel subcommand to create a tunnel to the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		startTunnel()
	},
}

func init() {
	rootCmd.AddCommand(tunnelCmd)
	tunnelCmd.Flags().StringVarP(&tunnelHost, "host", "H", "localhost", "Hostname")
	tunnelCmd.Flags().IntVarP(&tunnelPort, "port", "p", 2000, "Port to use")

	tunnelCmd.Flags().StringVarP(&tunnelSubdomain, "subdomain", "s", "", "Subdomain to use")
}

func startTunnel() {
	if tunnelSubdomain == "" {
		tunnelSubdomain = promptLine("Subdomain", "", "kubero")
	}

	tunnel, err := localtunnel.New(
		tunnelPort,
		tunnelHost,
		localtunnel.Options{
			Subdomain: tunnelSubdomain,
			//BaseURL:        "https://localtunnel.me",
			MaxConnections: 5,
		},
	)
	if err != nil {
		panic(err)
	}

	cfmt.Println("{{TUNNEL:}}::green Tunnel created at " + tunnel.URL())
	time.Sleep(3600 * time.Second)
	defer tunnel.Close()

}
