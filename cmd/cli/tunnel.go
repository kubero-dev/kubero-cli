package cli

import (
	"github.com/faelmori/kubero-cli/cmd/common"
	"github.com/faelmori/kubero-cli/internal/network"
	"github.com/i582/cfmt/cmd/cfmt"

	"github.com/spf13/cobra"
)

func TunnelCmds() []*cobra.Command {
	return []*cobra.Command{
		cmdTunnel(),
	}
}

func cmdTunnel() *cobra.Command {
	var tunnelHost string
	var tunnelPort int
	var tunnelSubdomain string
	var tunnelDuration string

	cmd := &cobra.Command{
		Use:   "tunnel",
		Short: cfmt.Sprint("Create a tunnel to the cluster in NATed infrastructures {{[BETA]}}::cyan "),
		Long:  `Use the tunnel subcommand to create a tunnel to the cluster in NATed infrastructures.`,
		Annotations: common.GetDescriptions([]string{
			cfmt.Sprint("Create a tunnel to the cluster in NATed infrastructures {{[BETA]}}::cyan "),
			`Use the tunnel subcommand to create a tunnel to the cluster in NATed infrastructures.`,
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			tunnel := network.NewTunnel(tunnelPort, tunnelHost, tunnelSubdomain, tunnelDuration)
			tunnel.StartTunnel()
		},
	}

	cmd.Flags().StringVarP(&tunnelHost, "host", "H", "localhost", "Hostname")
	cmd.Flags().IntVarP(&tunnelPort, "port", "p", 80, "Port to use")
	cmd.Flags().StringVarP(&tunnelDuration, "timeout", "t", "1h", "Timeout for the tunnel")

	cmd.Flags().StringVarP(&tunnelSubdomain, "subdomain", "s", "", "Subdomain to use ('-' to generate a random one)")

	return cmd
}
