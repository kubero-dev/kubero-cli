package kuberoCli

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
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

	tunnelCmd.Flags().StringVarP(&tunnelSubdomain, "subdomain", "s", "", "Subdomain to use ('-' to generate a random one)")
}

func startTunnel() {

	promptWarning("WARNING: your traffic will routed thru localtunnel.me")

	if tunnelSubdomain == "" {
		tunnelSubdomainSugestion := "kubero-" + generateRandomString(10, "abcdefghijklmnopqrstuvwxyz0123456789")

		if currentInstance.Tunnel.Subdomain != "" {
			tunnelSubdomainSugestion = currentInstance.Tunnel.Subdomain
		} else if currentInstance.Name != "" {
			tunnelSubdomainSugestion = "kubero-" + currentInstance.Name
		}

		tunnelSubdomain = promptLine("Subdomain", "", tunnelSubdomainSugestion)
	}

	// Check if subdomain is valid
	// localtunnel.me allows only lowercasae letters, numbers and dashes
	if !regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(tunnelSubdomain) {
		_, _ = cfmt.Println("{{✖}}::red Subdomain {{" + tunnelSubdomain + "}}::yellow can only contain lowercase letters, numbers and dashes")
		os.Exit(1)
	}

	// genereate a subdomain if the user entered "-"
	if tunnelSubdomain == "-" {
		tunnelSubdomain = ""
	}

	ipclient := resty.New().R()
	ipres, err := ipclient.Get("https://api.ipify.org")
	if err != nil {
		_, _ = cfmt.Println("{{✖}}::red Error getting your IP")
		os.Exit(1)
	}
	_, _ = cfmt.Println()
	_, _ = cfmt.Println("  Endpoint IP (Tunnel Password) : {{" + ipres.String() + "}}::cyan")
	_, _ = cfmt.Println("  Destination Host              : {{" + tunnelHost + "}}::cyan")
	_, _ = cfmt.Println("  Destination Port              : {{" + strconv.Itoa(tunnelPort) + "}}::cyan\n\n")

	spinnerObj := spinner.New()
	spinnerObj.Start("Waiting for tunnel to be ready")

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
		spinnerObj.Error("Error creating tunnel : " + err.Error())
		//cfmt.Println("{{✖}}::red Error creating tunnel : " + err.Error() + "")
		os.Exit(1)
	}
	defer func(tunnel *localtunnel.LocalTunnel) {
		localTunnelErr := tunnel.Close()
		if localTunnelErr != nil {
			fmt.Println("Error closing tunnel : " + localTunnelErr.Error())
		}
	}(tunnel)

	spinnerObj.UpdateMessage(cfmt.Sprint("Tunnel active at {{" + tunnel.URL() + "}}::cyan with an expiration of {{" + tunnelDuration + "}}::cyan"))
	//cfmt.Println("{{✔}}::green Tunnel created at {{" + tunnel.URL() + "}}::cyan with an expiration of {{" + tunnelDuration + "}}::cyan")

	tunnelTimeout, err := time.ParseDuration(tunnelDuration)
	if err != nil {
		_, _ = cfmt.Println("{{✗}}::red Error parsing timeout")
		os.Exit(1)
	}

	time.Sleep(tunnelTimeout * time.Second)

}
