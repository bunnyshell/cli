package remote_development

import (
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:     "remote-development",
	Aliases: []string{"rdev"},

	Short: "Remote Development",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		lib.LoadViperConfigIntoContext()
	},
}

func init() {
	lib.CLIContext.RequireTokenOnCommand(mainCmd)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
