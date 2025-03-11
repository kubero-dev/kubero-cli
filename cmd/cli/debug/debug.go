package debug

import (
	"github.com/kubero-dev/kubero-cli/cmd/common"
	"github.com/kubero-dev/kubero-cli/internal/debug"
	"github.com/spf13/cobra"
)

func DebugCmds() []*cobra.Command {
	return []*cobra.Command{
		cmdDebug(),
		cmdPrintCLIVersion(),
		cmdPrintOsArch(),
		cmdPrintKubernetesVersion(),
		cmdCheckKuberoOperator(),
		cmdCheckKuberoUI(),
		cmdCheckCertManager(),
		cmdCheckMetricsServer(),
	}
}

func cmdDebug() *cobra.Command {
	debugCmd := &cobra.Command{
		Use:     "debug",
		Aliases: []string{"dbg"},
		Short:   "Print debug information",
		Long: `This command will print debug information like:
    - Kubero CLI version
    - OS/Arch
    - Kubernetes version
    - Kubero operator version
    - Kubero operator namespace
    - Kubernetes metrics server version
    - Kubernetes cert-manager version`,
		Annotations: common.GetDescriptions([]string{
			"Print debug information",
			`This command will print debug information like:
    - Kubero CLI version
    - OS/Arch
    - Kubernetes version
    - Kubero operator version
    - Kubero operator namespace
    - Kubernetes metrics server version
    - Kubernetes cert-manager version`,
		}, false),
		RunE: func(cmd *cobra.Command, args []string) error {
			d := debug.NewDebug()
			return d.Run()
		},
	}

	return debugCmd
}
func cmdPrintCLIVersion() *cobra.Command {
	printCLIVersionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print Kubero CLI version",
		Run: func(cmd *cobra.Command, args []string) {
			d := debug.NewDebug()
			d.PrintCLIVersion()
		},
	}

	return printCLIVersionCmd
}
func cmdPrintOsArch() *cobra.Command {
	printOsArchCmd := &cobra.Command{
		Use:   "os-arch",
		Short: "Print OS/Arch information",
		Run: func(cmd *cobra.Command, args []string) {
			d := debug.NewDebug()
			d.PrintOsArch()
		},
	}

	return printOsArchCmd
}
func cmdPrintKubernetesVersion() *cobra.Command {
	printKubernetesVersionCmd := &cobra.Command{
		Use:   "kubernetes-version",
		Short: "Print Kubernetes version",
		RunE: func(cmd *cobra.Command, args []string) error {
			d := debug.NewDebug()
			return d.PrintKubernetesVersion()
		},
	}

	return printKubernetesVersionCmd
}
func cmdCheckKuberoOperator() *cobra.Command {
	checkKuberoOperatorCmd := &cobra.Command{
		Use:   "kubero-operator",
		Short: "Check Kubero Operator",
		Run: func(cmd *cobra.Command, args []string) {
			d := debug.NewDebug()
			d.CheckKuberoOperator()
		},
	}

	return checkKuberoOperatorCmd
}
func cmdCheckKuberoUI() *cobra.Command {
	checkKuberoUICmd := &cobra.Command{
		Use:   "kubero-ui",
		Short: "Check Kubero UI",
		Run: func(cmd *cobra.Command, args []string) {
			d := debug.NewDebug()
			d.CheckKuberoUI()
		},
	}

	return checkKuberoUICmd
}
func cmdCheckCertManager() *cobra.Command {
	checkCertManagerCmd := &cobra.Command{
		Use:   "cert-manager",
		Short: "Check Cert Manager",
		Run: func(cmd *cobra.Command, args []string) {
			d := debug.NewDebug()
			d.CheckCertManager()
		},
	}

	return checkCertManagerCmd
}
func cmdCheckMetricsServer() *cobra.Command {
	checkMetricsServerCmd := &cobra.Command{
		Use:   "metrics-server",
		Short: "Check Metrics Server",
		Run: func(cmd *cobra.Command, args []string) {
			d := debug.NewDebug()
			d.CheckMetricsServer()
		},
	}

	return checkMetricsServerCmd
}
