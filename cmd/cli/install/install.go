package install

import (
	"github.com/faelmori/kubero-cli/internal/install"
	"github.com/spf13/cobra"
)

func InstallCmds() []*cobra.Command {
	return []*cobra.Command{
		cmdInstall(),
		cmdInstallMetrics(),
		cmdInstallCertManager(),
	}
}

func cmdInstall() *cobra.Command {
	var (
		installOlm       bool
		argComponent     string
		clusterType      string
		argAdminPassword string
		argAdminUser     string
		argApiToken      string
		argPort          string
		argPortSecure    string
	)

	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Create a Kubernetes cluster and install all required components for kubero",
		Long: `This command will create a kubernetes cluster and install all required components 
for kubero on any kubernetes cluster.

required binaries:
 - kubectl
 - kind (optional)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			insMgr := install.NewManagerInstall(installOlm, argComponent, clusterType, argAdminPassword, argAdminUser, argApiToken, argPort, argPortSecure)
			return insMgr.FullInstallation()
		},
	}

	installCmd.Flags().StringVarP(&argComponent, "component", "c", "", "Component to install")
	installCmd.Flags().StringVarP(&clusterType, "cluster-type", "t", "", "Type of cluster to install")
	installCmd.Flags().StringVarP(&argAdminPassword, "admin-password", "p", "", "Admin password for kubero")
	installCmd.Flags().StringVarP(&argAdminUser, "admin-user", "u", "", "Admin user for kubero")
	installCmd.Flags().StringVarP(&argApiToken, "api-token", "a", "", "Api token for kubero")
	installCmd.Flags().StringVarP(&argPort, "port", "o", "", "Port for kubero")
	installCmd.Flags().StringVarP(&argPortSecure, "port-secure", "s", "", "Secure port for kubero")
	installCmd.Flags().BoolVarP(&installOlm, "olm", "l", false, "Install OLM for kubero")

	return installCmd
}

func cmdInstallMetrics() *cobra.Command {
	var installMetricsCmd = &cobra.Command{
		Use:   "metrics",
		Short: "Install metrics for kubero",
		Long:  `Install metrics for kubero`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// This logic didn't need any parameters, so I removed them and kept just the function call with component name self-contained
			insMgr := install.NewManagerInstall(false, "metrics", "", "", "", "", "", "")
			return insMgr.InstallMetrics()
		},
	}

	return installMetricsCmd
}

func cmdInstallCertManager() *cobra.Command {
	var installOlm bool

	var installCertManagerCmd = &cobra.Command{
		Use:   "cert-manager",
		Short: "Install cert-manager for kubero",
		Long:  `Install cert-manager for kubero`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// This logic contains some installOlm use case, so I kept it
			insMgr := install.NewManagerInstall(installOlm, "certManager", "", "", "", "", "", "")
			return insMgr.InstallCertManager()
		},
	}

	installCertManagerCmd.Flags().BoolVarP(&installOlm, "olm", "l", false, "Install OLM for kubero")

	return installCertManagerCmd
}
