package kuberoCli

import (
	"encoding/json"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var podsizesCmd = &cobra.Command{
	Use:   "podsizes",
	Short: "List the available pod sizes",
	Run: func(cmd *cobra.Command, args []string) {
		resp, _ := api.GetPodsize()
		printPodsizes(resp)
	},
}

func init() {
	configCmd.AddCommand(podsizesCmd)
}

func printPodsizes(r *resty.Response) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Description"})
	//table.SetBorder(false)

	var podsizeList []Podsize
	json.Unmarshal(r.Body(), &podsizeList)

	for _, podsize := range podsizeList {
		table.Append([]string{podsize.Name, podsize.Description})
	}

	printCLI(table, r)
}
