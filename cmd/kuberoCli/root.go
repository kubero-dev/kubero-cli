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

var currentInstanceName string
var instanceList map[string]Instance
var instanceNameList []string
var currentInstance Instance = Instance{}

//go:embed VERSION
var version string

var pipelineConfig *viper.Viper
var credentialsConfig *viper.Viper

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
	loadCredentials()
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
}

func printCLI(table *tablewriter.Table, r *resty.Response) {
	if outputFormat == "json" {
		fmt.Println(r)
	} else {
		table.Render()
	}
}

func promptWarning(msg string) {
	cfmt.Println("{{\n⚠️   " + msg + ".\n}}::yellow")
}

// promptBanner("✖ ERROR ..... do something")
func promptBanner(msg string) {
	cfmt.Printf(`
    {{                                                                            }}::bgRed
    {{  %-72s  }}::bgRed|#ffffff
    {{                                                                            }}::bgRed
	
	`, msg)
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

func selectFromList(question string, options []string, def string) string {
	cfmt.Println("")
	if def != "" && force {
		cfmt.Printf("\n{{?}}::green %s : {{%s}}::cyan\n", question, def)
		return def
	}
	prompt := &survey.Select{
		Message: question,
		Options: options,
	}
	survey.AskOne(prompt, &def)
	return def
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

	//load a personal config from the user's home directory

	gitdir := getGitdir()
	dir := gitdir

	repoConfig := viper.New()
	repoConfig.SetConfigName("kubero") // name of config file (without extension)
	repoConfig.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	repoConfig.AddConfigPath(dir)      // TODO this should search for the git repo root
	repoConfig.ConfigFileUsed()
	errCred := repoConfig.ReadInConfig()

	//load a default config from the current local git repository
	viper.SetDefault("api.url", "http://default:2000")
	viper.SetConfigName("kubero")         // name of config file (without extension)
	viper.SetConfigType("yaml")           // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/kubero/")   // path to look for the config file in
	viper.AddConfigPath("$HOME/.kubero/") // call multiple times to add many search paths
	err := viper.ReadInConfig()

	if err != nil && errCred != nil {

		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found; using defaults")
		} else {
			fmt.Printf("Error while loading config files: %v \n\n\n%v", err, errCred)
		}
	}

	viper.UnmarshalKey("instances", &instanceList)

	// iterate over all instances and and set the config path
	for instanceName, instance := range instanceList {
		instance.Name = instanceName
		instance.ConfigPath = viper.ConfigFileUsed()
		instanceList[instanceName] = instance
	}

	var repoInstancesList map[string]Instance
	repoConfig.UnmarshalKey("instances", &repoInstancesList)

	for instanceName, repoInstance := range repoInstancesList {
		repoInstance.Name = instanceName
		repoInstance.ConfigPath = repoConfig.ConfigFileUsed()
		instanceList[instanceName] = repoInstance
	}

	currentInstanceName = viper.GetString("currentInstance")

	// iterate over all instances and find the current one
	for instanceName, instance := range instanceList {
		instance.Name = instanceName
		instanceNameList = append(instanceNameList, instanceName)
		if instanceName == currentInstanceName {
			currentInstance = instance
		}
	}

}

func loadCredentials() {

	//load a personal config from the user's home directory
	credentialsConfig = viper.New()
	credentialsConfig.SetConfigName("credentials")    // name of config file (without extension)
	credentialsConfig.SetConfigType("yaml")           // REQUIRED if the config file does not have the extension in the name
	credentialsConfig.AddConfigPath("/etc/kubero/")   // path to look for the config file in
	credentialsConfig.AddConfigPath("$HOME/.kubero/") // call multiple times to add many search paths
	err := credentialsConfig.ReadInConfig()

	if err != nil {
		fmt.Println("Error while loading credentialsConfig file:", err)
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

	if appName == "" {
		fmt.Println("")

		prompt := &survey.Select{
			Message: "Select an app",
			Options: availableApps,
		}
		survey.AskOne(prompt, &appName)
	}
}
