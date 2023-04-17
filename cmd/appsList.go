package cmd

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var appsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List apps in a pipeline",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		if pipeline == "" {
			//iterate over all pipelines in pipelineConfigList
			for _, pc := range pipelineConfigList {

				pipeline = pc.GetString("spec.name")

				pipelineResp := getAppsList(pipeline)

				cfmt.Println("{{Pipeline: " + pipeline + "}}::bold|underline")
				printAppsList(pipelineResp)
			}
		} else {
			pipelineResp := getAppsList(pipeline)

			cfmt.Println("{{Pipeline: " + pipeline + "}}::bold|underline")
			printAppsList(pipelineResp)
		}
	},
}

func init() {
	appsListCmd.Flags().StringVarP(&pipeline, "pipeline", "p", "", "Name of the Pipeline")
	//appsListCmd.MarkFlagRequired("pipeline")
	appsCmd.AddCommand(appsListCmd)
}

func getAppsList(pipeline string) *resty.Response {
	checkConfig()
	pipelineResp, err := client.Get("/api/cli/pipelines/" + pipeline + "/apps")
	if err != nil {
		cfmt.Println("{{  Error getting pipeline list}}::red")
		os.Exit(1)
	}
	return pipelineResp
}

func printAppsList(r *resty.Response) {
	//fmt.Println(r)

	var pl Pipeline
	json.Unmarshal(r.Body(), &pl)

	for _, phase := range pl.Phases {
		if !phase.Enabled {
			continue
		}

		cfmt.Println("{{  " + strings.ToUpper(phase.Name) + "}}::bold")

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{
			"Name",
			"Phase",
			"Pipeline",
			"Repository",
			"Domain",
		})
		table.SetBorder(false)

		for _, app := range phase.Apps {
			table.Append([]string{
				app.Name,
				app.Phase,
				app.Pipeline,
				app.Gitrepo.CloneURL + ":" +
					app.Gitrepo.DefaultBranch,
				app.Domain,
			})
		}

		printCLI(table, r)
		print("\n")
	}
}
