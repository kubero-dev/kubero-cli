package utils

import (
	"bufio"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/faelmori/kubero-cli/internal/log"
	"os"
	"strings"
)

type ConsolePrompt struct{}

func (p *ConsolePrompt) PromptLine(question, options, def string) string {
	if def != "" {
		log.Info(fmt.Sprintf("%s %s : %s", question, options, def))
		return def
	}
	reader := bufio.NewReader(os.Stdin)
	log.Info(fmt.Sprintf("%s %s : %s", question, options, def))
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" {
		text = def
	}
	return text
}
func (p *ConsolePrompt) SelectFromList(question string, options []string, def string) string {
	log.Println("")
	if def != "" {
		log.Info(question, def)
		return def
	}
	prompt := &survey.Select{
		Message: question,
		Options: options,
	}
	askOneErr := survey.AskOne(prompt, &def)
	if askOneErr != nil {
		log.Error("failed to ask for input", askOneErr)
		return ""
	}
	return def
}

func NewConsolePrompt() *ConsolePrompt { return &ConsolePrompt{} }
