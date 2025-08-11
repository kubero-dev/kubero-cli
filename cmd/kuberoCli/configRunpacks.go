package kuberoCli

import (
	"encoding/json"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// runpacksCmd represents the runpacks command
var runpacksCmd = &cobra.Command{
	Use:   "runpacks",
	Short: "List the available runpacks",
	Run: func(cmd *cobra.Command, args []string) {
		resp, _ := api.GetRunpacks()
		printRunpacks(resp)
	},
}

func init() {
	configCmd.AddCommand(runpacksCmd)
}

var runPacksSimpleList []string

func loadRunpacks() {

	b, _ := api.GetRunpacks()

	var runPacks buildPacks
	json.Unmarshal(b.Body(), &runPacks)

	for _, runPack := range runPacks {
		runPacksSimpleList = append(runPacksSimpleList, runPack.Name)
	}

	//runPacks = []string{"java", "node", "python", "ruby", "php"}
}

// print the response as a table
func printRunpacks(r *resty.Response) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Language", "Phase", "Image", "Command"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	//table.SetBorder(false)

	var runPacksList buildPacks
	json.Unmarshal(r.Body(), &runPacksList)

	for _, podsize := range runPacksList {
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
