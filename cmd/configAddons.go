/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

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
		resp, _ := client.Get("/api/cli/addons")
		//fmt.Println(resp)
		printAddons(resp)
	},
}

func init() {
	configCmd.AddCommand(addonsCmd)
}

type AddonsList []struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
	Version struct {
		Latest    string `json:"latest"`
		Installed string `json:"installed"`
	} `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
	Readme      string `json:"readme,omitempty"`
	ArtifactURL string `json:"artifact_url"`
	Kind        string `json:"kind"`
	Install     string `json:"install"`
	Beta        bool   `json:"beta"`
}

// print the response as a table
func printAddons(r *resty.Response) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Description", "Version", "Beta", "Enabled"})
	table.SetRowLine(true)
	//table.SetBorder(false)

	var addonsList AddonsList
	json.Unmarshal(r.Body(), &addonsList)

	for _, addon := range addonsList {
		table.Append([]string{addon.ID, addon.Description, addon.Version.Installed, strconv.FormatBool(addon.Beta), strconv.FormatBool(addon.Enabled)})
	}

	printCLI(table, r)
}
