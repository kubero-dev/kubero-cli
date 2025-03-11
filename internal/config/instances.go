package config

import (
	"fmt"
	"github.com/i582/cfmt/cmd/cfmt"
	u "github.com/kubero-dev/kubero-cli/internal/utils"
	t "github.com/kubero-dev/kubero-cli/types"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
	"os"
)

var (
	utilsPrompt    = u.NewConsolePrompt()
	selectFromList = utilsPrompt.SelectFromList
)

type InstanceManager struct {
	credentialsCfg       *viper.Viper
	currentInstance      *t.Instance
	personalInstanceList map[string]*t.Instance
	globalInstanceList   []*t.Instance
}

func NewInstanceManager(credentialsCfg *viper.Viper) *InstanceManager {
	return &InstanceManager{
		credentialsCfg: credentialsCfg,
	}
}

func (i *InstanceManager) CreateInstanceForm() error {
	fmt.Println("Create a new instance")

	instanceNameArg := promptLine("Enter the name of the instance", "", "")
	instanceApiurlArg := promptLine("Enter the API URL of the instance", "", "http://localhost:80")
	instancePathArg := viper.ConfigFileUsed()
	personalInstanceList := viper.GetStringMap("instances")

	if i.personalInstanceList == nil {
		i.personalInstanceList = make(map[string]*t.Instance)
	}
	i.personalInstanceList[instanceNameArg] = &t.Instance{
		Name:       instanceNameArg,
		ApiUrl:     instanceApiurlArg,
		ConfigPath: instancePathArg,
	}

	viper.Set("instances", personalInstanceList)
	return i.SetCurrentInstance(instanceNameArg)
}

func (i *InstanceManager) PrintInstanceList() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Active", "Token", "Name", "API URL", "Path", "IAC Base Dir"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetRowLine(true)
	table.SetCenterSeparator("")
	table.SetRowSeparator("")
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	instanceNameList := i.GetInstanceNameList()
	for _, instanceName := range instanceNameList {
		active := ""
		if instanceName == i.currentInstance.Name {
			active = cfmt.Sprintf("   {{✔}}::green")
		}

		token := ""

		if i.credentialsCfg.GetString(instanceName) != "" {
			token = cfmt.Sprintf("  {{✔}}::green")
		}

		table.Append([]string{
			active,
			token,
			instanceName,
			i.personalInstanceList[instanceName].ApiUrl,
			i.personalInstanceList[instanceName].ConfigPath,
			i.personalInstanceList[instanceName].IacBaseDir,
		})
	}
	table.Render()
}

func (i *InstanceManager) SetCurrentInstance(instanceName string) error {
	currentInstanceName := instanceName
	currentInstance := i.personalInstanceList[currentInstanceName]
	viper.Set("currentInstance", currentInstance.Name)
	writeConfigErr := viper.WriteConfig()
	if writeConfigErr != nil {
		fmt.Println("Failed to save configuration:", writeConfigErr)
		return writeConfigErr
	}

	i.currentInstance = currentInstance

	return nil
}

func (i *InstanceManager) DeleteInstanceForm() error {
	instanceList := viper.GetStringMap("instances")
	instanceNameList := i.GetInstanceNameList()

	instanceName := selectFromList("Select an instance to delete", instanceNameList, "")

	delete(instanceList, instanceName)
	viper.Set("instances", instanceList)
	writeConfigErr := viper.WriteConfig()
	if writeConfigErr != nil {
		fmt.Println("Failed to save configuration:", writeConfigErr)
		return writeConfigErr
	}

	return nil
}

func (i *InstanceManager) GetInstanceNameList() []string {
	var instanceNameList = make([]string, 0)
	for _, instance := range i.globalInstanceList {
		instanceNameList = append(instanceNameList, instance.Name)
	}
	return instanceNameList
}

func (i *InstanceManager) GetInstance(instanceName string) *t.Instance {
	if instance, ok := i.personalInstanceList[instanceName]; ok {
		return instance
	}
	return &t.Instance{Name: instanceName}
}

func (i *InstanceManager) GetCurrentInstance() *t.Instance {
	if i.currentInstance == nil {
		currentInstanceName := viper.GetString("currentInstance")
		i.currentInstance = i.GetInstance(currentInstanceName)
	}
	return i.currentInstance
}

func (i *InstanceManager) GetPersonalInstanceList() map[string]*t.Instance {
	return i.personalInstanceList
}

func (i *InstanceManager) GetGlobalInstanceList() []*t.Instance {
	return i.globalInstanceList
}

func (i *InstanceManager) EnsureInstanceOrCreate() error {
	if i.currentInstance == nil {
		currentInstanceName := viper.GetString("currentInstance")
		if currentInstanceName == "" {
			return i.CreateInstanceForm()
		}
		i.currentInstance = i.GetInstance(currentInstanceName)
	}
	return nil
}
