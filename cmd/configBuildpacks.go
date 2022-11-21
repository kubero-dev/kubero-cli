package cmd

import (
	"encoding/json"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// buildpacksCmd represents the buildpacks command
var buildpacksCmd = &cobra.Command{
	Use:   "buildpacks",
	Short: "List the available buildpacks",
	Run: func(cmd *cobra.Command, args []string) {
		resp, _ := client.Get("/api/cli/config/buildpacks")
		printBuildpacks(resp)
	},
}

func init() {
	configCmd.AddCommand(buildpacksCmd)
}

var buildPacksSimpleList []string

type buildPacks []struct {
	Name     string `json:"name"`
	Language string `json:"language"`
	Fetch    struct {
		Repository string `json:"repository"`
		Tag        string `json:"tag"`
	} `json:"fetch"`
	Build struct {
		Repository string `json:"repository"`
		Tag        string `json:"tag"`
		Command    string `json:"command"`
	} `json:"build"`
	Run struct {
		Repository         string `json:"repository"`
		Tag                string `json:"tag"`
		ReadOnlyAppStorage bool   `json:"readOnlyAppStorage"`
		SecurityContext    *struct {
			AllowPrivilegeEscalation *bool `json:"allowPrivilegeEscalation"`
			ReadOnlyRootFilesystem   *bool `json:"readOnlyRootFilesystem"`
		} `json:"securityContext"`
		Command string `json:"command"`
	} `json:"run,omitempty"`
}

func loadBuildpacks() {

	b, _ := client.Get("/api/cli/config/buildpacks")

	var buildPacks buildPacks
	json.Unmarshal(b.Body(), &buildPacks)

	for _, buildPack := range buildPacks {
		buildPacksSimpleList = append(buildPacksSimpleList, buildPack.Name)
	}

	//buildPacks = []string{"java", "node", "python", "ruby", "php"}
}

// print the response as a table
func printBuildpacks(r *resty.Response) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Language", "Phase", "Image", "Command"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	//table.SetBorder(false)

	var buildPacksList buildPacks
	json.Unmarshal(r.Body(), &buildPacksList)

	for _, podsize := range buildPacksList {
		table.Append([]string{
			podsize.Name,
			podsize.Language,
			"Fetch",
			podsize.Fetch.Repository + ":" + podsize.Fetch.Tag,
			"git clone",
		})
		table.Append([]string{
			podsize.Name,
			podsize.Language,
			"Build",
			podsize.Build.Repository + ":" + podsize.Build.Tag,
			podsize.Build.Command,
		})
		table.Append([]string{
			podsize.Name,
			podsize.Language,
			"Run",
			podsize.Run.Repository + ":" + podsize.Run.Tag + " (ro)",
			podsize.Run.Command,
		})
	}

	printCLI(table, r)
}
