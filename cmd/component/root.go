package component

import (
	"bunnyshell.com/cli/cmd/component/action"
	"bunnyshell.com/cli/cmd/component/variable"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:     "components",
	Aliases: []string{"comp"},

	Short: "Components",
	Long:  "Bunnyshell Components",
}

var mainGroup = &cobra.Group{
	ID:    "components",
	Title: "Commands for Components:",
}

func init() {
	config.MainManager.CommandWithAPI(mainCmd)

	mainCmd.AddGroup(mainGroup)

	util.AddGroupedCommands(
		mainCmd,
		cobra.Group{
			ID:    "actions",
			Title: "Commands for Component Actions:",
		},
		action.GetMainCommand().Commands(),
	)

	util.AddGroupedCommands(
		mainCmd,
		cobra.Group{
			ID:    "variables",
			Title: "Commands for Component Variables:",
		},
		[]*cobra.Command{variable.GetMainCommand()},
	)

	config.MainManager.CommandWithGlobalOptions(mainCmd)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
