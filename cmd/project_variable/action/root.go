package action

import (
	"fmt"

	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/config/option"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{}

func GetMainCommand() *cobra.Command {
	return mainCmd
}

func GetIDOption(value *string) *option.String {
	help := fmt.Sprintf(
		`Find available project variables with "%s variables list"`,
		build.Name,
	)

	idOption := option.NewStringOption(value)

	idOption.AddFlagWithExtraHelp("id", "Project Variable Id", help)

	return idOption
}
