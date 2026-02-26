package services

type PromptService interface {
	Confirm(message string) (bool, error)
	Ask(message string) (bool, error)
	Select(message string, opt []string) (string, error)
}
