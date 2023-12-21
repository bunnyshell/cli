package interactive

import (
	"fmt"

	"bunnyshell.com/cli/pkg/util"
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

	if required[0] != util.StrTrue {
		return
	}

	ensure(flag)

	required[0] = util.StrFalse
}

func ensure(flag *pflag.Flag) {
	question := getQuestion(flag)

	for {
		if err := question(); err == nil {
			break
		}
	}
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

func getQuestion(flag *pflag.Flag) func() error {
	message := fmt.Sprintf("Provide a value for flag '%s':", flag.Name)

	minimumLength := 1
	if util.GetFlagBoolAnnotation(flag, util.FlagAllowBlank) {
		minimumLength = 0
	}

	validator := All(
		AssertMinimumLength(minimumLength),
		setterValidator(flag),
	)

	if !util.HasHelp(flag) {
		return func() error {
			_, err := Ask(message, validator)

			return err
		}
	}

	help := util.GetHelp(flag)

	if util.IsHidden(flag) {
		return func() error {
			_, err := AskSecretWithHelp(message, help, validator)

			return err
		}
	}

	return func() error {
		_, err := AskWithHelp(message, help, validator)

		return err
	}
}
