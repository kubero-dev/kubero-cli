package common

import (
	"os"
	"strings"
)

func GetDescriptions(descriptionArg []string, hideBanner bool) map[string]string {
	var description, banner string

	if strings.Contains(strings.Join(os.Args[0:], ""), "-h") {
		description = descriptionArg[0]
	} else {
		if len(descriptionArg) > 1 {
			description = descriptionArg[1]
		} else {
			description = descriptionArg[0]
		}
	}

	if !hideBanner {
		banner = `
	,--. ,--.        ,--.
	|  .'   /,--.,--.|  |-.  ,---. ,--.--. ,---.
	|  .   ' |  ||  || .-. '| .-. :|  .--'| .-. |
	|  |\   \'  ''  '| '-' |\   --.|  |   ' '-' '
	'--' '--' '----'  '---'  '----''--'    '---'
	    Documentation: https://docs.kubero.dev
`
	} else {
		banner = ""
	}
	return map[string]string{"banner": banner, "description": description}
}

func ConcatenateExamples(example ...string) string {
	examples := ""
	for _, exp := range example {
		examples += string(exp) + "\n  "
	}
	return examples
}
