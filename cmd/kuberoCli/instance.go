package kuberoCli

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

import (
	"fmt"
	"os"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// instanceCmd represents the instance command
var instanceCmd = &cobra.Command{
	Use:     "instance",
	Aliases: []string{"i"},
	Short:   "List available instances",
	Long:    `Print a list of available instances.`,
	Run: func(cmd *cobra.Command, args []string) {
		/*
			fmt.Println("current instance : " + currentInstanceName)
			fmt.Println(instanceList)
			if currentInstance.ApiUrl == "" {
				fmt.Println("No current instance api URL")
			} else {
				fmt.Println("Current instance api URL : " + currentInstance.ApiUrl)
			}
		*/

		printInstanceList()
	},
}

func init() {
	rootCmd.AddCommand(instanceCmd)
}

func printInstanceList() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Active", "Token", "Name", "API URL", "Path", "IAC Base Dir"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	//table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	//table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowLine(true)
	table.SetCenterSeparator("")
	//table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)
	for _, instanceName := range instanceNameList {
		active := ""
		if instanceName == currentInstanceName {
			active = cfmt.Sprintf("   {{✔}}::green")
		}

		token := ""

		// check if instanceName is in credentialsConfig
		if credentialsConfig.GetString(instanceName) != "" {
			token = cfmt.Sprintf("  {{✔}}::green")
		}

		table.Append([]string{
			active,
			token,
			instanceName,
			instanceList[instanceName].ApiUrl,
			instanceList[instanceName].ConfigPath,
			instanceList[instanceName].IacBaseDir,
		})
	}
	table.Render()
}

func createInstanceForm() {
	fmt.Println("Create a new instance")

	instanceName := promptLine("Enter the name of the instance", "", "")
	instanceApiurl := promptLine("Enter the API URL of the instance", "", "http://localhost:80")
	instancePath := viper.ConfigFileUsed()

	personalInstanceList := viper.GetStringMap("instances")

	personalInstanceList[instanceName] = Instance{
		Name:       instanceName,
		ApiUrl:     instanceApiurl,
		ConfigPath: instancePath,
	}

	viper.Set("instances", personalInstanceList)

	instanceNameList = append(instanceNameList, instanceName)

	setCurrentInstance(instanceName)

}

func setCurrentInstance(instanceName string) {
	currentInstanceName = instanceName
	currentInstance = instanceList[instanceName]
	viper.Set("currentInstance", instanceName)
	writeConfigErr := viper.WriteConfig()
	if writeConfigErr != nil {
		fmt.Println("Failed to save configuration:", writeConfigErr)
		return
	}
}

func deleteInstanceForm() {
	instanceName := selectFromList("Select an instance to delete", instanceNameList, "")

	delete(instanceList, instanceName)
	viper.Set("instances", instanceList)
	writeConfigErr := viper.WriteConfig()
	if writeConfigErr != nil {
		fmt.Println("Failed to save configuration:", writeConfigErr)
		return
	}

}
