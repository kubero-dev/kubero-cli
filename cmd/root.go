/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var outputFormat string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubero",
	Short: "A brief description of your application",
	Long: `
	 ___ ___       __               _______
	|   Y   .--.--|  |--.-----.----|   _   |
	|.  1  /|  |  |  _  |  -__|   _|.  |   |
	|.  _  \|_____|_____|_____|__| |.  |   |
	|:  |   \                      |:  1   |
	|::.| .  )         CLI         |::.. . |
	'--- ---'                      '-------'

	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "output format [table, json]")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func printCLI(table *tablewriter.Table, r *resty.Response) {
	if outputFormat == "json" {
		fmt.Println(r)
	} else {
		table.Render()
	}
}

func promptLine(question string) string {
	reader := bufio.NewReader(os.Stdin)
	cfmt.Printf("\n  {{%s}}::lightWhite : ", question)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	return text
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

type Repositories struct {
	Github    bool `json:"github"`
	Gitea     bool `json:"gitea"`
	Gitlab    bool `json:"gitlab"`
	Bitbucket bool `json:"bitbucket"`
	Docker    bool `json:"docker"`
}

var repoSimpleList []string

func loadRepositories() {

	rep, _ := client.Get("/api/cli/config/repositories")

	var availRep Repositories
	json.Unmarshal(rep.Body(), &availRep)

	t := reflect.TypeOf(availRep)

	repoSimpleList = make([]string, t.NumField())
	for i := range repoSimpleList {
		if reflect.ValueOf(availRep).Field(i).Bool() {
			repoSimpleList[i] = t.Field(i).Name
		}
	}
}

type Contexts []struct {
	Cluster string `json:"cluster"`
	Name    string `json:"name"`
	User    string `json:"user"`
}

var contextSimpleList []string

func loadContexts() {

	cont, _ := client.Get("/api/cli/config/k8s/context")

	var contexts Contexts
	json.Unmarshal(cont.Body(), &contexts)

	for _, context := range contexts {
		contextSimpleList = append(contextSimpleList, context.Name)
	}
}
