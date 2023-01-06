package interactive

import (
	"errors"
	"log"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

func Ask(question string, validate survey.Validator) (string, error) {
	var answer string

	return answer, askPrompt(&survey.Input{
		Message: question,
	}, &answer, validate)
}

func AskInt32(question string, validate survey.Validator) (int32, error) {
	var answer int32

	return answer, askPrompt(&survey.Input{
		Message: question,
	}, &answer, validate)
}

func AskSecretWithHelp(question string, help string, validate survey.Validator) (string, error) {
	var answer string

	return answer, askPrompt(&survey.Password{
		Message: question,
		Help:    help,
	}, &answer, validate)
}

func AskPath(question string, value string, validate survey.Validator) (string, error) {
	var answer string

	return answer, askPrompt(&survey.Input{
		Message: question,
		Default: value,
		Suggest: suggestPaths,
	}, &answer, validate)
}

func Confirm(question string) (bool, error) {
	var answer bool

	return answer, askPrompt(&survey.Confirm{
		Message: question,
	}, &answer, nil)
}

func Choose(question string, items []string) (int, string, error) {
	var answerIndex int

	return answerIndex, items[answerIndex], askPrompt(&survey.Select{
		Message: question,
		Options: items,
	}, &answerIndex, nil)
}

func askPrompt(input survey.Prompt, answer any, validate survey.Validator) error {
	// safeguard: it should really be handled upstream
	if getSettings().NonInteractive {
		return ErrNonInteractive
	}

	err := survey.AskOne(input, answer, withValidator(validate))

	if errors.Is(err, terminal.InterruptErr) {
		log.Fatal("interrupted")
	}

	return err
}

func withValidator(validate survey.Validator) survey.AskOpt {
	if validate == nil {
		return nil
	}

	return survey.WithValidator(validate)
}

func suggestPaths(toComplete string) []string {
	files, _ := filepath.Glob(toComplete + "*")

	return files
}
