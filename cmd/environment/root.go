package environment

import (
	"bunnyshell.com/cli/cmd/environment/action"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

var mainGroup = &cobra.Group{
	ID:    "environments",
	Title: "Commands for Environment:",
}

var mainCmd = &cobra.Command{
	Use:     "environments",
	Aliases: []string{"env"},

	Short: "Environments",
	Long:  "Bunnyshell Environments",

	ValidArgsFunction: cobra.NoFileCompletions,
}

func init() {
	config.MainManager.CommandWithAPI(mainCmd)

	mainCmd.AddGroup(mainGroup)

	util.AddGroupedCommands(
		mainCmd,
		cobra.Group{
			ID:    "actions",
			Title: "Commands for Environment Actions:",
		},
		action.GetMainCommand().Commands(),
	)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
