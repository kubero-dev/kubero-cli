package utils

// Prompt defines the interface for prompting user input
type Prompt interface {
	PromptLine(question, options, def string) string
	SelectFromList(question string, options []string, def string) string
	ConfirmationLine(question string, def string) bool
}
