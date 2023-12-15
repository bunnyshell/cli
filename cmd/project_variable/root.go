package project_variable

import (
	"bunnyshell.com/cli/cmd/project_variable/action"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:     "project-variables",
	Aliases: []string{"pvar"},

	Short: "Project Variables",
	Long:  "Bunnyshell Project Variables",
}

func init() {
	config.MainManager.CommandWithAPI(mainCmd)

	util.AddGroupedCommands(
		mainCmd,
		cobra.Group{
			ID:    "actions",
			Title: "Commands for Project Variables:",
		},
		action.GetMainCommand().Commands(),
	)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
