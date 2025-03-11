package config

import (
	"encoding/json"
	"fmt"
	"github.com/kubero-dev/kubero-cli/cmd/common"
	a "github.com/kubero-dev/kubero-cli/internal/api"
	u "github.com/kubero-dev/kubero-cli/internal/utils"
	t "github.com/kubero-dev/kubero-cli/types"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func ConfigPodsizesCmds() []*cobra.Command {
	return []*cobra.Command{
		cmdPodSizes(),
	}
}

func cmdPodSizes() *cobra.Command {
	var podSizesCmd = &cobra.Command{
		Use:   "podsizes",
		Short: "List the available pod sizes",
		Long: `List the available pod sizes. This command will list all available pod sizes.
You can use the 'config' command to show your current configuration.`,
		Annotations: common.GetDescriptions([]string{
			"List the available pod sizes",
			`List the available pod sizes. This command will list all available pod sizes.
You can use the 'config' command to show your current configuration.`,
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			client := a.NewClient()
			resp, _ := client.GetPodsize()
			printPodSizes(resp)
		},
	}

	return podSizesCmd
}

func printPodSizes(r *resty.Response) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Description"})
	//table.SetBorder(false)

	var podsizeList []t.PodSize
	unmarshalErr := json.Unmarshal(r.Body(), &podsizeList)
	if unmarshalErr != nil {
		fmt.Println("Failed to unmarshal the response body:", unmarshalErr)
		return
	}

	for _, podsize := range podsizeList {
		table.Append([]string{podsize.Name, podsize.Description})
	}

	utils := u.NewUtils()
	utils.PrintCLI(table, r, "table")
}
