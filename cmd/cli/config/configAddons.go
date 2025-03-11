package config

import (
	"encoding/json"
	a "github.com/faelmori/kubero-cli/internal/api"
	t "github.com/faelmori/kubero-cli/types"
	"github.com/kubero-dev/kubero-cli/cmd/common"
	u "github.com/kubero-dev/kubero-cli/internal/utils"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func ConfigAddonsCmds() []*cobra.Command {
	return []*cobra.Command{
		cmdAddons(),
	}
}

func cmdAddons() *cobra.Command {
	var addonsCmd = &cobra.Command{
		Use:   "addons",
		Short: "A brief description of your command",
		Annotations: common.GetDescriptions([]string{
			"A brief description of your command",
			"Use the addons subcommand to list all available addons.",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			client := a.NewClient()
			resp, _ := client.GetAddons()
			//fmt.Println(resp)
			printAddons(resp)
		},
	}

	return addonsCmd
}

// print the response as a table
func printAddons(r *resty.Response) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Description", "Version", "Beta", "Enabled"})
	table.SetRowLine(true)
	//table.SetBorder(false)

	var addonsList []t.Addon
	unmarshalErr := json.Unmarshal(r.Body(), &addonsList)
	if unmarshalErr != nil {
		return
	}

	for _, addon := range addonsList {
		table.Append([]string{addon.ID, addon.Description, addon.Version.Installed, strconv.FormatBool(addon.Beta), strconv.FormatBool(addon.Enabled)})
	}
	utils := u.NewUtils()
	utils.PrintCLI(table, r, "table")
}
