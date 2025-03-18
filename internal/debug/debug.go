package debug

import (
	"fmt"
	l "github.com/kubero-dev/kubero-cli/internal/log"
	u "github.com/kubero-dev/kubero-cli/internal/utils"
	v "github.com/kubero-dev/kubero-cli/version"
	"os"
	"os/exec"
	"runtime"
)

var (
	utils = u.NewUtils()
	log   = l.Logger()
)

type Debug struct{}

func NewDebug() *Debug { return &Debug{} }

func (d *Debug) Run() error {
	log.Info("Kubero CLI Debug Information", nil)

	d.PrintCLIVersion()
	d.PrintOsArch()

	log.Info("Kubernetes Debug Information", nil)
	if err := d.PrintKubernetesVersion(); err != nil {
		return err
	}

	log.Info("Kubero Operator Debug Information", nil)
	d.CheckKuberoOperator()

	log.Info("Kubero UI Debug Information", nil)
	d.CheckKuberoUI()

	log.Info("Cert Manager Debug Information", nil)
	d.CheckCertManager()

	return nil
}
func (d *Debug) PrintCLIVersion() {
	log.Info(fmt.Sprintf("Kubero CLI Version: %s", v.Version()), nil)
}
func (d *Debug) PrintOsArch() {
	log.Info("OS/Arch", nil)
	log.Info(fmt.Sprintf("OS: %s", runtime.GOOS), nil)
	log.Info(fmt.Sprintf("Arch: %s", runtime.GOARCH), nil)
	log.Info("Runtime Information", nil)
	log.Info(fmt.Sprintf("Version: %s", runtime.Version()), nil)
	log.Info(fmt.Sprintf("Compiler: %s", runtime.Compiler), nil)
	log.Info("Performance Information", nil)
	log.Info(fmt.Sprintf("NumCPU: %d", runtime.NumCPU()), nil)
	log.Info(fmt.Sprintf("Memory: %d MB", runtime.MemStats{}.TotalAlloc/1024/1024), nil)
	log.Info("Process Information", nil)
	log.Info(fmt.Sprintf("Go NumGoroutine: %d", runtime.NumGoroutine()), nil)
	log.Info(fmt.Sprintf("Go NumCgoCall: %d", runtime.NumCgoCall()), nil)
	log.Info(fmt.Sprintf("Go GOMAXPROCS: %d", runtime.GOMAXPROCS(0)), nil)
}
func (d *Debug) PrintKubernetesVersion() error {
	if !utils.CheckBinary("kubectl") {
		log.Error("kubectl is not installed. Please install kubectl and try again.", map[string]interface{}{
			"context": "kubero-cli",
			"action":  "PrintKubernetesVersion",
			"error":   "kubectl is not installed",
		})
		return os.ErrNotExist
	}

	kVersion, err := exec.Command("kubectl", "version", "--client=true", "-o", "yaml").Output()
	if err != nil {
		log.Error("Failed to fetch Kubernetes version.", map[string]interface{}{
			"context": "kubero-cli",
			"action":  "PrintKubernetesVersion",
			"error":   err.Error(),
		})
		return err
	}

	log.Info(fmt.Sprintf("Kubernetes Version:\n%s", string(kVersion)), nil)

	return nil
}
func (d *Debug) CheckKuberoOperator() {
	output, _ := exec.Command("kubectl", "get", "deployments.apps", "-n", "kubero-operator-system").Output()
	log.Info(string(output), nil)

	log.Info("Kubero Operator Image", nil)
	imageOutput, _ := exec.Command("kubectl", "get", "deployment", "kubero-operator-controller-manager", "-o=jsonpath={$.spec.template.spec.containers[:1].image}", "-n", "kubero-operator-system").Output()
	log.Info(string(imageOutput), nil)
}
func (d *Debug) CheckKuberoUI() {
	output, _ := exec.Command("kubectl", "get", "deployments.apps", "-n", "kubero").Output()
	log.Info(string(output), nil)

	log.Info("Kubero UI Ingress", nil)
	ingressOutput, _ := exec.Command("kubectl", "get", "ingress", "-n", "kubero").Output()
	log.Info(string(ingressOutput), nil)
}
func (d *Debug) CheckCertManager() {
	output, _ := exec.Command("kubectl", "get", "deployments.apps", "-n", "cert-manager").Output()
	log.Info(string(output), nil)

	log.Info("Cert Manager Issuers", nil)
	issuersOutput, _ := exec.Command("kubectl", "get", "clusterissuers.cert-manager.io").Output()
	log.Info(string(issuersOutput), nil)
}
func (d *Debug) CheckMetricsServer() {
	output, _ := exec.Command("kubectl", "get", "deployments.apps", "metrics-server", "-n", "kube-system").Output()
	log.Info(string(output), nil)

	log.Info("Metrics Server Image", nil)
	imageOutput, _ := exec.Command("kubectl", "get", "deployment", "metrics-server", "-o=jsonpath={$.spec.template.spec.containers[:1].image}", "-n", "kube-system").Output()
	log.Info(string(imageOutput), nil)
}
