package util

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
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

func AskInt32(question string, validate survey.Validator) (int32, error) {
	var answer int32 = 0

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

func AssertBetween(min int32, max int32) survey.Validator {
	return func(input interface{}) error {
		str, ok := input.(string)
		if !ok {
			return errors.New("invalid value")
		}

		i, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return errors.New("input must be an integer")
		}

		val := int32(i)
		if val < min || val > max {
			return fmt.Errorf("input must be between %d and %d", min, max)
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
