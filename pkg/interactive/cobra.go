package interactive

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func AskMissingRequiredFlags(command *cobra.Command) {
	if getSettings().NonInteractive {
		return
	}

	command.Flags().VisitAll(flagOrInteractive)
}

func flagOrInteractive(flag *pflag.Flag) {
	if flag.Changed {
		return
	}

	required, hasRequired := flag.Annotations[cobra.BashCompOneRequiredFlag]
	if !hasRequired {
		return
	}

	if required[0] != "true" {
		return
	}

	question := fmt.Sprintf("Provide a value for '%s':", flag.Name)
	validator := All(
		AssertMinimumLength(1),
		setterValidator(flag),
	)

	for {
		if _, err := Ask(question, validator); err == nil {
			break
		}
	}

	required[0] = "false"
}

func setterValidator(flag *pflag.Flag) survey.Validator {
	return func(input interface{}) error {
		str, ok := input.(string)
		if !ok {
			return ErrInvalidValue
		}

		if err := flag.Value.Set(str); err != nil {
			return fmt.Errorf("%w: needs to be a %s", ErrInvalidValue, flag.Value.Type())
		}

		flag.Changed = true

		return nil
	}
}
