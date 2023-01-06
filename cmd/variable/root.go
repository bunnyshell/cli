package variable

import (
	"bunnyshell.com/cli/pkg/config"
	"github.com/spf13/cobra"
)

var mainGroup = cobra.Group{
	ID:    "variables",
	Title: "Commands for Environment Variables:",
}

var mainCmd = &cobra.Command{
	Use:     "variables",
	Aliases: []string{"var"},

	Short: "Environment Variables",
	Long:  "Bunnyshell Environment Variables",
}

func init() {
	config.MainManager.CommandWithAPI(mainCmd)

	mainCmd.AddGroup(&mainGroup)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
