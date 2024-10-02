package variable_group

import (
	"bunnyshell.com/cli/cmd/variable_group/action"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:     "variables-groups",
	Aliases: []string{"variables-group", "var-groups", "var-group", "var-g"},

	Short: "Grouped Environment Variables",
	Long:  "Bunnyshell Environment Variables in Groups",
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
