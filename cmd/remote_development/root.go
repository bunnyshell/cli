package remote_development

import (
	cfgCommand "bunnyshell.com/cli/cmd/remote_development/config"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:     "remote-development",
	Aliases: []string{"rdev"},

	Short: "Remote Development",
}

func init() {
	config.MainManager.CommandWithAPI(mainCmd)

	util.AddGroupedCommands(
		mainCmd,
		cobra.Group{
			ID:    "Config",
			Title: "Commands for config management:",
		},
		[]*cobra.Command{
			cfgCommand.GetMainCommand(),
		},
	)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
