package cli

import (
	"github.com/kubero-dev/kubero-cli/cmd/common"
	c "github.com/kubero-dev/kubero-cli/internal/config"
	"github.com/spf13/cobra"
)

func LoginCmds() []*cobra.Command {
	return []*cobra.Command{
		cmdLogin(),
	}
}

func cmdLogin() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Login to your Kubero instance",
		Long:  `Use the login subcommand to login to your Kubero instance.`,
		Annotations: common.GetDescriptions([]string{
			"Login to your Kubero instance",
			`Use the login subcommand to login to your Kubero instance.`,
		}, false),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := c.NewViperConfig("", "")
			if ensureOrCreateErr := cfg.GetInstanceManager().EnsureInstanceOrCreate(); ensureOrCreateErr != nil {
				return ensureOrCreateErr
			}
			cfg.GetCredentialsManager().SetCredentials(cfg.GetCredentialsManager().GetCredentials())
			return nil
		},
	}
}
