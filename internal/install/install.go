package install

import (
	u "github.com/faelmori/kubero-cli/internal/utils"
	"github.com/kubero-dev/kubero-cli/internal/log"
	"github.com/spf13/cobra"
	"math/rand"
	"time"
)

var (
	utils          = u.NewUtils()
	prompt         = u.NewConsolePrompt()
	promptLine     = prompt.PromptLine
	selectFromList = prompt.SelectFromList
)

type ManagerInstall struct {
	cmd                      *cobra.Command
	installOlm               bool
	argComponent             string
	clusterType              string
	argAdminPassword         string
	argAdminUser             string
	argApiToken              string
	argPort                  string
	argPortSecure            string
	ingressControllerVersion string
	clusterTypeList          []string
}

func NewManagerInstall(installOlm bool, argComponent, clusterType, argAdminPassword, argAdminUser, argApiToken, argPort, argPortSecure string) *ManagerInstall {
	return &ManagerInstall{
		installOlm:               installOlm,
		argComponent:             argComponent,
		clusterType:              clusterType,
		argAdminPassword:         argAdminPassword,
		argAdminUser:             argAdminUser,
		argApiToken:              argApiToken,
		argPort:                  argPort,
		argPortSecure:            argPortSecure,
		ingressControllerVersion: "v1.10.0", // https://github.com/kubernetes/ingress-nginx/tags -> controller-v1.5.1
		clusterTypeList:          []string{"kind", "linode", "scaleway", "gke", "digitalocean"},
	}
}

func (m *ManagerInstall) FullInstallation() error {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	if checkAllBinariesErr := utils.CheckAllBinaries(); checkAllBinariesErr != nil {
		log.Error("Failed to check all binaries")
		return checkAllBinariesErr
	}

	switch m.argComponent {
	case "metrics":
		if insMetricsErr := m.InstallMetrics(); insMetricsErr != nil {
			return insMetricsErr
		}
		return nil
	case "certManager":
		return m.InstallCertManager()
	case "olm":
		return m.InstallKuberoOLMOperator()
	case "kubero-operator":
		return m.InstallKuberoOperator()
	case "kubero-ui":
		return m.InstallKuberoUi()
	case "ingress":
		return m.InstallIngress()
	case "monitoring":
		return m.InstallMonitoring()
	case "kubernetes":
		return m.InstallKubernetes()
	case "":
		//printInstallSteps()
		if insKu8sErr := m.InstallKubernetes(); insKu8sErr != nil { // 1
			return insKu8sErr
		}
		if checkClusterErr := utils.CheckClusters(); checkClusterErr != nil { // 2
			return checkClusterErr
		}
		if installOLMErr := m.InstallKuberoOLMOperator(); installOLMErr != nil { // 3
			return installOLMErr
		}
		if insKubeOperErr := m.InstallKuberoOperator(); insKubeOperErr != nil { // 3
			return insKubeOperErr
		}
		if installIngressErr := m.InstallIngress(); installIngressErr != nil { // 4
			return installIngressErr
		}
		if insMetricsErr := m.InstallMetrics(); insMetricsErr != nil { // 5
			return insMetricsErr
		}
		if installCertManagerErr := m.InstallCertManager(); installCertManagerErr != nil { // 6
			return installCertManagerErr
		}
		if installMonitoringErr := m.InstallMonitoring(); installMonitoringErr != nil { // 7
			return installMonitoringErr
		}
		if installKuberoUiErr := m.InstallKuberoUi(); installKuberoUiErr != nil { // 8
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
