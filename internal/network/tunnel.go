package network

import (
	"fmt"
	c "github.com/faelmori/kubero-cli/internal/config"
	"github.com/faelmori/kubero-cli/internal/log"
	u "github.com/faelmori/kubero-cli/internal/utils"
	"github.com/go-resty/resty/v2"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/jonasfj/go-localtunnel"
	"github.com/leaanthony/spinner"
	"os"
	"regexp"
	"strconv"
	"time"
)

var (
	utilsPrompt          = u.NewConsolePrompt()
	utils                = u.NewUtils()
	promptLine           = utilsPrompt.PromptLine
	selectFromList       = utilsPrompt.SelectFromList
	generateRandomString = utils.GenerateRandomString
)

type Tunnel struct {
	tunnelPort      int
	tunnelHost      string
	tunnelSubdomain string
	tunnelDuration  string
}

func NewTunnel(tunnelPort int, tunnelHost string, tunnelSubdomain string, tunnelDuration string) *Tunnel {
	return &Tunnel{
		tunnelPort:      tunnelPort,
		tunnelHost:      tunnelHost,
		tunnelSubdomain: tunnelSubdomain,
		tunnelDuration:  tunnelDuration,
	}
}

func (t *Tunnel) StartTunnel() {
	log.Warn("WARNING: your traffic will routed thru localtunnel.me")
	cfg := c.NewViperConfig("", "")

	if t.tunnelSubdomain == "" {
		tunnelSubdomainSugestion := "kubero-" + generateRandomString(10, "abcdefghijklmnopqrstuvwxyz0123456789")

		if cfg.GetInstanceManager().GetCurrentInstance().Tunnel.Subdomain != "" {
			tunnelSubdomainSugestion = cfg.GetInstanceManager().GetCurrentInstance().Tunnel.Subdomain
		} else if cfg.GetInstanceManager().GetCurrentInstance().Name != "" {
			tunnelSubdomainSugestion = "kubero-" + cfg.GetInstanceManager().GetCurrentInstance().Name
		}

		t.tunnelSubdomain = promptLine("Subdomain", "", tunnelSubdomainSugestion)
	}

	// Check if subdomain is valid
	// localtunnel.me allows only lowercasae letters, numbers and dashes
	if !regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(t.tunnelSubdomain) {
		_, _ = cfmt.Println("{{✖}}::red Subdomain {{" + t.tunnelSubdomain + "}}::yellow can only contain lowercase letters, numbers and dashes")
		os.Exit(1)
	}

	// genereate a subdomain if the user entered "-"
	if t.tunnelSubdomain == "-" {
		t.tunnelSubdomain = ""
	}

	ipclient := resty.New().R()
	ipres, err := ipclient.Get("https://api.ipify.org")
	if err != nil {
		_, _ = cfmt.Println("{{✖}}::red Error getting your IP")
		os.Exit(1)
	}
	_, _ = cfmt.Println()
	_, _ = cfmt.Println("  Endpoint IP (Tunnel Password) : {{" + ipres.String() + "}}::cyan")
	_, _ = cfmt.Println("  Destination Host              : {{" + t.tunnelHost + "}}::cyan")
	_, _ = cfmt.Println("  Destination Port              : {{" + strconv.Itoa(t.tunnelPort) + "}}::cyan\n\n")

	spinnerObj := spinner.New()
	spinnerObj.Start("Waiting for tunnel to be ready")

	tunnel, err := localtunnel.New(
		t.tunnelPort,
		t.tunnelHost,
		localtunnel.Options{
			Subdomain: t.tunnelSubdomain,
			//BaseURL:        "https://localtunnel.me",
			MaxConnections: 5,
		},
	)
	if err != nil {
		spinnerObj.Error("Error creating tunnel : " + err.Error())

		os.Exit(1)
	}
	defer func(tunnel *localtunnel.LocalTunnel) {
		localTunnelErr := tunnel.Close()
		if localTunnelErr != nil {
			fmt.Println("Error closing tunnel : " + localTunnelErr.Error())
		}
	}(tunnel)

	spinnerObj.UpdateMessage(cfmt.Sprint("Tunnel active at {{" + tunnel.URL() + "}}::cyan with an expiration of {{" + t.tunnelDuration + "}}::cyan"))
	//cfmt.Println("{{✔}}::green Tunnel created at {{" + tunnel.URL() + "}}::cyan with an expiration of {{" + tunnelDuration + "}}::cyan")

	tunnelTimeout, err := time.ParseDuration(t.tunnelDuration)
	if err != nil {
		_, _ = cfmt.Println("{{✗}}::red Error parsing timeout")
		os.Exit(1)
	}

	time.Sleep(tunnelTimeout * time.Second)

}
