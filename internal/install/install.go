package install

import (
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

var (
	installOlm          bool
	monitoringInstalled bool
)
var (
	argComponent     string
	argDomain        string
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
	Run: func(cmd *cobra.Command, args []string) {
		rand.New(rand.NewSource(time.Now().UnixNano()))

		utils.CheckAllBinaries()

		switch argComponent {
		case "metrics":
			installMetrics()
			return
		case "certManager":
			installCertManager()
			return
		case "olm":
			installKuberoOLMOperator()
			//installOLM()
			return
		case "kubero-operator":
			installKuberoOperator()
			return
		case "kubero-ui":
			installKuberoUi()
			return
		case "ingress":
			installIngress()
			return
		case "monitoring":
			installMonitoring()
			return
		case "kubernetes":
			installKubernetes()
			//checkCluster()
			return
		case "":
			//printInstallSteps()
			installKubernetes() // 1
			//checkCluster()          //
			//installOLM()            // 2
			installKuberoOperator() // 3
			installIngress()        // 4
			installMetrics()        // 5
			installCertManager()    // 6
			installMonitoring()     // 7
			installKuberoUi()       // 8
			//writeCLIConfig()        // 9
			//printDNSinfo()
			//finalMessage()
			return
		default:
			return
		}
	},
}

type ManagerInstall struct{}
