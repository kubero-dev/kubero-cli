package kuberoCli

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"

	"kubero/pkg/kuberoApi"
	"os"
	"reflect"
	"strings"

	"github.com/AlecAivazis/survey/v2"
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
var repoSimpleList []string
var client *resty.Request
var api *kuberoApi.KuberoClient
var contextSimpleList []string

//go:embed VERSION
var version string

var pipelineConfig *viper.Viper

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "kubero",
	Short:   "Kubero is a platform as a service (PaaS) that enables developers to build, run, and operate applications on Kubernetes.",
	Version: version,
	Long: `

	,--. ,--.        ,--.
	|  .'   /,--.,--.|  |-.  ,---. ,--.--. ,---.
	|  .   ' |  ||  || .-. '| .-. :|  .--'| .-. |
	|  |\   \'  ''  '| '-' |\   --.|  |   ' '-' '
	'--' '--' '----'  '---'  '----''--'    '---'


Documentation:
  https://docs.kubero.dev
`,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	loadCLIConfig()
	api = new(kuberoApi.KuberoClient)
	client = api.Init(viper.GetString("api.url"), viper.GetString("api.token"))

	//client = kuberoApi.Init(viper.GetString("api.url"), viper.GetString("api.token"))
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

	//rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "output format [table, json]")

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
		cfmt.Printf("\n{{?}}::green %s %s : {{%s}}::cyan\n", question, options, def)
		return def
	}
	reader := bufio.NewReader(os.Stdin)
	cfmt.Printf("\n{{?}}::green|bold {{%s %s}}::bold {{%s}}::cyan : ", question, options, def)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	if text == "" {
		text = def
	}
	return text
}

func confirmationLine(question string, def string) bool {
	confirmation := promptLine(question, "[y,n]", def)
	if confirmation != "y" {
		cfmt.Println("{{\n✗ Aborted\n}}::red")
		os.Exit(0)
		return false
	} else {
		return true
	}
}

func loadRepositories() {

	res, err := client.Get("/api/cli/config/repositories")
	if res.StatusCode() != 200 {
		fmt.Println("Error: ", res.StatusCode(), "Can't reach Kubero API. Make sure, you are logged in.")
		os.Exit(1)
	}

	if err != nil {
		fmt.Println("Error: ", "Unable to load repositories")
		fmt.Println(err)
		os.Exit(1)
	}

	var availRep Repositories
	json.Unmarshal(res.Body(), &availRep)

	t := reflect.TypeOf(availRep)

	repoSimpleList = make([]string, t.NumField())
	for i := range repoSimpleList {
		if reflect.ValueOf(availRep).Field(i).Bool() {
			repoSimpleList[i] = t.Field(i).Name
		}
	}
}

func loadContexts() {

	cont, _ := client.Get("/api/cli/config/k8s/context")

	var contexts Contexts
	json.Unmarshal(cont.Body(), &contexts)

	for _, context := range contexts {
		contextSimpleList = append(contextSimpleList, context.Name)
	}
}

func getGitRemote() string {
	gitdir := getGitdir() + "/.git"
	fs := osfs.New(gitdir)
	s := filesystem.NewStorageWithOptions(fs, cache.NewObjectLRUDefault(), filesystem.Options{KeepDescriptors: true})
	r, err := git.Open(s, fs)
	if err == nil {
		remotes, _ := r.Remotes()
		return remotes[0].Config().URLs[0]
	}
	return ""
}

func getGitdir() string {
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
				return strings.Join(path[:i], "/")
			} else {
				//fmt.Println(subpath + "/.git not a dir")
				continue
			}
		}

	}
	return ""
}

func loadConfigs(basePath string, pipelineName string) {

	gitdir := getGitdir()
	dir := gitdir + basePath + pipelineName
	//fmt.Println(dir)

	pipelineConfig = viper.New()
	pipelineConfig.SetConfigName("pipeline") // name of config file (without extension)
	pipelineConfig.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name
	pipelineConfig.AddConfigPath(dir)        // path to look for the config file in
	pipelineConfig.ReadInConfig()

	//fmt.Println("Using config file:", viper.ConfigFileUsed())
	//fmt.Println("Using config file:", pipelineConfig.ConfigFileUsed())
}

// create recursive folder if not exists
/*
func createFolder(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
}
*/

func loadCLIConfig() {

	gitdir := getGitdir()
	basePath := "/.kubero/" //TODO Make it dynamic
	dir := gitdir + basePath

	//load a default config from the current local git repository
	viper.SetDefault("api.url", "http://localhost:2000")
	viper.SetConfigName("kubero") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(dir)      // TODO this should search for the git repo root
	err := viper.ReadInConfig()

	//load a personal config from the user's home directory
	personal := viper.New()
	personal.SetConfigName("kubero")        // name of config file (without extension)
	personal.SetConfigType("yaml")          // REQUIRED if the config file does not have the extension in the name
	personal.AddConfigPath("/etc/kubero/")  // path to look for the config file in
	personal.AddConfigPath("$HOME/.kubero") // call multiple times to add many search paths
	errCred := personal.ReadInConfig()

	viper.MergeConfigMap(personal.AllSettings())
	if err != nil && errCred != nil {

		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found; using defaults")
		} else {
			fmt.Printf("Error while loading config files: %v \n\n\n%v", err, errCred)
		}
	}
}

func boolToEmoji(b bool) string {
	if b {
		return "✅"
	} else {
		return "❌"
	}
}

// pipelinesList := getAllLocalPipelines()
func ensurePipelineIsSet(pipelinesList []string) {
	if pipelineName == "" {
		fmt.Println("")
		prompt := &survey.Select{
			Message: "Select a pipeline",
			Options: pipelinesList,
		}
		survey.AskOne(prompt, &pipelineName)
	}
}

func ensureAppNameIsSet() {
	if appName == "" {
		appName = promptLine("Define a app name", "", appName)
	}
}

func ensureStageNameIsSet() {
	if stageName == "" {
		fmt.Println("")

		pipelineConfig := loadPipelineConfig(pipelineName)
		availablePhases := getPipelinePhases(pipelineConfig)

		prompt := &survey.Select{
			Message: "Select a stage",
			Options: availablePhases,
		}
		survey.AskOne(prompt, &stageName)
	}
}

func ensureAppNameIsSelected(availableApps []string) {

	prompt := &survey.Select{
		Message: "Select a app",
		Options: availableApps,
	}
	survey.AskOne(prompt, &appName)
}

func getAllRemoteApps() []string {
	apps, _ := api.GetApps()
	var appShortList []appShort
	json.Unmarshal(apps.Body(), &appShortList)

	var appsList []string
	for _, app := range appShortList {
		if pipelineName != "" && app.Pipeline != pipelineName {
			continue
		}
		if stageName != "" && app.Phase != stageName {
			continue
		}
		if appName != "" && app.Name != appName {
			continue
		}
		appsList = append(appsList, app.Name)
	}

	return appsList
}
