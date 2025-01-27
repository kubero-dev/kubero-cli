package kuberoCli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var podSizesCmd = &cobra.Command{
	Use:   "podsizes",
	Short: "List the available pod sizes",
	Run: func(cmd *cobra.Command, args []string) {
		resp, _ := api.GetPodsize()
		printPodSizes(resp)
	},
}

func init() {
	configCmd.AddCommand(podSizesCmd)
}

func printPodSizes(r *resty.Response) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Description"})
	//table.SetBorder(false)

	var podsizeList []PodSize
	unmarshalErr := json.Unmarshal(r.Body(), &podsizeList)
	if unmarshalErr != nil {
		fmt.Println("Failed to unmarshal the response body:", unmarshalErr)
		return
	}

	for _, podsize := range podsizeList {
		table.Append([]string{podsize.Name, podsize.Description})
	}

	printCLI(table, r)
}
