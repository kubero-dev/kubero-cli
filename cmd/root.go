package cmd

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-resty/resty/v2"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var outputFormat string
var force bool

//go:embed VERSION
var version string

var pipelineConfig *viper.Viper
var pipelineConfigList []*viper.Viper

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "kubero",
	Short:   "A brief description of your application",
	Version: version,
	Long: `

	,--. ,--.        ,--.
	|  .'   /,--.,--.|  |-.  ,---. ,--.--. ,---.
	|  .   ' |  ||  || .-. '| .-. :|  .--'| .-. |
	|  |\   \'  ''  '| '-' |\   --.|  |   ' '-' '
	'--' '--' '----'  '---'  '----''--'    '---'


Documentation:
  https://docs.kubero.dev

Env vars:
  KUBERO_API_URL - Kubero api URL
  KUBERO_API_TOKEN - Kubero api token
`,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	getPipelinesInRepository()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
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

// question, options/example, default
func promptLine(question string, options string, def string) string {
	if def != "" && force {
		cfmt.Printf("\n  %s %s : {{%s}}::green\n", question, options, def)
		return def
	}
	reader := bufio.NewReader(os.Stdin)
	cfmt.Printf("\n  %s %s {{%s}}::green : ", question, options, def)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	if text == "" {
		text = def
	}
	return text
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

func getGitRemote() string {
	gitdir := GetGitdir()
	fs := osfs.New(gitdir)
	s := filesystem.NewStorageWithOptions(fs, cache.NewObjectLRUDefault(), filesystem.Options{KeepDescriptors: true})
	r, err := git.Open(s, fs)
	if err == nil {
		remotes, _ := r.Remotes()
		return remotes[0].Config().URLs[0]
	}
	return ""
}

func GetGitdir() string {
	wd, _ := os.Getwd()

	path := strings.Split(wd, "/")
	for i := len(path); i >= 0; i-- {

		subpath := strings.Join(path[:i], "/")
		fileInfo, err := os.Stat(subpath + "/.git")

		if err != nil {
			//fmt.Println(subpath + "/.git not a dir")
			continue
		} else {
			if fileInfo.IsDir() {
				//fmt.Println(subpath + "/.git is a dir")
				return strings.Join(path[:i], "/") + "/.git"
			} else {
				//fmt.Println(subpath + "/.git not a dir")
				continue
			}
		}

	}
	return ""
}

func getPipelinesInRepository() {

	gitdir := GetGitdir() + "/../.kubero"

	// iterate over all directories and append pipeline.yaml config to pipelineConfig
	err := filepath.Walk(gitdir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			//fmt.Println("Checking dir: " + path)
			if _, err := os.Stat(path + "/pipeline.yaml"); err == nil {
				//fmt.Println("Found pipeline.yaml in: " + path)
				pc := viper.New()
				pc.SetConfigName("pipeline") // name of config file (without extension)
				pc.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name
				pc.AddConfigPath(path)
				pc.ReadInConfig()

				pipelineConfigList = append(pipelineConfigList, pc)
				//fmt.Println("Using config file:", pc.ConfigFileUsed())
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	//fmt.Println("Pipeline Configs: " + strconv.Itoa(len(pipelineConfigList)))
}

func checkConfig() {
	if viper.GetString("api.url") == "" {
		fmt.Println("No API URL set. Create a .kubero/credentials.yaml file in your repository or set the KUBERO_API_URL environment variable.")
		os.Exit(1)
	}

	if viper.GetString("api.token") == "" {
		fmt.Println("No API token set. Create a .kubero/credentials.yaml file in your repository or set the KUBERO_API_TOKEN environment variable.")
		os.Exit(1)
	}
}
