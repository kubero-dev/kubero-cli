package config

import (
	"encoding/json"
	a "github.com/faelmori/kubero-cli/internal/api"
	t "github.com/faelmori/kubero-cli/types"
	"github.com/kubero-dev/kubero-cli/cmd/common"
	u "github.com/kubero-dev/kubero-cli/internal/utils"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func ConfigBuildpacksCmds() []*cobra.Command {
	return []*cobra.Command{
		cmdBuildpacks(),
	}
}

func cmdBuildpacks() *cobra.Command {
	var buildpacksCmd = &cobra.Command{
		Use:   "buildpacks",
		Short: "List the available buildpacks",
		Long: `List the available buildpacks. This command will list all available buildpacks.
You can use the 'config' command to show your current configuration.`,
		Annotations: common.GetDescriptions([]string{
			"List the available buildpacks",
			`List the available buildpacks. This command will list all available buildpacks.
You can use the 'config' command to show your current configuration.`,
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			client := a.NewClient()
			resp, _ := client.GetBuildpacks()
			printBuildpacks(resp)
		},
	}

	return buildpacksCmd
}

var buildPacksSimpleList []string

func loadBuildpacks() {
	client := a.NewClient()
	b, _ := client.GetBuildpacks()

	var buildPacks t.BuildPacks
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

	var buildPacksList t.BuildPacks
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

	utils := u.NewUtils()
	utils.PrintCLI(table, r, "table")
}
