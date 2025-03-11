package install

import (
	"github.com/faelmori/kubero-cli/internal/log"
	"github.com/faelmori/kubero-cli/internal/utils"
	"github.com/spf13/cobra"
	"math/rand"
	"time"
)

var (
	prompt         = utils.NewConsolePrompt()
	promptLine     = prompt.PromptLine
	selectFromList = prompt.SelectFromList
)

var installOlm bool
var (
	argComponent     string
	clusterType      string
	argAdminPassword string
	argAdminUser     string
	argApiToken      string
	argPort          string
	argPortSecure    string
)
var (
	ingressControllerVersion = "v1.10.0" // https://github.com/kubernetes/ingress-nginx/tags -> controller-v1.5.1
	clusterTypeList          = []string{"kind", "linode", "scaleway", "gke", "digitalocean"}
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Create a Kubernetes cluster and install all required components for kubero",
	Long: `This command will create a kubernetes cluster and install all required components 
for kubero on any kubernetes cluster.

required binaries:
 - kubectl
 - kind (optional)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return NewManagerInstall().runAll()
	},
}

type ManagerInstall struct {
	cmd *cobra.Command
}

func NewManagerInstall() *ManagerInstall {
	return &ManagerInstall{
		cmd: installCmd,
	}
}

func (m *ManagerInstall) InstallKubernetes() error        { return installKubernetes() }
func (m *ManagerInstall) InstallKuberoOperator() error    { return installKuberoOperator() }
func (m *ManagerInstall) InstallKuberoUi() error          { return installKuberoUi() }
func (m *ManagerInstall) InstallIngress() error           { return installIngress() }
func (m *ManagerInstall) InstallMetrics() error           { return installMetrics() }
func (m *ManagerInstall) InstallCertManager() error       { return installCertManager() }
func (m *ManagerInstall) InstallMonitoring() error        { return installMonitoring() }
func (m *ManagerInstall) InstallKuberoOLMOperator() error { return installKuberoOLMOperator() }

func (m *ManagerInstall) runAll() error {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	if checkAllBinariesErr := utils.CheckAllBinaries(); checkAllBinariesErr != nil {
		log.Error("Failed to check all binaries")
		return checkAllBinariesErr
	}

	switch argComponent {
	case "metrics":
		if insMetricsErr := installMetrics(); insMetricsErr != nil {
			return insMetricsErr
		}
		return nil
	case "certManager":
		return installCertManager()
	case "olm":
		return installKuberoOLMOperator()
	case "kubero-operator":
		return installKuberoOperator()
	case "kubero-ui":
		return installKuberoUi()
	case "ingress":
		return installIngress()
	case "monitoring":
		return installMonitoring()
	case "kubernetes":
		return installKubernetes()
	case "":
		//printInstallSteps()
		if insKu8sErr := installKubernetes(); insKu8sErr != nil { // 1
			return insKu8sErr
		}
		if checkClusterErr := utils.CheckClusters(); checkClusterErr != nil { // 2
			return checkClusterErr
		}
		if installOLMErr := installKuberoOLMOperator(); installOLMErr != nil { // 3
			return installOLMErr
		}
		if insKubeOperErr := installKuberoOperator(); insKubeOperErr != nil { // 3
			return insKubeOperErr
		}
		if installIngressErr := installIngress(); installIngressErr != nil { // 4
			return installIngressErr
		}
		if insMetricsErr := installMetrics(); insMetricsErr != nil { // 5
			return insMetricsErr
		}
		if installCertManagerErr := installCertManager(); installCertManagerErr != nil { // 6
			return installCertManagerErr
		}
		if installMonitoringErr := installMonitoring(); installMonitoringErr != nil { // 7
			return installMonitoringErr
		}
		if installKuberoUiErr := installKuberoUi(); installKuberoUiErr != nil { // 8
			return installKuberoUiErr
		}
		//if installDNSErr := installDNS(); installDNSErr != nil { // 9
		//	return installDNSErr
		//}
		//if writeKuberoConfigErr := writeCliConfig(); writeKuberoConfigErr != nil { // 10
		//	return writeKuberoConfigErr
		//}
		//printDNSinfo()
		//finalMessage()
		return nil
	default:
		return nil
	}
}

func init() {
	installCmd.Flags().StringVarP(&argComponent, "component", "c", "", "Component to install")
	installCmd.Flags().StringVarP(&clusterType, "cluster-type", "t", "", "Type of cluster to install")
	installCmd.Flags().StringVarP(&argAdminPassword, "admin-password", "p", "", "Admin password for kubero")
	installCmd.Flags().StringVarP(&argAdminUser, "admin-user", "u", "", "Admin user for kubero")
	installCmd.Flags().StringVarP(&argApiToken, "api-token", "a", "", "Api token for kubero")
	installCmd.Flags().StringVarP(&argPort, "port", "o", "", "Port for kubero")
	installCmd.Flags().StringVarP(&argPortSecure, "port-secure", "s", "", "Secure port for kubero")
	installCmd.Flags().BoolVarP(&installOlm, "olm", "l", false, "Install OLM for kubero")
}
