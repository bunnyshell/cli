package action

import (
	"errors"
	"fmt"

	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/config/option"
	"github.com/spf13/cobra"
)

var (
	errMissingValue        = errors.New("the plain value must be provided")
	errMultipleValueInputs = errors.New("the value must be provided either by argument or by stdin, not both")
)

var mainCmd = &cobra.Command{}

func GetMainCommand() *cobra.Command {
	return mainCmd
}

func GetIDOption(value *string) *option.String {
	help := fmt.Sprintf(
		`Find available environment variables with "%s variables list"`,
		build.Name,
	)

	idOption := option.NewStringOption(value)

	idOption.AddFlagWithExtraHelp("id", "Environment Variable Id", help)

	return idOption
}
