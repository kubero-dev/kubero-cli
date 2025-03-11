package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/i582/cfmt/cmd/cfmt"
)

type ConsolePrompt struct{}

func (p *ConsolePrompt) PromptLine(question, options, def string) string {
	if def != "" {
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

func (p *ConsolePrompt) SelectFromList(question string, options []string, def string) string {
	_, _ = cfmt.Println("")
	if def != "" {
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

func NewConsolePrompt() *ConsolePrompt {
	return &ConsolePrompt{}
}
