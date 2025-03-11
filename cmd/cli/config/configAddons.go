package config

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// addonsCmd represents the addons command
var addonsCmd = &cobra.Command{
	Use:   "addons",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		resp, _ := api.GetAddons()
		//fmt.Println(resp)
		printAddons(resp)
	},
}

func init() {
	configCmd.AddCommand(addonsCmd)
}

// print the response as a table
func printAddons(r *resty.Response) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Description", "Version", "Beta", "Enabled"})
	table.SetRowLine(true)
	//table.SetBorder(false)

	var addonsList []Addon
	json.Unmarshal(r.Body(), &addonsList)

	for _, addon := range addonsList {
		table.Append([]string{addon.ID, addon.Description, addon.Version.Installed, strconv.FormatBool(addon.Beta), strconv.FormatBool(addon.Enabled)})
	}

	printCLI(table, r)
}
