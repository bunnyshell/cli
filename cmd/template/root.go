package template

import (
	"bunnyshell.com/cli/cmd/template/repository"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:     "templates",
	Aliases: []string{"tpl"},

	Short: "Template",
	Long:  "Bunnyshell Template",
}

var mainGroup = &cobra.Group{
	ID:    "templates",
	Title: "Commands for Templates:",
}

func init() {
	config.MainManager.CommandWithAPI(mainCmd)

	mainCmd.AddGroup(mainGroup)

	util.AddGroupedCommands(
		mainCmd,
		cobra.Group{
			ID:    "subresources",
			Title: "Commands for Template subresources:",
		},
		[]*cobra.Command{
			repository.GetMainCommand(),
		},
	)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
