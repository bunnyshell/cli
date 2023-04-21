package variable

import (
	"fmt"

	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/config/option"
	"github.com/spf13/cobra"
)

var mainGroup = cobra.Group{
	ID:    "variables",
	Title: "Commands for Environment Variables:",
}

var mainCmd = &cobra.Command{
	Use:     "variables",
	Aliases: []string{"var"},

	Short: "Environment Variables",
	Long:  "Bunnyshell Environment Variables",
}

func init() {
	config.MainManager.CommandWithAPI(mainCmd)

	mainCmd.AddGroup(&mainGroup)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}

func getIDOption(value *string) *option.String {
	help := fmt.Sprintf(
		`Find available variables with "%s variables list"`,
		build.Name,
	)

	idOption := option.NewStringOption(value)

	idOption.AddFlagWithExtraHelp("id", "Environment Variable Id", help)

	return idOption
}
