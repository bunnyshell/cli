package environment

import (
	"bunnyshell.com/cli/cmd/environment/action"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:     "environments",
	Aliases: []string{"env"},

	Short: "Bunnyshell Environments",

	ValidArgsFunction: cobra.NoFileCompletions,
}

func init() {
	config.MainManager.CommandWithAPI(mainCmd)

	util.AddGroupedCommands(
		mainCmd,
		cobra.Group{
			ID:    "actions",
			Title: "Environment Actions",
		},
		action.GetMainCommand().Commands(),
	)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
