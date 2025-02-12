package utils

import (
	"bunnyshell.com/cli/cmd/git"
	"bunnyshell.com/cli/cmd/component_debug"
	"bunnyshell.com/cli/cmd/remote_development"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{}

func init() {
	mainCmd.AddCommand(git.GetMainCommand())
	mainCmd.AddCommand(remote_development.GetMainCommand())
	mainCmd.AddCommand(component_debug.GetMainCommand())
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
