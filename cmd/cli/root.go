package cli

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/faelmori/kubero-cli/pkg/kuberoApi"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-resty/resty/v2"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	_ "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	outputFormat        string
	force               bool
	repoSimpleList      []string
	client              *resty.Request
	api                 *kuberoApi.KuberoClient
	contextSimpleList   []string
	currentInstanceName string
	instanceList        map[string]Instance
	instanceNameList    []string
	currentInstance     Instance
	kuberoCliVersion    string
	pipelineConfig      *viper.Viper
	credentialsConfig   *viper.Viper
	db                  *gorm.DB
)

var rootCmd = &cobra.Command{
	Use:   "kubero",
	Short: "Kubero is a platform as a service (PaaS) that enables developers to build, run, and operate applications on Kubernetes.",
	Long: `
	,--. ,--.        ,--.
	|  .'   /,--.,--.|  |-.  ,---. ,--.--. ,---.
	|  .   ' |  ||  || .-. '| .-. :|  .--'| .-. |
	|  |\   \'  ''  '| '-' |\   --.|  |   ' '-' '
	'--' '--' '----'  '---'  '----''--'    '---'
Documentation:
  https://docs.kubero.dev
`,
	Example: `kubero install`,
	Aliases: []string{"kbr"},
}

func Execute() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = false

	SetUsageDefinition(rootCmd)

	loadCLIConfig()
	loadCredentials()
	api = new(kuberoApi.KuberoClient)
	client = api.Init(currentInstance.ApiUrl, credentialsConfig.GetString(currentInstanceName))
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	initDB()
}

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("kubero.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	autoMigrateErr := db.AutoMigrate(&Instance{})
	if autoMigrateErr != nil {
		log.Fatal("Failed to migrate database:", autoMigrateErr)
		return
	}
}

func printCLI(table *tablewriter.Table, r *resty.Response) {
	if outputFormat == "json" {
		fmt.Println(r)
	} else {
		table.Render()
	}
}

func promptWarning(msg string) {
	_, _ = cfmt.Println("{{\n⚠️   " + msg + ".\n}}::yellow")
}

//func promptBanner(msg string) {
//	_, _ = cfmt.Printf(`
//    {{                                                                            }}::bgRed
//    {{  %-72s  }}::bgRed|#ffffff
//    {{                                                                            }}::bgRed
//	`, msg)
//}

func promptLine(question, options, def string) string {
	if def != "" && force {
		_, _ = cfmt.Printf("\n{{?}}::green %s %s : {{%s}}::cyan\n", question, options, def)
		return def
	}
	reader := bufio.NewReader(os.Stdin)
	_, _ = cfmt.Printf("\n{{?}}::green|bold {{%s %s}}::bold {{%s}}::cyan : ", question, options, def)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" {
		text = def
	}
	return text
}

func selectFromList(question string, options []string, def string) string {
	_, _ = cfmt.Println("")
	if def != "" && force {
		_, _ = cfmt.Printf("\n{{?}}::green %s : {{%s}}::cyan\n", question, def)
		return def
	}
	prompt := &survey.Select{
		Message: question,
		Options: options,
	}
	askOneErr := survey.AskOne(prompt, &def)
	if askOneErr != nil {
		fmt.Println("Error while selecting:", askOneErr)
		return ""
	}
	return def
}

func confirmationLine(question, def string) bool {
	confirmation := promptLine(question, "[y,n]", def)
	if confirmation != "y" {
		_, _ = cfmt.Println("{{\n✗ Aborted\n}}::red")
		os.Exit(0)
		return false
	}
	return true
}

func loadRepositories() {
	res, err := api.GetRepositories()
	if res == nil {
		fmt.Println("Error: Can't reach Kubero API. Make sure, you are logged in.")
		os.Exit(1)
	}
	if res.StatusCode() != 200 {
		fmt.Println("Error:", res.StatusCode(), "Can't reach Kubero API. Make sure, you are logged in.")
		os.Exit(1)
	}
	if err != nil {
		fmt.Println("Error: Unable to load repositories")
		fmt.Println(err)
		os.Exit(1)
	}
	var availRep Repositories
	jsonUnmarshalErr := json.Unmarshal(res.Body(), &availRep)
	if jsonUnmarshalErr != nil {
		fmt.Println("Error: Unable to load repositories")
		return
	}
	t := reflect.TypeOf(availRep)
	repoSimpleList = make([]string, t.NumField())
	for i := range repoSimpleList {
		if reflect.ValueOf(availRep).Field(i).Bool() {
			repoSimpleList[i] = t.Field(i).Name
		}
	}
}

func loadContexts() {
	cont, _ := api.GetContexts()
	var contexts Contexts
	jsonUnmarshalErr := json.Unmarshal(cont.Body(), &contexts)
	if jsonUnmarshalErr != nil {
		fmt.Println("Error: Unable to load contexts")
		return
	}
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
		subPath := strings.Join(path[:i], "/")
		fileInfo, err := os.Stat(subPath + "/.git")
		if err == nil && fileInfo.IsDir() {
			return subPath
		}
	}
	return ""
}

func getIACBaseDir() string {
	basePath := "."
	if currentInstance.IacBaseDir == "" {
		currentInstance.IacBaseDir = ".kubero"
		basePath += "/" + currentInstance.IacBaseDir
	}
	gitdir := getGitdir()
	if gitdir != "" {
		basePath = gitdir + "/" + currentInstance.IacBaseDir
	}
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		_, _ = cfmt.Println("{{Creating directory}}::yellow " + basePath)
		mkDirAllErr := os.MkdirAll(basePath, 0755)
		if mkDirAllErr != nil {
			fmt.Println("Error while creating directory:", mkDirAllErr)
			return ""
		}
	}
	return basePath
}

//func loadConfigs(basePath, pipelineName string) {
//	baseDir := getIACBaseDir()
//	dir := baseDir + "/" + pipelineName
//	pipelineConfig = viper.New()
//	pipelineConfig.SetConfigName("pipeline")
//	pipelineConfig.SetConfigType("yaml")
//	pipelineConfig.AddConfigPath(dir)
//	readInConfigErr := pipelineConfig.ReadInConfig()
//	if readInConfigErr != nil {
//		fmt.Println("Error while loading pipeline config file:", readInConfigErr)
//		return
//	}
//}

func loadConfigs(pipelineName string) {
	baseDir := getIACBaseDir()
	dir := baseDir + "/" + pipelineName
	pipelineConfig = viper.New()
	pipelineConfig.SetConfigName("pipeline")
	pipelineConfig.SetConfigType("yaml")
	pipelineConfig.AddConfigPath(dir)
	readInConfigErr := pipelineConfig.ReadInConfig()
	if readInConfigErr != nil {
		fmt.Println("Error while loading pipeline config file:", readInConfigErr)
		return
	}
}

func loadCLIConfig() {
	dir := getGitdir()
	repoConfig := viper.New()
	repoConfig.SetConfigName("kubero")
	repoConfig.SetConfigType("yaml")
	repoConfig.AddConfigPath(dir)
	repoConfig.ConfigFileUsed()
	errCred := repoConfig.ReadInConfig()

	viper.SetDefault("api.url", "http://default:2000")
	viper.SetConfigName("kubero")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/kubero/")
	viper.AddConfigPath("$HOME/.kubero/")
	err := viper.ReadInConfig()

	if err != nil && errCred != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			fmt.Println("No config file found; using defaults")
		} else {
			fmt.Println("Error while loading config file:", err)
			return
		}
	}

	viperUnmarshalErr := viper.UnmarshalKey("instances", &instanceList)
	if viperUnmarshalErr != nil {
		fmt.Println("Error while unmarshalling instances:", viperUnmarshalErr)
		return
	}
	for instanceName, instance := range instanceList {
		instance.Name = instanceName
		instance.ConfigPath = viper.ConfigFileUsed()
		instanceList[instanceName] = instance
	}

	var repoInstancesList map[string]Instance
	unmarshalKeyErr := repoConfig.UnmarshalKey("instances", &repoInstancesList)
	if unmarshalKeyErr != nil {
		fmt.Println("Error while unmarshalling instances:", unmarshalKeyErr)
		return
	}
	for instanceName, repoInstance := range repoInstancesList {
		repoInstance.Name = instanceName
		repoInstance.ConfigPath = repoConfig.ConfigFileUsed()
		instanceList[instanceName] = repoInstance
	}

	currentInstanceName = viper.GetString("currentInstance")
	for instanceName, instance := range instanceList {
		instance.Name = instanceName
		instanceNameList = append(instanceNameList, instanceName)
		if instanceName == currentInstanceName {
			currentInstance = instance
		}
	}
}

func loadCredentials() {
	credentialsConfig = viper.New()
	credentialsConfig.SetConfigName("credentials")
	credentialsConfig.SetConfigType("yaml")
	credentialsConfig.AddConfigPath("/etc/kubero/")
	credentialsConfig.AddConfigPath("$HOME/.kubero/")
	err := credentialsConfig.ReadInConfig()
	if err != nil {
		fmt.Println("Error while loading credentialsConfig file:", err)
	}
}

func boolToEmoji(b bool) string {
	if b {
		return "✅"
	}
	return "❌"
}

func ensurePipelineIsSet(pipelinesList []string) {
	if pipelineName == "" {
		fmt.Println("")
		prompt := &survey.Select{
			Message: "Select a pipeline",
			Options: pipelinesList,
		}
		askOneErr := survey.AskOne(prompt, &pipelineName)
		if askOneErr != nil {
			fmt.Println("Error while selecting pipeline:", askOneErr)
			return
		}
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
		askOneErr := survey.AskOne(prompt, &stageName)
		if askOneErr != nil {
			fmt.Println("Error while selecting stage:", askOneErr)
			return
		}
	}
}

func ensureAppNameIsSelected(availableApps []string) {
	if appName == "" {
		fmt.Println("")
		prompt := &survey.Select{
			Message: "Select an app",
			Options: availableApps,
		}
		askOneErr := survey.AskOne(prompt, &appName)
		if askOneErr != nil {
			fmt.Println("Error while selecting app:", askOneErr)
			return
		}
	}
}
