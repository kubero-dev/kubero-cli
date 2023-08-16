package kuberoCli

import (
	"encoding/json"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// podsizesCmd represents the podsizes command
var podsizesCmd = &cobra.Command{
	Use:   "podsizes",
	Short: "List the available pod sizes",
	Run: func(cmd *cobra.Command, args []string) {
		resp, _ := client.Get("/api/cli/config/podsize")
		printPodsizes(resp)
	},
}

func init() {
	configCmd.AddCommand(podsizesCmd)
}

type PodsizeList []struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Default     bool   `json:"default,omitempty"`
	Resources   struct {
		Requests struct {
			Memory string `json:"memory"`
			CPU    string `json:"cpu"`
		} `json:"requests"`
		Limits struct {
			Memory string `json:"memory"`
			CPU    string `json:"cpu"`
		} `json:"limits,omitempty"`
	} `json:"resources,omitempty"`
	Active bool `json:"active,omitempty"`
}

// print the response as a table
func printPodsizes(r *resty.Response) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Description"})
	//table.SetBorder(false)

	var podsizeList PodsizeList
	json.Unmarshal(r.Body(), &podsizeList)

	for _, podsize := range podsizeList {
		table.Append([]string{podsize.Name, podsize.Description})
	}

	printCLI(table, r)
}
