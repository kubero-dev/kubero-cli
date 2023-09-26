package kuberoCli

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
		resp, _ := api.GetBuildpacks()
		printBuildpacks(resp)
	},
}

func init() {
	configCmd.AddCommand(buildpacksCmd)
}

var buildPacksSimpleList []string

func loadBuildpacks() {

	b, _ := api.GetBuildpacks()

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
