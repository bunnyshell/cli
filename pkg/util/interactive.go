package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
)

func Ask(question string, validate promptui.ValidateFunc) (string, error) {
	prompt := promptui.Prompt{
		Label: question,
		Templates: &promptui.PromptTemplates{
			Prompt:  "{{ . }} ",
			Valid:   "{{ . | green }} ",
			Invalid: "{{ . | red }} ",
			Success: "{{ . | bold }} ",
		},
		Validate: validate,
	}

	s, e := prompt.Run()
	if e == promptui.ErrInterrupt || e == promptui.ErrEOF {
		os.Exit(0)
	}

	return s, e
}

func AskDefault(question string, value string, validate promptui.ValidateFunc) string {
	prompt := promptui.Prompt{
		Label: question,
		Templates: &promptui.PromptTemplates{
			Prompt:  "{{ . }} ",
			Valid:   "{{ . | green }} ",
			Invalid: "{{ . | red }} ",
			Success: "{{ . | bold }} ",
		},
		Default:   value,
		AllowEdit: false,
		Validate:  validate,
	}

	s, e := prompt.Run()
	if e == promptui.ErrInterrupt || e == promptui.ErrEOF {
		os.Exit(0)
	}

	if s == "" {
		return value
	}

	return s
}

func AskSecret(question string, validate promptui.ValidateFunc) (string, error) {
	prompt := promptui.Prompt{
		Label: question,
		Templates: &promptui.PromptTemplates{
			Prompt:  "{{ . }} ",
			Valid:   "{{ . | green }} ",
			Invalid: "{{ . | red }} ",
			Success: "{{ . | bold }} ",
		},
		Mask:     '*',
		Validate: validate,
	}

	s, e := prompt.Run()
	if e == promptui.ErrInterrupt || e == promptui.ErrEOF {
		os.Exit(0)
	}

	return s, e
}

func AskWithDefault(question, defaultInput string) (string, error) {
	prompt := promptui.Prompt{
		Label:   question,
		Default: defaultInput,
	}

	return prompt.Run()
}

func Confirm(question string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     question,
		IsConfirm: true,
	}

	r, e := prompt.Run()

	if e != nil {
		if e == promptui.ErrInterrupt || e == promptui.ErrEOF {
			os.Exit(0)
		}

		if e.Error() != "" {
			return false, e
		}
	}

	return r == "y" || r == "Y", nil
}

func Choose(question string, items []string) (int, string, error) {
	prompt := promptui.Select{
		Label: question,
		Items: items,
	}

	s, r, e := prompt.Run()
	if e == promptui.ErrInterrupt || e == promptui.ErrEOF {
		os.Exit(0)
	}

	return s, r, e
}

func ChooseOrOther(question string, items []string, other string) (int, string, error) {
	prompt := promptui.SelectWithAdd{
		Label:    question,
		Items:    items,
		AddLabel: other,
	}

	i, s, e := prompt.Run()
	if e == promptui.ErrInterrupt || e == promptui.ErrEOF {
		os.Exit(0)
	}
	return i, s, e
}

func ConfigFileValidation(input string) error {
	if input == "" {
		return nil
	}

	ext := filepath.Ext(input)
	// sumimasen
	switch ext {
	case ".json", ".yaml":
		return nil
	}

	return fmt.Errorf("supported extensions: json or yaml")
}

func All(funcs ...promptui.ValidateFunc) promptui.ValidateFunc {
	return func(input string) error {
		for _, callable := range funcs {
			err := callable(input)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func Lowercase() promptui.ValidateFunc {
	return func(input string) error {
		if strings.ToLower(input) != input {
			return fmt.Errorf("profile names should be lowercase only")
		}

		return nil
	}
}

func RequiredExtension(extensions ...string) promptui.ValidateFunc {
	return func(input string) error {
		ext := filepath.Ext(input)
		// sumimasen

		for _, allowed := range extensions {
			if ext == allowed {
				return nil
			}
		}

		return fmt.Errorf("supported extensions: %v", extensions)
	}
}

func OptionalMinimumLength(length int) promptui.ValidateFunc {
	return func(input string) error {
		if input == "" {
			return nil
		}

		return minimumLength(input, length)
	}
}

func AssertMinimumLength(length int) promptui.ValidateFunc {
	return func(input string) error {
		return minimumLength(input, length)
	}
}

func minimumLength(input string, length int) error {
	if len(input) <= length {
		return fmt.Errorf("input at least %d characters", length)
	}

	return nil
}
