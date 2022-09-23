package util

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

func Ask(question string, validate survey.Validator) (string, error) {
	answer := ""

	err := survey.AskOne(&survey.Input{
		Message: question,
	}, &answer, withValidator(validate))

	if err == terminal.InterruptErr {
		log.Fatal("interrupted")
	}

	return answer, err
}

func AskSecretWithHelp(question string, help string, validate survey.Validator) (string, error) {
	answer := ""

	err := survey.AskOne(&survey.Password{
		Message: question,
		Help:    help,
	}, &answer, withValidator(validate))

	if err == terminal.InterruptErr {
		log.Fatal("interrupted")
	}

	return answer, err
}

func AskPath(question string, value string, validate survey.Validator) (string, error) {
	answer := ""

	err := survey.AskOne(&survey.Input{
		Message: question,
		Default: value,
		Suggest: suggestPaths,
	}, &answer, withValidator(validate))

	if err == terminal.InterruptErr {
		log.Fatal("interrupted")
	}

	return answer, err
}

func Confirm(question string) (bool, error) {
	answer := false

	err := survey.AskOne(&survey.Confirm{
		Message: question,
	}, &answer)

	if err == terminal.InterruptErr {
		log.Fatal("interrupted")
	}

	return answer, err
}

func Choose(question string, items []string) (int, string, error) {
	answerIndex := 0

	err := survey.AskOne(&survey.Select{
		Message: question,
		Options: items,
	}, &answerIndex)

	if err == terminal.InterruptErr {
		log.Fatal("interrupted")
	}

	return answerIndex, items[answerIndex], err
}

func All(funcs ...survey.Validator) survey.Validator {
	return func(input interface{}) error {
		for _, callable := range funcs {
			err := callable(input)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func Lowercase() survey.Validator {
	return func(input interface{}) error {
		if strings.ToLower(input.(string)) != input {
			return fmt.Errorf("profile names should be lowercase only")
		}

		return nil
	}
}

func AssertMinimumLength(length int) survey.Validator {
	return survey.MinLength(length)
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
