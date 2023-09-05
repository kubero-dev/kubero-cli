package kuberoCli

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/olekukonko/tablewriter"
)

func appsList() {

	pipelineResp, _ := client.Get("/api/cli/pipelines/" + pipelineName + "/apps")

	var pl Pipeline
	json.Unmarshal(pipelineResp.Body(), &pl)

	for _, phase := range pl.Phases {
		if !phase.Enabled {
			continue
		}
		cfmt.Print("\n")

		cfmt.Println("{{  " + strings.ToUpper(phase.Name) + "}}::bold|white" + " (" + phase.Context + ")")

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

		printCLI(table, pipelineResp)
	}
}
