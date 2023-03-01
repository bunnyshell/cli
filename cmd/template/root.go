package template

import (
	"bunnyshell.com/cli/pkg/config"
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
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
