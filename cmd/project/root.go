package project

import (
	"bunnyshell.com/cli/cmd/project/action"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:     "projects",
	Aliases: []string{"proj"},

	Short: "Projects",
	Long:  "Bunnyshell Projects",

	ValidArgsFunction: cobra.NoFileCompletions,
}

func init() {
	config.MainManager.CommandWithAPI(mainCmd)

	util.AddGroupedCommands(
		mainCmd,
		cobra.Group{
			ID:    "actions",
			Title: "Commands for Project Actions:",
		},
		action.GetMainCommand().Commands(),
	)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
