package kuberoCli

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func colorYellow(s string) string {
	return color.New(color.FgYellow).SprintFunc()(s)
}

func colorGreen(s string) string {
	return color.New(color.FgGreen).SprintFunc()(s)
}

func colorBlue(s string) string {
	return color.New(color.FgBlue).SprintFunc()(s)
}

func colorRed(s string) string {
	return color.New(color.FgRed).SprintFunc()(s)
}

func colorHelp(s string) string {
	return color.New(color.FgCyan).SprintFunc()(s)
}

func hasServiceCommands(cmds []*cobra.Command) bool {
	for _, cmd := range cmds {
		if cmd.Annotations["service"] == "true" {
			return true
		}
	}
	return false
}

func hasModuleCommands(cmds []*cobra.Command) bool {
	for _, cmd := range cmds {
		if cmd.Annotations["service"] != "true" {
			return true
		}
	}
	return false
}

var cliUsageTemplate = `{{colorYellow "Usage:"}}{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command] [args]{{end}}{{if gt (len .Aliases) 0}}

{{colorYellow "Aliases:"}}
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

{{colorYellow "Example:"}}
  {{.Example}}{{end}}{{if .HasAvailableSubCommands}}

{{colorYellow "Available Commands:"}}
{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{colorGreen (rpad .Name .NamePadding) }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

{{colorYellow "Flags:"}}
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces | colorHelp}}{{end}}{{if .HasAvailableInheritedFlags}}

{{colorYellow "Global Options:"}}
  {{.InheritedFlags.FlagUsages | trimTrailingWhitespaces | colorHelp}}{{end}}{{if .HasHelpSubCommands}}

{{colorYellow "Additional help topics:"}}
{{range .Commands}}{{if .IsHelpCommand}}
  {{colorGreen (rpad .CommandPath .CommandPathPadding) }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasSubCommands}}

{{colorYellow (printf "Use \"%s [command] --help\" for more information about a command." .CommandPath)}}{{end}}
`

func SetUsageDefinition(cmd *cobra.Command) {
	cobra.AddTemplateFunc("colorYellow", colorYellow)
	cobra.AddTemplateFunc("colorGreen", colorGreen)
	cobra.AddTemplateFunc("colorRed", colorRed)
	cobra.AddTemplateFunc("colorBlue", colorBlue)
	cobra.AddTemplateFunc("colorHelp", colorHelp)
	cobra.AddTemplateFunc("hasServiceCommands", hasServiceCommands)
	cobra.AddTemplateFunc("hasModuleCommands", hasModuleCommands)

	// Altera o template de uso do cobra
	cmd.SetUsageTemplate(cliUsageTemplate)
}
