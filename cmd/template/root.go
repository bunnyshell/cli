package template

import (
	"fmt"

	"bunnyshell.com/cli/cmd/template/repository"
	"bunnyshell.com/cli/pkg/build"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/config/option"
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

func getIDOption(value *string) *option.String {
	help := fmt.Sprintf(
		`Find available Templates with "%s templates list"`,
		build.Name,
	)

	option := option.NewStringOption(value)

	option.AddFlagWithExtraHelp("id", "Template ID", help)

	return option
}
