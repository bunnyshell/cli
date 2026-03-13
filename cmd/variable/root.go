package variable

import (
	"bunnyshell.com/cli/cmd/variable/action"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:     "variables",
	Aliases: []string{"variable", "var", "vars"},

	Short: "Environment Variables",
	Long:  "Bunnyshell Environment Variables",
}

func init() {
	config.MainManager.CommandWithAPI(mainCmd)

	util.AddGroupedCommands(
		mainCmd,
		cobra.Group{
			ID:    "actions",
			Title: "Commands for Environment Variables:",
		},
		action.GetMainCommand().Commands(),
	)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
