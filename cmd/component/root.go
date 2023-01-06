package component

import (
	"bunnyshell.com/cli/cmd/component/action"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:     "components",
	Aliases: []string{"comp"},

	Short: "Bunnyshell Components",
}

func init() {
	config.MainManager.CommandWithAPI(mainCmd)

	util.AddGroupedCommands(
		mainCmd,
		cobra.Group{
			ID:    "actions",
			Title: "Component Actions",
		},
		action.GetMainCommand().Commands(),
	)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
