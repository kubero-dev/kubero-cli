package version

import (
	_ "embed"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of kuberoCli",
		Long:  "Print the version number of kuberoCli",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(GetVersionInfo())
		},
	}
	subLatestCmd = &cobra.Command{
		Use:   "latest",
		Short: "Print the latest version number of kuberoCli",
		Long:  "Print the latest version number of kuberoCli",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(GetLatestVersionInfo())
		},
	}
	subCmdCheck = &cobra.Command{
		Use:   "check",
		Short: "Check if the current version is the latest version of kuberoCli",
		Long:  "Check if the current version is the latest version of kuberoCli",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(GetVersionInfoWithLatestAndCheck())
		},
	}
	subCmdUpgrade = &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade kuberoCli to the latest version",
		Long:  "Upgrade kuberoCli to the latest version",
		Run: func(cmd *cobra.Command, args []string) {
			syncCmd := CmdUpgradeCLIAndCheck()
			if err := syncCmd.Start(); err != nil {
				fmt.Println("Error: " + err.Error())
			}
		},
	}
	subCmdUpgradeCheck = &cobra.Command{
		Use:    "check",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return UpgradeCLI()
		},
	}
)

const gitModelUrl = "https://github.com/kubero-dev/kubero-cli"
const currentVersionFallback = "v2.4.2" // First version with the version file

//go:embed CLI_VERSION
var currentVersion string

func GetVersion() string {
	if currentVersion == "" {
		return currentVersionFallback
	}
	return currentVersion
}

func GetGitModelUrl() string {
	return gitModelUrl
}

func GetVersionInfo() string {
	return "Version: " + GetVersion() + "\n" + "Git repository: " + GetGitModelUrl()
}

func GetLatestVersionFromGit() string {
	netClient := &http.Client{
		Timeout: time.Second * 10,
	}

	response, err := netClient.Get(gitModelUrl + "/releases/latest")
	if err != nil {
		return "Error: " + err.Error()
	}

	if response.StatusCode != 200 {
		return "Error: " + response.Status
	}

	tag := strings.Split(response.Request.URL.Path, "/")

	return tag[len(tag)-1]
}

func GetLatestVersionInfo() string {
	return "Latest version: " + GetLatestVersionFromGit()
}

func GetVersionInfoWithLatestAndCheck() string {
	if GetVersion() == GetLatestVersionFromGit() {
		return GetVersionInfo() + "\n" + GetLatestVersionInfo() + "\n" + "You are using the latest version."
	} else {
		return GetVersionInfo() + "\n" + GetLatestVersionInfo() + "\n" + "You are using an outdated version.\n" + "Please upgrade your kuberoCli to prevent any issues."
	}
}

func UpgradeCLI() error {
	if GetVersion() == GetLatestVersionFromGit() {
		fmt.Println("You are using the latest version.")
		return nil
	} else {
		netClient := &http.Client{
			Timeout: time.Second * 10,
		}

		response, err := netClient.Get(gitModelUrl + "/releases/latest")
		if err != nil {
			return err
		}

		latestVersion := response.Status

		if latestVersion == "200 OK" {
			return fmt.Errorf("error: %s", latestVersion)
		}

		fileUrl := gitModelUrl + "/releases/download/" + latestVersion + "/kuberoCli"

		// Download the file
		response, err = netClient.Get(fileUrl)
		if err != nil {
			return fmt.Errorf("error: %s", err)
		}

		// Save the file
		writeFile, err := os.Create("kuberoCli")
		if err != nil {
			return fmt.Errorf("error: %s", err)
		}
		defer func(writeFile *os.File) {
			_ = writeFile.Close()
		}(writeFile)

		fileInfo, err := writeFile.Stat()
		if err != nil {
			return fmt.Errorf("error: %s", err)
		}

		fileMode := fileInfo.Mode()
		if err := writeFile.Chmod(fileMode); err != nil {
			return fmt.Errorf("error: %s", err)
		}

		currentExecutable, err := os.Executable()
		if err != nil {
			return fmt.Errorf("error: %s", err)
		}

		cmdCopy := "cp " + currentExecutable + " " + currentExecutable + ".old"
		cmdRemove := "rm " + currentExecutable
		cmdRename := "mv kuberoCli " + currentExecutable
		cmdUpgradeSpawner := cmdCopy + " && " + cmdRemove + " && " + cmdRename

		spawner := os.Getenv("SHELL")
		if spawner == "" {
			spawner = "/bin/sh"
		}

		cmd := exec.Command(spawner, "-c", cmdUpgradeSpawner)
		return cmd.Run()
	}
}

func CmdUpgradeCLIAndCheck() *exec.Cmd {
	cmd := exec.Command("kubero", "upgrade", "check")
	return cmd
}

func CliCommand() *cobra.Command {
	versionCmd.AddCommand(subLatestCmd)
	versionCmd.AddCommand(subCmdCheck)
	//subCmdUpgrade.AddCommand(subCmdUpgradeCheck)
	//versionCmd.AddCommand(subCmdUpgrade)
	return versionCmd
}
