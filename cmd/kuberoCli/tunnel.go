/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package kuberoCli

import (
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/jonasfj/go-localtunnel"
	"github.com/leaanthony/spinner"

	"github.com/spf13/cobra"
)

var tunnelHost string
var tunnelPort int
var tunnelSubdomain string
var tunnelDuration string

// tunnelCmd represents the tunnel command
var tunnelCmd = &cobra.Command{
	Use:   "tunnel",
	Short: cfmt.Sprint("Create a tunnel to the cluster in NATed infrastructures {{[BETA]}}::cyan "),
	Long:  `Use the tunnel subcommand to create a tunnel to the cluster in NATed infrastructures.`,
	Run: func(cmd *cobra.Command, args []string) {
		startTunnel()
	},
}

func init() {
	rootCmd.AddCommand(tunnelCmd)
	tunnelCmd.Flags().StringVarP(&tunnelHost, "host", "H", "localhost", "Hostname")
	tunnelCmd.Flags().IntVarP(&tunnelPort, "port", "p", 80, "Port to use")
	tunnelCmd.Flags().StringVarP(&tunnelDuration, "timeout", "t", "1h", "Timeout for the tunnel")

	tunnelCmd.Flags().StringVarP(&tunnelSubdomain, "subdomain", "s", "", "Subdomain to use")
}

func startTunnel() {

	promptWarning("WARNING: your traffic will routed thru localtunnel.me")

	if tunnelSubdomain == "" {
		tunnelSubdomain = promptLine("Subdomain", "", "kubero-"+generateRandomString(13, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"))
	}

	if tunnelSubdomain == "-" {
		tunnelSubdomain = ""
	}

	ipclient := resty.New().R()
	ipres, err := ipclient.Get("https://api.ipify.org")
	if err != nil {
		cfmt.Println("{{✖}}::red Error getting your IP")
		os.Exit(1)
	}
	cfmt.Print("\n  Your 'Endpoint IP' is {{" + ipres.String() + "}}::cyan\n\n")

	spinner := spinner.New()
	spinner.Start("Waiting for tunnel to be ready")

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
		spinner.Error("Error creating tunnel : " + err.Error())
		//cfmt.Println("{{✖}}::red Error creating tunnel : " + err.Error() + "")
		os.Exit(1)
	}
	defer tunnel.Close()

	spinner.UpdateMessage(cfmt.Sprint("Tunnel active at {{" + tunnel.URL() + "}}::cyan with an expiration of {{" + tunnelDuration + "}}::cyan"))
	//cfmt.Println("{{✔}}::green Tunnel created at {{" + tunnel.URL() + "}}::cyan with an expiration of {{" + tunnelDuration + "}}::cyan")

	tunnelTimeout, err := time.ParseDuration(tunnelDuration)
	if err != nil {
		cfmt.Println("{{✗}}::red Error parsing timeout")
		os.Exit(1)
	}

	time.Sleep(tunnelTimeout * time.Second)

}
