package variable

import (
	"bunnyshell.com/cli/cmd/component/variable/action"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:     "variables",
	Aliases: []string{"vars"},

	Short: "Component Variables",
}

var mainGroup = &cobra.Group{
	ID:    "variables",
	Title: "Commands for Component Variables:",
}

func init() {
	config.MainManager.CommandWithAPI(mainCmd)

	mainCmd.AddGroup(mainGroup)

	util.AddGroupedCommands(
		mainCmd,
		cobra.Group{
			ID:    "actions",
			Title: "Commands for Component variables Actions:",
		},
		action.GetMainCommand().Commands(),
	)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
